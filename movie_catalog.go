// movie_catalog.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Movie struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Director  string `json:"director"`
	Year      int    `json:"year"`
	Genre     string `json:"genre"`
	Watched   bool   `json:"watched"`
	Rating    int    `json:"rating"`
	AddedDate string `json:"added_date"`
}

type MoviesData struct {
	Movies []Movie `json:"movies"`
}

type MovieCatalog struct {
	movies  []Movie
	nextID  int
}

func NewMovieCatalog() *MovieCatalog {
	return &MovieCatalog{
		movies:  []Movie{},
		nextID:  1,
	}
}

func (c *MovieCatalog) AddMovie(title, director string, year int, genre string, watched bool, rating int) (Movie, error) {
	if rating < 1 || rating > 10 {
		return Movie{}, fmt.Errorf("оценка должна быть от 1 до 10")
	}
	currentYear := time.Now().Year()
	if year < 1900 || year > currentYear {
		return Movie{}, fmt.Errorf("год должен быть от 1900 до %d", currentYear)
	}
	if title == "" || director == "" {
		return Movie{}, fmt.Errorf("название и режиссёр не могут быть пустыми")
	}
	if genre == "" {
		genre = "Неизвестно"
	}
	movie := Movie{
		ID:        c.nextID,
		Title:     title,
		Director:  director,
		Year:      year,
		Genre:     genre,
		Watched:   watched,
		Rating:    rating,
		AddedDate: time.Now().Format("2006-01-02"),
	}
	c.movies = append(c.movies, movie)
	c.nextID++
	return movie, nil
}

func (c *MovieCatalog) FindMovie(id int) *Movie {
	for i := range c.movies {
		if c.movies[i].ID == id {
			return &c.movies[i]
		}
	}
	return nil
}

func (c *MovieCatalog) EditMovie(id int, updates map[string]interface{}) bool {
	movie := c.FindMovie(id)
	if movie == nil {
		return false
	}
	for key, value := range updates {
		switch key {
		case "title":
			if v, ok := value.(string); ok {
				movie.Title = v
			}
		case "director":
			if v, ok := value.(string); ok {
				movie.Director = v
			}
		case "year":
			if v, ok := value.(int); ok {
				movie.Year = v
			}
		case "genre":
			if v, ok := value.(string); ok {
				movie.Genre = v
			}
		case "watched":
			if v, ok := value.(bool); ok {
				movie.Watched = v
			}
		case "rating":
			if v, ok := value.(int); ok {
				movie.Rating = v
			}
		}
	}
	return true
}

func (c *MovieCatalog) DeleteMovie(id int) bool {
	for i, m := range c.movies {
		if m.ID == id {
			c.movies = append(c.movies[:i], c.movies[i+1:]...)
			return true
		}
	}
	return false
}

func (c *MovieCatalog) SearchMovies(query string) []Movie {
	q := strings.ToLower(query)
	var result []Movie
	for _, m := range c.movies {
		if strings.Contains(strings.ToLower(m.Title), q) || strings.Contains(strings.ToLower(m.Director), q) {
			result = append(result, m)
		}
	}
	return result
}

func (c *MovieCatalog) FilterByWatched(watched bool) []Movie {
	var result []Movie
	for _, m := range c.movies {
		if m.Watched == watched {
			result = append(result, m)
		}
	}
	return result
}

func (c *MovieCatalog) FilterByGenre(genre string) []Movie {
	var result []Movie
	for _, m := range c.movies {
		if strings.EqualFold(m.Genre, genre) {
			result = append(result, m)
		}
	}
	return result
}

func (c *MovieCatalog) FilterByYear(year int) []Movie {
	var result []Movie
	for _, m := range c.movies {
		if m.Year == year {
			result = append(result, m)
		}
	}
	return result
}

func (c *MovieCatalog) GetStats() map[string]interface{} {
	total := len(c.movies)
	watched := len(c.FilterByWatched(true))
	unwatched := total - watched
	ratings := []int{}
	for _, m := range c.movies {
		if m.Watched {
			ratings = append(ratings, m.Rating)
		}
	}
	avgRating := 0.0
	if len(ratings) > 0 {
		sum := 0
		for _, r := range ratings {
			sum += r
		}
		avgRating = float64(sum) / float64(len(ratings))
	}
	genres := make(map[string]int)
	for _, m := range c.movies {
		genres[m.Genre]++
	}
	return map[string]interface{}{
		"total":       total,
		"watched":     watched,
		"unwatched":   unwatched,
		"avg_rating":  avgRating,
		"genres":      genres,
	}
}

func (c *MovieCatalog) SaveToFile(filename string) error {
	data := MoviesData{Movies: c.movies}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}

func (c *MovieCatalog) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var md MoviesData
	if err := json.Unmarshal(data, &md); err != nil {
		return err
	}
	c.movies = md.Movies
	for _, m := range c.movies {
		if m.ID >= c.nextID {
			c.nextID = m.ID + 1
		}
	}
	return nil
}

func readString(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt(prompt string) int {
	for {
		input := readString(prompt)
		if val, err := strconv.Atoi(input); err == nil {
			return val
		}
		fmt.Println("Введите число.")
	}
}

func readBool(prompt string) bool {
	for {
		input := readString(prompt)
		if input == "1" {
			return true
		} else if input == "0" {
			return false
		}
		fmt.Println("Введите 1 или 0.")
	}
}

func printMovie(movie Movie) {
	status := "✅ Просмотрен"
	if !movie.Watched {
		status = "⏳ Не просмотрен"
	}
	fmt.Printf("#%d - %s (%d)\n", movie.ID, movie.Title, movie.Year)
	fmt.Printf("   Режиссёр: %s, Жанр: %s\n", movie.Director, movie.Genre)
	fmt.Printf("   %s, Оценка: %d/10\n", status, movie.Rating)
	fmt.Printf("   Добавлен: %s\n", movie.AddedDate)
}

func main() {
	catalog := NewMovieCatalog()
	if err := catalog.LoadFromFile("movies_data.json"); err != nil {
		fmt.Println("Ошибка загрузки:", err)
	}

	for {
		fmt.Println("\n===== КАТАЛОГ ФИЛЬМОВ (Go) =====")
		fmt.Println("1. Добавить фильм")
		fmt.Println("2. Показать все фильмы")
		fmt.Println("3. Показать просмотренные фильмы")
		fmt.Println("4. Показать непросмотренные фильмы")
		fmt.Println("5. Найти фильмы по названию или режиссёру")
		fmt.Println("6. Отметить фильм как просмотренный")
		fmt.Println("7. Редактировать фильм")
		fmt.Println("8. Удалить фильм")
		fmt.Println("9. Показать статистику")
		fmt.Println("10. Сохранить в файл")
		fmt.Println("11. Загрузить из файла")
		fmt.Println("0. Выход")
		choice := readString("Выберите действие: ")

		switch choice {
		case "0":
			return
		case "1":
			title := readString("Название: ")
			if title == "" {
				fmt.Println("Название не может быть пустым.")
				continue
			}
			director := readString("Режиссёр: ")
			if director == "" {
				fmt.Println("Режиссёр не может быть пустым.")
				continue
			}
			year := readInt("Год выпуска: ")
			genre := readString("Жанр: ")
			if genre == "" {
				genre = "Неизвестно"
			}
			watched := readBool("Статус (1-просмотрен, 0-не просмотрен): ")
			rating := readInt("Оценка (1-10): ")
			movie, err := catalog.AddMovie(title, director, year, genre, watched, rating)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Printf("Фильм добавлен с ID %d\n", movie.ID)
			}
		case "2":
			if len(catalog.movies) == 0 {
				fmt.Println("Нет фильмов.")
			} else {
				for _, m := range catalog.movies {
					printMovie(m)
				}
			}
		case "3":
			movies := catalog.FilterByWatched(true)
			if len(movies) == 0 {
				fmt.Println("Нет просмотренных фильмов.")
			} else {
				for _, m := range movies {
					printMovie(m)
				}
			}
		case "4":
			movies := catalog.FilterByWatched(false)
			if len(movies) == 0 {
				fmt.Println("Нет непросмотренных фильмов.")
			} else {
				for _, m := range movies {
					printMovie(m)
				}
			}
		case "5":
			query := readString("Введите часть названия или режиссёра: ")
			results := catalog.SearchMovies(query)
			if len(results) == 0 {
				fmt.Println("Фильмы не найдены.")
			} else {
				for _, m := range results {
					printMovie(m)
				}
			}
		case "6":
			id := readInt("Введите ID фильма: ")
			movie := catalog.FindMovie(id)
			if movie == nil {
				fmt.Println("Фильм не найден.")
				continue
			}
			movie.Watched = true
			fmt.Println("Фильм отмечен как просмотренный.")
		case "7":
			id := readInt("Введите ID фильма для редактирования: ")
			movie := catalog.FindMovie(id)
			if movie == nil {
				fmt.Println("Фильм не найден.")
				continue
			}
			fmt.Println("Оставьте поле пустым, чтобы не менять.")
			newTitle := readString(fmt.Sprintf("Название (%s): ", movie.Title))
			newDirector := readString(fmt.Sprintf("Режиссёр (%s): ", movie.Director))
			newYear := readString(fmt.Sprintf("Год (%d): ", movie.Year))
			newGenre := readString(fmt.Sprintf("Жанр (%s): ", movie.Genre))
			newWatched := readString(fmt.Sprintf("Статус (1-просмотрен, 0-не просмотрен) сейчас: %d: ", map[bool]int{true: 1, false: 0}[movie.Watched]))
			newRating := readString(fmt.Sprintf("Оценка (%d): ", movie.Rating))
			updates := make(map[string]interface{})
			if newTitle != "" {
				updates["title"] = newTitle
			}
			if newDirector != "" {
				updates["director"] = newDirector
			}
			if newYear != "" {
				if y, err := strconv.Atoi(newYear); err == nil {
					updates["year"] = y
				} else {
					fmt.Println("Год должен быть числом, пропускаем.")
				}
			}
			if newGenre != "" {
				updates["genre"] = newGenre
			}
			if newWatched != "" {
				updates["watched"] = newWatched == "1"
			}
			if newRating != "" {
				if r, err := strconv.Atoi(newRating); err == nil {
					updates["rating"] = r
				} else {
					fmt.Println("Оценка должна быть числом, пропускаем.")
				}
			}
			if catalog.EditMovie(id, updates) {
				fmt.Println("Фильм обновлён.")
			} else {
				fmt.Println("Ошибка обновления.")
			}
		case "8":
			id := readInt("Введите ID фильма для удаления: ")
			if catalog.DeleteMovie(id) {
				fmt.Println("Фильм удалён.")
			} else {
				fmt.Println("Фильм не найден.")
			}
		case "9":
			stats := catalog.GetStats()
			fmt.Println("\n=== СТАТИСТИКА ===")
			fmt.Printf("Всего фильмов: %d\n", stats["total"])
			fmt.Printf("Просмотрено: %d\n", stats["watched"])
			fmt.Printf("Не просмотрено: %d\n", stats["unwatched"])
			fmt.Printf("Средняя оценка (только просмотренные): %.2f\n", stats["avg_rating"])
			fmt.Println("По жанрам:")
			genres := stats["genres"].(map[string]int)
			for g, c := range genres {
				fmt.Printf("  %s: %d\n", g, c)
			}
		case "10":
			if err := catalog.SaveToFile("movies_data.json"); err != nil {
				fmt.Println("Ошибка сохранения:", err)
			} else {
				fmt.Println("Сохранено.")
			}
		case "11":
			if err := catalog.LoadFromFile("movies_data.json"); err != nil {
				fmt.Println("Ошибка загрузки:", err)
			} else {
				fmt.Println("Загружено.")
			}
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
