// movie_catalog.cpp
#include <iostream>
#include <vector>
#include <string>
#include <fstream>
#include <sstream>
#include <algorithm>
#include <iomanip>
#include <ctime>
#include <map>
#include <variant>
#include <regex>
#include <cctype>

using namespace std;

struct Movie {
    int id;
    string title;
    string director;
    int year;
    string genre;
    bool watched;
    int rating;
    string addedDate;

    Movie(int id, const string& title, const string& director, int year, const string& genre,
          bool watched, int rating, const string& addedDate = "")
        : id(id), title(title), director(director), year(year), genre(genre),
          watched(watched), rating(rating), addedDate(addedDate) {
        if (addedDate.empty()) {
            time_t now = time(nullptr);
            tm* tm_now = localtime(&now);
            char buf[11];
            strftime(buf, sizeof(buf), "%Y-%m-%d", tm_now);
            this->addedDate = string(buf);
        }
    }
};

class MovieCatalog {
private:
    vector<Movie> movies;
    int nextId = 1;

    int currentYear() {
        time_t now = time(nullptr);
        tm* tm_now = localtime(&now);
        return tm_now->tm_year + 1900;
    }

public:
    Movie addMovie(const string& title, const string& director, int year, const string& genre, bool watched, int rating) {
        if (rating < 1 || rating > 10) throw invalid_argument("Оценка должна быть от 1 до 10");
        int cy = currentYear();
        if (year < 1900 || year > cy) throw invalid_argument("Год должен быть от 1900 до " + to_string(cy));
        if (title.empty() || director.empty()) throw invalid_argument("Название и режиссёр не могут быть пустыми");
        string g = genre.empty() ? "Неизвестно" : genre;
        Movie movie(nextId, title, director, year, g, watched, rating);
        movies.push_back(movie);
        nextId++;
        return movie;
    }

    Movie* findMovie(int id) {
        auto it = find_if(movies.begin(), movies.end(), [id](const Movie& m) { return m.id == id; });
        return it != movies.end() ? &(*it) : nullptr;
    }

    bool editMovie(int id, const map<string, string>& updates) {
        Movie* movie = findMovie(id);
        if (!movie) return false;
        for (const auto& [key, value] : updates) {
            if (key == "title") movie->title = value;
            else if (key == "director") movie->director = value;
            else if (key == "year") movie->year = stoi(value);
            else if (key == "genre") movie->genre = value;
            else if (key == "watched") movie->watched = (value == "1");
            else if (key == "rating") movie->rating = stoi(value);
        }
        return true;
    }

    bool deleteMovie(int id) {
        auto it = find_if(movies.begin(), movies.end(), [id](const Movie& m) { return m.id == id; });
        if (it == movies.end()) return false;
        movies.erase(it);
        return true;
    }

    vector<Movie> searchMovies(const string& query) {
        string q = query;
        transform(q.begin(), q.end(), q.begin(), ::tolower);
        vector<Movie> result;
        for (const auto& m : movies) {
            string titleLower = m.title, directorLower = m.director;
            transform(titleLower.begin(), titleLower.end(), titleLower.begin(), ::tolower);
            transform(directorLower.begin(), directorLower.end(), directorLower.begin(), ::tolower);
            if (titleLower.find(q) != string::npos || directorLower.find(q) != string::npos) {
                result.push_back(m);
            }
        }
        return result;
    }

    vector<Movie> filterByWatched(bool watched) const {
        vector<Movie> result;
        for (const auto& m : movies) {
            if (m.watched == watched) result.push_back(m);
        }
        return result;
    }

    vector<Movie> filterByGenre(const string& genre) const {
        vector<Movie> result;
        for (const auto& m : movies) {
            if (m.genre == genre) result.push_back(m);
        }
        return result;
    }

    vector<Movie> filterByYear(int year) const {
        vector<Movie> result;
        for (const auto& m : movies) {
            if (m.year == year) result.push_back(m);
        }
        return result;
    }

    map<string, variant<int, double, map<string, int>>> getStats() const {
        int total = movies.size();
        int watchedCount = filterByWatched(true).size();
        int unwatched = total - watchedCount;
        int sumRating = 0;
        int ratingCount = 0;
        for (const auto& m : movies) {
            if (m.watched) { sumRating += m.rating; ratingCount++; }
        }
        double avgRating = ratingCount > 0 ? static_cast<double>(sumRating) / ratingCount : 0.0;
        map<string, int> genres;
        for (const auto& m : movies) {
            genres[m.genre]++;
        }
        map<string, variant<int, double, map<string, int>>> stats;
        stats["total"] = total;
        stats["watched"] = watchedCount;
        stats["unwatched"] = unwatched;
        stats["avg_rating"] = avgRating;
        stats["genres"] = genres;
        return stats;
    }

    void saveToFile(const string& filename = "movies_data.txt") {
        ofstream out(filename);
        if (!out) return;
        for (const auto& m : movies) {
            out << m.id << '|'
                << m.title << '|'
                << m.director << '|'
                << m.year << '|'
                << m.genre << '|'
                << m.watched << '|'
                << m.rating << '|'
                << m.addedDate << '\n';
        }
    }

    void loadFromFile(const string& filename = "movies_data.txt") {
        ifstream in(filename);
        if (!in) return;
        movies.clear();
        string line;
        while (getline(in, line)) {
            stringstream ss(line);
            string idStr, title, director, yearStr, genre, watchedStr, ratingStr, addedDate;
            getline(ss, idStr, '|');
            getline(ss, title, '|');
            getline(ss, director, '|');
            getline(ss, yearStr, '|');
            getline(ss, genre, '|');
            getline(ss, watchedStr, '|');
            getline(ss, ratingStr, '|');
            getline(ss, addedDate, '|');
            int id = stoi(idStr);
            int year = stoi(yearStr);
            bool watched = (watchedStr == "1");
            int rating = stoi(ratingStr);
            movies.emplace_back(id, title, director, year, genre, watched, rating, addedDate);
            if (id >= nextId) nextId = id + 1;
        }
    }

    const vector<Movie>& getMovies() const { return movies; }
};

string readString(const string& prompt) {
    cout << prompt;
    string input;
    getline(cin, input);
    return input;
}

int readInt(const string& prompt) {
    while (true) {
        cout << prompt;
        string input;
        getline(cin, input);
        try {
            return stoi(input);
        } catch (...) {
            cout << "Введите число.\n";
        }
    }
}

bool readBool(const string& prompt) {
    while (true) {
        string input = readString(prompt);
        if (input == "1") return true;
        if (input == "0") return false;
        cout << "Введите 1 или 0.\n";
    }
}

void printMovie(const Movie& movie) {
    string status = movie.watched ? "✅ Просмотрен" : "⏳ Не просмотрен";
    cout << "#" << movie.id << " - " << movie.title << " (" << movie.year << ")\n";
    cout << "   Режиссёр: " << movie.director << ", Жанр: " << movie.genre << "\n";
    cout << "   " << status << ", Оценка: " << movie.rating << "/10\n";
    cout << "   Добавлен: " << movie.addedDate << "\n";
}

int main() {
    MovieCatalog catalog;
    catalog.loadFromFile();

    while (true) {
        cout << "\n===== КАТАЛОГ ФИЛЬМОВ (C++) =====" << endl;
        cout << "1. Добавить фильм\n";
        cout << "2. Показать все фильмы\n";
        cout << "3. Показать просмотренные фильмы\n";
        cout << "4. Показать непросмотренные фильмы\n";
        cout << "5. Найти фильмы по названию или режиссёру\n";
        cout << "6. Отметить фильм как просмотренный\n";
        cout << "7. Редактировать фильм\n";
        cout << "8. Удалить фильм\n";
        cout << "9. Показать статистику\n";
        cout << "10. Сохранить в файл\n";
        cout << "11. Загрузить из файла\n";
        cout << "0. Выход\n";
        string choice = readString("Выберите действие: ");

        if (choice == "0") break;

        if (choice == "1") {
            string title = readString("Название: ");
            if (title.empty()) { cout << "Название не может быть пустым.\n"; continue; }
            string director = readString("Режиссёр: ");
            if (director.empty()) { cout << "Режиссёр не может быть пустым.\n"; continue; }
            int year = readInt("Год выпуска: ");
            string genre = readString("Жанр: ");
            if (genre.empty()) genre = "Неизвестно";
            bool watched = readBool("Статус (1-просмотрен, 0-не просмотрен): ");
            int rating = readInt("Оценка (1-10): ");
            try {
                Movie movie = catalog.addMovie(title, director, year, genre, watched, rating);
                cout << "Фильм добавлен с ID " << movie.id << "\n";
            } catch (const exception& e) {
                cout << "Ошибка: " << e.what() << "\n";
            }
        } else if (choice == "2") {
            if (catalog.getMovies().empty()) {
                cout << "Нет фильмов.\n";
            } else {
                for (const auto& m : catalog.getMovies()) printMovie(m);
            }
        } else if (choice == "3") {
            auto movies = catalog.filterByWatched(true);
            if (movies.empty()) cout << "Нет просмотренных фильмов.\n";
            else for (const auto& m : movies) printMovie(m);
        } else if (choice == "4") {
            auto movies = catalog.filterByWatched(false);
            if (movies.empty()) cout << "Нет непросмотренных фильмов.\n";
            else for (const auto& m : movies) printMovie(m);
        } else if (choice == "5") {
            string query = readString("Введите часть названия или режиссёра: ");
            auto results = catalog.searchMovies(query);
            if (results.empty()) cout << "Фильмы не найдены.\n";
            else for (const auto& m : results) printMovie(m);
        } else if (choice == "6") {
            int id = readInt("Введите ID фильма: ");
            Movie* movie = catalog.findMovie(id);
            if (!movie) { cout << "Фильм не найден.\n"; continue; }
            movie->watched = true;
            cout << "Фильм отмечен как просмотренный.\n";
        } else if (choice == "7") {
            int id = readInt("Введите ID фильма для редактирования: ");
            Movie* movie = catalog.findMovie(id);
            if (!movie) { cout << "Фильм не найден.\n"; continue; }
            cout << "Оставьте поле пустым, чтобы не менять.\n";
            string newTitle = readString("Название (" + movie->title + "): ");
            string newDirector = readString("Режиссёр (" + movie->director + "): ");
            string newYear = readString("Год (" + to_string(movie->year) + "): ");
            string newGenre = readString("Жанр (" + movie->genre + "): ");
            string newWatched = readString("Статус (1-просмотрен, 0-не просмотрен) сейчас: " + string(movie->watched ? "1" : "0") + ": ");
            string newRating = readString("Оценка (" + to_string(movie->rating) + "): ");
            map<string, string> updates;
            if (!newTitle.empty()) updates["title"] = newTitle;
            if (!newDirector.empty()) updates["director"] = newDirector;
            if (!newYear.empty()) updates["year"] = newYear;
            if (!newGenre.empty()) updates["genre"] = newGenre;
            if (!newWatched.empty()) updates["watched"] = newWatched;
            if (!newRating.empty()) updates["rating"] = newRating;
            if (catalog.editMovie(id, updates)) {
                cout << "Фильм обновлён.\n";
            } else {
                cout << "Ошибка обновления.\n";
            }
        } else if (choice == "8") {
            int id = readInt("Введите ID фильма для удаления: ");
            if (catalog.deleteMovie(id)) {
                cout << "Фильм удалён.\n";
            } else {
                cout << "Фильм не найден.\n";
            }
        } else if (choice == "9") {
            auto stats = catalog.getStats();
            cout << "\n=== СТАТИСТИКА ===\n";
            cout << "Всего фильмов: " << get<int>(stats["total"]) << "\n";
            cout << "Просмотрено: " << get<int>(stats["watched"]) << "\n";
            cout << "Не просмотрено: " << get<int>(stats["unwatched"]) << "\n";
            cout << "Средняя оценка (только просмотренные): " << fixed << setprecision(2) << get<double>(stats["avg_rating"]) << "\n";
            cout << "По жанрам:\n";
            auto genres = get<map<string, int>>(stats["genres"]);
            for (const auto& [g, c] : genres) {
                cout << "  " << g << ": " << c << "\n";
            }
        } else if (choice == "10") {
            catalog.saveToFile();
            cout << "Сохранено.\n";
        } else if (choice == "11") {
            catalog.loadFromFile();
            cout << "Загружено.\n";
        } else {
            cout << "Неизвестная команда.\n";
        }
    }
    return 0;
}
