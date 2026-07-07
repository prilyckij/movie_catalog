// MovieCatalog.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Serialization;

public record Movie(
    int Id,
    string Title,
    string Director,
    int Year,
    string Genre,
    bool Watched,
    int Rating,
    string AddedDate
);

public class MoviesData
{
    public List<Movie> Movies { get; set; } = new();
}

public class MovieCatalog
{
    private List<Movie> movies = new();
    private int nextId = 1;

    public IReadOnlyList<Movie> Movies => movies.AsReadOnly();

    public Movie AddMovie(string title, string director, int year, string genre, bool watched, int rating)
    {
        if (rating < 1 || rating > 10) throw new ArgumentException("Оценка должна быть от 1 до 10");
        if (year < 1900 || year > DateTime.Now.Year)
            throw new ArgumentException($"Год должен быть от 1900 до {DateTime.Now.Year}");
        if (string.IsNullOrWhiteSpace(title) || string.IsNullOrWhiteSpace(director))
            throw new ArgumentException("Название и режиссёр не могут быть пустыми");
        if (string.IsNullOrWhiteSpace(genre)) genre = "Неизвестно";
        var movie = new Movie(nextId, title, director, year, genre, watched, rating, DateTime.Now.ToString("yyyy-MM-dd"));
        movies.Add(movie);
        nextId++;
        return movie;
    }

    public Movie? FindMovie(int id) => movies.FirstOrDefault(m => m.Id == id);

    public bool EditMovie(int id, Dictionary<string, object> updates)
    {
        var old = FindMovie(id);
        if (old == null) return false;
        movies.Remove(old);
        string title = updates.ContainsKey("title") ? (string)updates["title"] : old.Title;
        string director = updates.ContainsKey("director") ? (string)updates["director"] : old.Director;
        int year = updates.ContainsKey("year") ? (int)updates["year"] : old.Year;
        string genre = updates.ContainsKey("genre") ? (string)updates["genre"] : old.Genre;
        bool watched = updates.ContainsKey("watched") ? (bool)updates["watched"] : old.Watched;
        int rating = updates.ContainsKey("rating") ? (int)updates["rating"] : old.Rating;
        var updated = new Movie(old.Id, title, director, year, genre, watched, rating, old.AddedDate);
        movies.Add(updated);
        return true;
    }

    public bool DeleteMovie(int id) => movies.RemoveAll(m => m.Id == id) > 0;

    public List<Movie> SearchMovies(string query)
    {
        var q = query.ToLower();
        return movies.Where(m => m.Title.ToLower().Contains(q) || m.Director.ToLower().Contains(q)).ToList();
    }

    public List<Movie> FilterByWatched(bool watched) => movies.Where(m => m.Watched == watched).ToList();

    public List<Movie> FilterByGenre(string genre) =>
        movies.Where(m => string.Equals(m.Genre, genre, StringComparison.OrdinalIgnoreCase)).ToList();

    public List<Movie> FilterByYear(int year) => movies.Where(m => m.Year == year).ToList();

    public Dictionary<string, object> GetStats()
    {
        int total = movies.Count;
        int watchedCount = FilterByWatched(true).Count;
        int unwatched = total - watchedCount;
        double avgRating = FilterByWatched(true).Any() ? FilterByWatched(true).Average(m => m.Rating) : 0;
        var genres = movies.GroupBy(m => m.Genre).ToDictionary(g => g.Key, g => g.Count());
        return new Dictionary<string, object>
        {
            ["total"] = total,
            ["watched"] = watchedCount,
            ["unwatched"] = unwatched,
            ["avg_rating"] = avgRating,
            ["genres"] = genres
        };
    }

    public void SaveToFile(string filename)
    {
        var data = new MoviesData { Movies = movies };
        var options = new JsonSerializerOptions { WriteIndented = true };
        string json = JsonSerializer.Serialize(data, options);
        File.WriteAllText(filename, json);
    }

    public void LoadFromFile(string filename)
    {
        if (!File.Exists(filename)) return;
        string json = File.ReadAllText(filename);
        var data = JsonSerializer.Deserialize<MoviesData>(json);
        if (data != null)
        {
            movies = data.Movies;
            nextId = movies.Any() ? movies.Max(m => m.Id) + 1 : 1;
        }
    }
}

public static class Program
{
    private static string ReadString(string prompt)
    {
        Console.Write(prompt);
        return Console.ReadLine()?.Trim() ?? "";
    }

    private static int ReadInt(string prompt)
    {
        while (true)
        {
            Console.Write(prompt);
            if (int.TryParse(Console.ReadLine(), out int result))
                return result;
            Console.WriteLine("Введите число.");
        }
    }

    private static bool ReadBool(string prompt)
    {
        while (true)
        {
            string input = ReadString(prompt);
            if (input == "1") return true;
            if (input == "0") return false;
            Console.WriteLine("Введите 1 или 0.");
        }
    }

    private static void PrintMovie(Movie movie)
    {
        string status = movie.Watched ? "✅ Просмотрен" : "⏳ Не просмотрен";
        Console.WriteLine($"#{movie.Id} - {movie.Title} ({movie.Year})");
        Console.WriteLine($"   Режиссёр: {movie.Director}, Жанр: {movie.Genre}");
        Console.WriteLine($"   {status}, Оценка: {movie.Rating}/10");
        Console.WriteLine($"   Добавлен: {movie.AddedDate}");
    }

    public static void Main()
    {
        var catalog = new MovieCatalog();
        try { catalog.LoadFromFile("movies_data.json"); }
        catch { Console.WriteLine("Не удалось загрузить данные."); }

        while (true)
        {
            Console.WriteLine("\n===== КАТАЛОГ ФИЛЬМОВ (C#) =====");
            Console.WriteLine("1. Добавить фильм");
            Console.WriteLine("2. Показать все фильмы");
            Console.WriteLine("3. Показать просмотренные фильмы");
            Console.WriteLine("4. Показать непросмотренные фильмы");
            Console.WriteLine("5. Найти фильмы по названию или режиссёру");
            Console.WriteLine("6. Отметить фильм как просмотренный");
            Console.WriteLine("7. Редактировать фильм");
            Console.WriteLine("8. Удалить фильм");
            Console.WriteLine("9. Показать статистику");
            Console.WriteLine("10. Сохранить в файл");
            Console.WriteLine("11. Загрузить из файла");
            Console.WriteLine("0. Выход");
            string choice = ReadString("Выберите действие: ");

            switch (choice)
            {
                case "0": return;
                case "1":
                    string title = ReadString("Название: ");
                    if (string.IsNullOrWhiteSpace(title)) { Console.WriteLine("Название не может быть пустым."); continue; }
                    string director = ReadString("Режиссёр: ");
                    if (string.IsNullOrWhiteSpace(director)) { Console.WriteLine("Режиссёр не может быть пустым."); continue; }
                    int year = ReadInt("Год выпуска: ");
                    string genre = ReadString("Жанр: ");
                    if (string.IsNullOrWhiteSpace(genre)) genre = "Неизвестно";
                    bool watched = ReadBool("Статус (1-просмотрен, 0-не просмотрен): ");
                    int rating = ReadInt("Оценка (1-10): ");
                    try
                    {
                        var movie = catalog.AddMovie(title, director, year, genre, watched, rating);
                        Console.WriteLine($"Фильм добавлен с ID {movie.Id}");
                    }
                    catch (Exception ex) { Console.WriteLine($"Ошибка: {ex.Message}"); }
                    break;
                case "2":
                    if (!catalog.Movies.Any()) Console.WriteLine("Нет фильмов.");
                    else foreach (var m in catalog.Movies) PrintMovie(m);
                    break;
                case "3":
                    var watchedMovies = catalog.FilterByWatched(true);
                    if (!watchedMovies.Any()) Console.WriteLine("Нет просмотренных фильмов.");
                    else foreach (var m in watchedMovies) PrintMovie(m);
                    break;
                case "4":
                    var unwatchedMovies = catalog.FilterByWatched(false);
                    if (!unwatchedMovies.Any()) Console.WriteLine("Нет непросмотренных фильмов.");
                    else foreach (var m in unwatchedMovies) PrintMovie(m);
                    break;
                case "5":
                    string query = ReadString("Введите часть названия или режиссёра: ");
                    var results = catalog.SearchMovies(query);
                    if (!results.Any()) Console.WriteLine("Фильмы не найдены.");
                    else foreach (var m in results) PrintMovie(m);
                    break;
                case "6":
                    int id = ReadInt("Введите ID фильма: ");
                    var movie = catalog.FindMovie(id);
                    if (movie == null) { Console.WriteLine("Фильм не найден."); continue; }
                    catalog.EditMovie(id, new Dictionary<string, object> { ["watched"] = true });
                    Console.WriteLine("Фильм отмечен как просмотренный.");
                    break;
                case "7":
                    int eid = ReadInt("Введите ID фильма для редактирования: ");
                    var old = catalog.FindMovie(eid);
                    if (old == null) { Console.WriteLine("Фильм не найден."); continue; }
                    Console.WriteLine("Оставьте поле пустым, чтобы не менять.");
                    string newTitle = ReadString($"Название ({old.Title}): ");
                    string newDirector = ReadString($"Режиссёр ({old.Director}): ");
                    string newYearStr = ReadString($"Год ({old.Year}): ");
                    string newGenre = ReadString($"Жанр ({old.Genre}): ");
                    string newWatchedStr = ReadString($"Статус (1-просмотрен, 0-не просмотрен) сейчас: {(old.Watched ? "1" : "0")}: ");
                    string newRatingStr = ReadString($"Оценка ({old.Rating}): ");
                    var updates = new Dictionary<string, object>();
                    if (!string.IsNullOrWhiteSpace(newTitle)) updates["title"] = newTitle;
                    if (!string.IsNullOrWhiteSpace(newDirector)) updates["director"] = newDirector;
                    if (!string.IsNullOrWhiteSpace(newYearStr))
                    {
                        if (int.TryParse(newYearStr, out int y)) updates["year"] = y;
                        else Console.WriteLine("Год должен быть числом, пропускаем.");
                    }
                    if (!string.IsNullOrWhiteSpace(newGenre)) updates["genre"] = newGenre;
                    if (!string.IsNullOrWhiteSpace(newWatchedStr)) updates["watched"] = newWatchedStr == "1";
                    if (!string.IsNullOrWhiteSpace(newRatingStr))
                    {
                        if (int.TryParse(newRatingStr, out int r)) updates["rating"] = r;
                        else Console.WriteLine("Оценка должна быть числом, пропускаем.");
                    }
                    if (catalog.EditMovie(eid, updates)) Console.WriteLine("Фильм обновлён.");
                    else Console.WriteLine("Ошибка обновления.");
                    break;
                case "8":
                    int delId = ReadInt("Введите ID фильма для удаления: ");
                    if (catalog.DeleteMovie(delId)) Console.WriteLine("Фильм удалён.");
                    else Console.WriteLine("Фильм не найден.");
                    break;
                case "9":
                    var stats = catalog.GetStats();
                    Console.WriteLine("\n=== СТАТИСТИКА ===");
                    Console.WriteLine($"Всего фильмов: {stats["total"]}");
                    Console.WriteLine($"Просмотрено: {stats["watched"]}");
                    Console.WriteLine($"Не просмотрено: {stats["unwatched"]}");
                    Console.WriteLine($"Средняя оценка (только просмотренные): {stats["avg_rating"]:F2}");
                    Console.WriteLine("По жанрам:");
                    var genres = (Dictionary<string, int>)stats["genres"];
                    foreach (var kv in genres) Console.WriteLine($"  {kv.Key}: {kv.Value}");
                    break;
                case "10":
                    try { catalog.SaveToFile("movies_data.json"); Console.WriteLine("Сохранено."); }
                    catch (Exception ex) { Console.WriteLine($"Ошибка: {ex.Message}"); }
                    break;
                case "11":
                    try { catalog.LoadFromFile("movies_data.json"); Console.WriteLine("Загружено."); }
                    catch (Exception ex) { Console.WriteLine($"Ошибка: {ex.Message}"); }
                    break;
                default: Console.WriteLine("Неизвестная команда."); break;
            }
        }
    }
}
