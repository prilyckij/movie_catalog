// MovieCatalog.java
import java.io.*;
import java.nio.file.*;
import java.time.LocalDate;
import java.util.*;
import java.util.stream.Collectors;

record Movie(int id, String title, String director, int year, String genre, boolean watched, int rating, String addedDate) implements Serializable {}

class MoviesData implements Serializable {
    private static final long serialVersionUID = 1L;
    public List<Movie> movies;
}

class MovieCatalog implements Serializable {
    private static final long serialVersionUID = 1L;
    private List<Movie> movies = new ArrayList<>();
    private int nextId = 1;

    public Movie addMovie(String title, String director, int year, String genre, boolean watched, int rating) {
        if (rating < 1 || rating > 10) throw new IllegalArgumentException("Оценка должна быть от 1 до 10");
        if (year < 1900 || year > LocalDate.now().getYear())
            throw new IllegalArgumentException("Год должен быть от 1900 до " + LocalDate.now().getYear());
        if (title.isBlank() || director.isBlank())
            throw new IllegalArgumentException("Название и режиссёр не могут быть пустыми");
        if (genre.isBlank()) genre = "Неизвестно";
        Movie movie = new Movie(nextId, title, director, year, genre, watched, rating, LocalDate.now().toString());
        movies.add(movie);
        nextId++;
        return movie;
    }

    public Optional<Movie> findMovie(int id) {
        return movies.stream().filter(m -> m.id() == id).findFirst();
    }

    public boolean editMovie(int id, Map<String, Object> updates) {
        Optional<Movie> opt = findMovie(id);
        if (opt.isEmpty()) return false;
        Movie old = opt.get();
        movies.remove(old);
        String title = (String) updates.getOrDefault("title", old.title());
        String director = (String) updates.getOrDefault("director", old.director());
        int year = (int) updates.getOrDefault("year", old.year());
        String genre = (String) updates.getOrDefault("genre", old.genre());
        boolean watched = (boolean) updates.getOrDefault("watched", old.watched());
        int rating = (int) updates.getOrDefault("rating", old.rating());
        Movie updated = new Movie(old.id(), title, director, year, genre, watched, rating, old.addedDate());
        movies.add(updated);
        return true;
    }

    public boolean deleteMovie(int id) {
        return movies.removeIf(m -> m.id() == id);
    }

    public List<Movie> searchMovies(String query) {
        String q = query.toLowerCase();
        return movies.stream()
                .filter(m -> m.title().toLowerCase().contains(q) || m.director().toLowerCase().contains(q))
                .collect(Collectors.toList());
    }

    public List<Movie> filterByWatched(boolean watched) {
        return movies.stream().filter(m -> m.watched() == watched).collect(Collectors.toList());
    }

    public List<Movie> filterByGenre(String genre) {
        return movies.stream().filter(m -> m.genre().equalsIgnoreCase(genre)).collect(Collectors.toList());
    }

    public List<Movie> filterByYear(int year) {
        return movies.stream().filter(m -> m.year() == year).collect(Collectors.toList());
    }

    public Map<String, Object> getStats() {
        int total = movies.size();
        int watchedCount = filterByWatched(true).size();
        int unwatched = total - watchedCount;
        double avgRating = filterByWatched(true).stream().mapToInt(Movie::rating).average().orElse(0);
        Map<String, Integer> genres = new HashMap<>();
        movies.forEach(m -> genres.put(m.genre(), genres.getOrDefault(m.genre(), 0) + 1));
        Map<String, Object> stats = new HashMap<>();
        stats.put("total", total);
        stats.put("watched", watchedCount);
        stats.put("unwatched", unwatched);
        stats.put("avg_rating", avgRating);
        stats.put("genres", genres);
        return stats;
    }

    public void saveToFile(String filename) throws IOException {
        MoviesData data = new MoviesData();
        data.movies = new ArrayList<>(movies);
        try (ObjectOutputStream oos = new ObjectOutputStream(Files.newOutputStream(Paths.get(filename)))) {
            oos.writeObject(data);
        }
    }

    public void loadFromFile(String filename) throws IOException, ClassNotFoundException {
        try (ObjectInputStream ois = new ObjectInputStream(Files.newInputStream(Paths.get(filename)))) {
            MoviesData data = (MoviesData) ois.readObject();
            movies = new ArrayList<>(data.movies);
            for (Movie m : movies) {
                if (m.id() >= nextId) nextId = m.id() + 1;
            }
        }
    }

    public List<Movie> getMovies() { return Collections.unmodifiableList(movies); }
}

public class MovieCatalogApp {
    private static final Scanner scanner = new Scanner(System.in);

    private static String readString(String prompt) {
        System.out.print(prompt);
        return scanner.nextLine().trim();
    }

    private static int readInt(String prompt) {
        while (true) {
            try {
                System.out.print(prompt);
                return Integer.parseInt(scanner.nextLine().trim());
            } catch (NumberFormatException e) {
                System.out.println("Введите число.");
            }
        }
    }

    private static boolean readBool(String prompt) {
        while (true) {
            String input = readString(prompt);
            if (input.equals("1")) return true;
            if (input.equals("0")) return false;
            System.out.println("Введите 1 или 0.");
        }
    }

    private static void printMovie(Movie movie) {
        String status = movie.watched() ? "✅ Просмотрен" : "⏳ Не просмотрен";
        System.out.printf("#%d - %s (%d)%n", movie.id(), movie.title(), movie.year());
        System.out.printf("   Режиссёр: %s, Жанр: %s%n", movie.director(), movie.genre());
        System.out.printf("   %s, Оценка: %d/10%n", status, movie.rating());
        System.out.printf("   Добавлен: %s%n", movie.addedDate());
    }

    public static void main(String[] args) {
        MovieCatalog catalog = new MovieCatalog();
        try {
            catalog.loadFromFile("movies_data.ser");
        } catch (IOException | ClassNotFoundException e) {
            System.out.println("Не удалось загрузить данные.");
        }

        while (true) {
            System.out.println("\n===== КАТАЛОГ ФИЛЬМОВ (Java) =====");
            System.out.println("1. Добавить фильм");
            System.out.println("2. Показать все фильмы");
            System.out.println("3. Показать просмотренные фильмы");
            System.out.println("4. Показать непросмотренные фильмы");
            System.out.println("5. Найти фильмы по названию или режиссёру");
            System.out.println("6. Отметить фильм как просмотренный");
            System.out.println("7. Редактировать фильм");
            System.out.println("8. Удалить фильм");
            System.out.println("9. Показать статистику");
            System.out.println("10. Сохранить в файл");
            System.out.println("11. Загрузить из файла");
            System.out.println("0. Выход");
            String choice = readString("Выберите действие: ");

            switch (choice) {
                case "0" -> { return; }
                case "1" -> {
                    String title = readString("Название: ");
                    if (title.isBlank()) { System.out.println("Название не может быть пустым."); continue; }
                    String director = readString("Режиссёр: ");
                    if (director.isBlank()) { System.out.println("Режиссёр не может быть пустым."); continue; }
                    int year = readInt("Год выпуска: ");
                    String genre = readString("Жанр: ");
                    if (genre.isBlank()) genre = "Неизвестно";
                    boolean watched = readBool("Статус (1-просмотрен, 0-не просмотрен): ");
                    int rating = readInt("Оценка (1-10): ");
                    try {
                        Movie movie = catalog.addMovie(title, director, year, genre, watched, rating);
                        System.out.println("Фильм добавлен с ID " + movie.id());
                    } catch (Exception e) {
                        System.out.println("Ошибка: " + e.getMessage());
                    }
                }
                case "2" -> {
                    if (catalog.getMovies().isEmpty()) System.out.println("Нет фильмов.");
                    else catalog.getMovies().forEach(MovieCatalogApp::printMovie);
                }
                case "3" -> {
                    var movies = catalog.filterByWatched(true);
                    if (movies.isEmpty()) System.out.println("Нет просмотренных фильмов.");
                    else movies.forEach(MovieCatalogApp::printMovie);
                }
                case "4" -> {
                    var movies = catalog.filterByWatched(false);
                    if (movies.isEmpty()) System.out.println("Нет непросмотренных фильмов.");
                    else movies.forEach(MovieCatalogApp::printMovie);
                }
                case "5" -> {
                    String query = readString("Введите часть названия или режиссёра: ");
                    var results = catalog.searchMovies(query);
                    if (results.isEmpty()) System.out.println("Фильмы не найдены.");
                    else results.forEach(MovieCatalogApp::printMovie);
                }
                case "6" -> {
                    int id = readInt("Введите ID фильма: ");
                    var opt = catalog.findMovie(id);
                    if (opt.isEmpty()) { System.out.println("Фильм не найден."); continue; }
                    Movie old = opt.get();
                    catalog.editMovie(id, Map.of("watched", true));
                    System.out.println("Фильм отмечен как просмотренный.");
                }
                case "7" -> {
                    int id = readInt("Введите ID фильма для редактирования: ");
                    var opt = catalog.findMovie(id);
                    if (opt.isEmpty()) { System.out.println("Фильм не найден."); continue; }
                    Movie old = opt.get();
                    System.out.println("Оставьте поле пустым, чтобы не менять.");
                    String newTitle = readString("Название (" + old.title() + "): ");
                    String newDirector = readString("Режиссёр (" + old.director() + "): ");
                    String newYearStr = readString("Год (" + old.year() + "): ");
                    String newGenre = readString("Жанр (" + old.genre() + "): ");
                    String newWatchedStr = readString("Статус (1-просмотрен, 0-не просмотрен) сейчас: " + (old.watched() ? "1" : "0") + ": ");
                    String newRatingStr = readString("Оценка (" + old.rating() + "): ");
                    Map<String, Object> updates = new HashMap<>();
                    if (!newTitle.isBlank()) updates.put("title", newTitle);
                    if (!newDirector.isBlank()) updates.put("director", newDirector);
                    if (!newYearStr.isBlank()) {
                        try { updates.put("year", Integer.parseInt(newYearStr)); }
                        catch (NumberFormatException e) { System.out.println("Год должен быть числом, пропускаем."); }
                    }
                    if (!newGenre.isBlank()) updates.put("genre", newGenre);
                    if (!newWatchedStr.isBlank()) updates.put("watched", newWatchedStr.equals("1"));
                    if (!newRatingStr.isBlank()) {
                        try { updates.put("rating", Integer.parseInt(newRatingStr)); }
                        catch (NumberFormatException e) { System.out.println("Оценка должна быть числом, пропускаем."); }
                    }
                    if (catalog.editMovie(id, updates)) System.out.println("Фильм обновлён.");
                    else System.out.println("Ошибка обновления.");
                }
                case "8" -> {
                    int id = readInt("Введите ID фильма для удаления: ");
                    if (catalog.deleteMovie(id)) System.out.println("Фильм удалён.");
                    else System.out.println("Фильм не найден.");
                }
                case "9" -> {
                    var stats = catalog.getStats();
                    System.out.println("\n=== СТАТИСТИКА ===");
                    System.out.println("Всего фильмов: " + stats.get("total"));
                    System.out.println("Просмотрено: " + stats.get("watched"));
                    System.out.println("Не просмотрено: " + stats.get("unwatched"));
                    System.out.printf("Средняя оценка (только просмотренные): %.2f%n", stats.get("avg_rating"));
                    System.out.println("По жанрам:");
                    @SuppressWarnings("unchecked")
                    Map<String, Integer> genres = (Map<String, Integer>) stats.get("genres");
                    genres.forEach((g, c) -> System.out.println("  " + g + ": " + c));
                }
                case "10" -> {
                    try {
                        catalog.saveToFile("movies_data.ser");
                        System.out.println("Сохранено.");
                    } catch (IOException e) {
                        System.out.println("Ошибка сохранения: " + e.getMessage());
                    }
                }
                case "11" -> {
                    try {
                        catalog.loadFromFile("movies_data.ser");
                        System.out.println("Загружено.");
                    } catch (IOException | ClassNotFoundException e) {
                        System.out.println("Ошибка загрузки: " + e.getMessage());
                    }
                }
                default -> System.out.println("Неизвестная команда.");
            }
        }
    }
}
