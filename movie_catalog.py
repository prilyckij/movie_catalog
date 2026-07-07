# movie_catalog.py
import json
from dataclasses import dataclass, asdict
from datetime import date
from typing import List, Optional
from pathlib import Path

@dataclass
class Movie:
    id: int
    title: str
    director: str
    year: int
    genre: str
    watched: bool
    rating: int  # 1-10
    added_date: str

class MovieCatalog:
    def __init__(self):
        self.movies: List[Movie] = []
        self.next_id = 1

    def add_movie(self, title: str, director: str, year: int, genre: str, watched: bool, rating: int) -> Movie:
        if rating < 1 or rating > 10:
            raise ValueError("Оценка должна быть от 1 до 10")
        if year < 1900 or year > date.today().year:
            raise ValueError(f"Год должен быть от 1900 до {date.today().year}")
        if not title or not director:
            raise ValueError("Название и режиссёр не могут быть пустыми")
        movie = Movie(
            id=self.next_id,
            title=title,
            director=director,
            year=year,
            genre=genre or "Неизвестно",
            watched=watched,
            rating=rating,
            added_date=date.today().isoformat()
        )
        self.movies.append(movie)
        self.next_id += 1
        return movie

    def find_movie(self, movie_id: int) -> Optional[Movie]:
        return next((m for m in self.movies if m.id == movie_id), None)

    def edit_movie(self, movie_id: int, **kwargs) -> bool:
        movie = self.find_movie(movie_id)
        if not movie:
            return False
        for key, value in kwargs.items():
            if hasattr(movie, key) and value is not None:
                setattr(movie, key, value)
        return True

    def delete_movie(self, movie_id: int) -> bool:
        movie = self.find_movie(movie_id)
        if movie:
            self.movies.remove(movie)
            return True
        return False

    def search_movies(self, query: str) -> List[Movie]:
        q = query.lower()
        return [m for m in self.movies if q in m.title.lower() or q in m.director.lower()]

    def filter_by_watched(self, watched: bool) -> List[Movie]:
        return [m for m in self.movies if m.watched == watched]

    def filter_by_genre(self, genre: str) -> List[Movie]:
        return [m for m in self.movies if m.genre.lower() == genre.lower()]

    def filter_by_year(self, year: int) -> List[Movie]:
        return [m for m in self.movies if m.year == year]

    def get_stats(self) -> dict:
        total = len(self.movies)
        watched_count = len(self.filter_by_watched(True))
        unwatched = total - watched_count
        watched_ratings = [m.rating for m in self.movies if m.watched]
        avg_rating = sum(watched_ratings) / len(watched_ratings) if watched_ratings else 0.0
        genres = {}
        for m in self.movies:
            genres[m.genre] = genres.get(m.genre, 0) + 1
        return {
            "total": total,
            "watched": watched_count,
            "unwatched": unwatched,
            "avg_rating": avg_rating,
            "genres": genres
        }

    def save_to_file(self, filename: str = "movies_data.json") -> None:
        data = {"movies": [asdict(m) for m in self.movies]}
        with open(filename, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=2)

    def load_from_file(self, filename: str = "movies_data.json") -> None:
        path = Path(filename)
        if not path.exists():
            return
        with open(filename, "r", encoding="utf-8") as f:
            data = json.load(f)
            self.movies.clear()
            for item in data.get("movies", []):
                movie = Movie(
                    id=item["id"],
                    title=item["title"],
                    director=item["director"],
                    year=item["year"],
                    genre=item["genre"],
                    watched=item["watched"],
                    rating=item["rating"],
                    added_date=item["added_date"]
                )
                self.movies.append(movie)
                if movie.id >= self.next_id:
                    self.next_id = movie.id + 1

def print_movie(movie: Movie) -> None:
    status = "✅ Просмотрен" if movie.watched else "⏳ Не просмотрен"
    print(f"#{movie.id} - {movie.title} ({movie.year})")
    print(f"   Режиссёр: {movie.director}, Жанр: {movie.genre}")
    print(f"   {status}, Оценка: {movie.rating}/10")
    print(f"   Добавлен: {movie.added_date}")

def main():
    catalog = MovieCatalog()
    catalog.load_from_file()

    while True:
        print("\n===== КАТАЛОГ ФИЛЬМОВ (Python) =====")
        print("1. Добавить фильм")
        print("2. Показать все фильмы")
        print("3. Показать просмотренные фильмы")
        print("4. Показать непросмотренные фильмы")
        print("5. Найти фильмы по названию или режиссёру")
        print("6. Отметить фильм как просмотренный")
        print("7. Редактировать фильм")
        print("8. Удалить фильм")
        print("9. Показать статистику")
        print("10. Сохранить в файл")
        print("11. Загрузить из файла")
        print("0. Выход")
        choice = input("Выберите действие: ").strip()

        if choice == "0":
            break
        elif choice == "1":
            title = input("Название: ").strip()
            if not title:
                print("Название не может быть пустым.")
                continue
            director = input("Режиссёр: ").strip()
            if not director:
                print("Режиссёр не может быть пустым.")
                continue
            try:
                year = int(input("Год выпуска: ").strip())
            except ValueError:
                print("Введите число.")
                continue
            genre = input("Жанр: ").strip()
            if not genre:
                genre = "Неизвестно"
            watched_input = input("Статус (1-просмотрен, 0-не просмотрен): ").strip()
            watched = watched_input == "1"
            try:
                rating = int(input("Оценка (1-10): ").strip())
            except ValueError:
                rating = 0
            try:
                movie = catalog.add_movie(title, director, year, genre, watched, rating)
                print(f"Фильм добавлен с ID {movie.id}")
            except Exception as e:
                print("Ошибка:", e)
        elif choice == "2":
            if not catalog.movies:
                print("Нет фильмов.")
            else:
                for m in catalog.movies:
                    print_movie(m)
        elif choice == "3":
            movies = catalog.filter_by_watched(True)
            if not movies:
                print("Нет просмотренных фильмов.")
            else:
                for m in movies:
                    print_movie(m)
        elif choice == "4":
            movies = catalog.filter_by_watched(False)
            if not movies:
                print("Нет непросмотренных фильмов.")
            else:
                for m in movies:
                    print_movie(m)
        elif choice == "5":
            query = input("Введите часть названия или режиссёра: ").strip()
            if not query:
                print("Введите текст.")
                continue
            results = catalog.search_movies(query)
            if not results:
                print("Фильмы не найдены.")
            else:
                for m in results:
                    print_movie(m)
        elif choice == "6":
            try:
                mid = int(input("Введите ID фильма: ").strip())
            except ValueError:
                print("Некорректный ID.")
                continue
            movie = catalog.find_movie(mid)
            if not movie:
                print("Фильм не найден.")
                continue
            movie.watched = True
            print("Фильм отмечен как просмотренный.")
        elif choice == "7":
            try:
                mid = int(input("Введите ID фильма для редактирования: ").strip())
            except ValueError:
                print("Некорректный ID.")
                continue
            movie = catalog.find_movie(mid)
            if not movie:
                print("Фильм не найден.")
                continue
            print("Оставьте поле пустым, чтобы не менять.")
            new_title = input(f"Название ({movie.title}): ").strip()
            new_director = input(f"Режиссёр ({movie.director}): ").strip()
            new_year = input(f"Год ({movie.year}): ").strip()
            new_genre = input(f"Жанр ({movie.genre}): ").strip()
            new_watched = input(f"Статус (1-просмотрен, 0-не просмотрен) сейчас: {'1' if movie.watched else '0'}: ").strip()
            new_rating = input(f"Оценка ({movie.rating}): ").strip()
            updates = {}
            if new_title: updates["title"] = new_title
            if new_director: updates["director"] = new_director
            if new_year:
                try:
                    updates["year"] = int(new_year)
                except ValueError:
                    print("Год должен быть числом, пропускаем.")
            if new_genre: updates["genre"] = new_genre
            if new_watched: updates["watched"] = new_watched == "1"
            if new_rating:
                try:
                    updates["rating"] = int(new_rating)
                except ValueError:
                    print("Оценка должна быть числом, пропускаем.")
            if catalog.edit_movie(mid, **updates):
                print("Фильм обновлён.")
            else:
                print("Ошибка обновления.")
        elif choice == "8":
            try:
                mid = int(input("Введите ID фильма для удаления: ").strip())
            except ValueError:
                print("Некорректный ID.")
                continue
            if catalog.delete_movie(mid):
                print("Фильм удалён.")
            else:
                print("Фильм не найден.")
        elif choice == "9":
            stats = catalog.get_stats()
            print("\n=== СТАТИСТИКА ===")
            print(f"Всего фильмов: {stats['total']}")
            print(f"Просмотрено: {stats['watched']}")
            print(f"Не просмотрено: {stats['unwatched']}")
            print(f"Средняя оценка (только просмотренные): {stats['avg_rating']:.2f}")
            print("По жанрам:")
            for genre, count in stats['genres'].items():
                print(f"  {genre}: {count}")
        elif choice == "10":
            catalog.save_to_file()
            print("Сохранено.")
        elif choice == "11":
            catalog.load_from_file()
            print("Загружено.")
        else:
            print("Неизвестная команда.")

if __name__ == "__main__":
    main()
