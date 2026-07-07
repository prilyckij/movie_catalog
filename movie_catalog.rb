# movie_catalog.rb
require 'json'
require 'date'

class Movie
  attr_accessor :id, :title, :director, :year, :genre, :watched, :rating, :added_date

  def initialize(id, title, director, year, genre, watched, rating, added_date = Date.today.to_s)
    @id = id
    @title = title
    @director = director
    @year = year
    @genre = genre
    @watched = watched
    @rating = rating
    @added_date = added_date
  end

  def to_h
    { id: @id, title: @title, director: @director, year: @year, genre: @genre,
      watched: @watched, rating: @rating, added_date: @added_date }
  end

  def self.from_h(hash)
    Movie.new(hash[:id], hash[:title], hash[:director], hash[:year],
              hash[:genre], hash[:watched], hash[:rating], hash[:added_date])
  end
end

class MovieCatalog
  attr_reader :movies

  def initialize
    @movies = []
    @next_id = 1
  end

  def add_movie(title, director, year, genre, watched, rating)
    raise "Оценка должна быть от 1 до 10" unless (1..10).include?(rating)
    raise "Год должен быть от 1900 до #{Date.today.year}" unless (1900..Date.today.year).include?(year)
    raise "Название и режиссёр не могут быть пустыми" if title.empty? || director.empty?
    genre = "Неизвестно" if genre.empty?
    movie = Movie.new(@next_id, title, director, year, genre, watched, rating)
    @movies << movie
    @next_id += 1
    movie
  end

  def find_movie(id)
    @movies.find { |m| m.id == id }
  end

  def edit_movie(id, **kwargs)
    movie = find_movie(id)
    return false unless movie
    kwargs.each do |key, value|
      movie.send("#{key}=", value) if movie.respond_to?("#{key}=")
    end
    true
  end

  def delete_movie(id)
    movie = find_movie(id)
    return false unless movie
    @movies.delete(movie)
    true
  end

  def search_movies(query)
    q = query.downcase
    @movies.select { |m| m.title.downcase.include?(q) || m.director.downcase.include?(q) }
  end

  def filter_by_watched(watched)
    @movies.select { |m| m.watched == watched }
  end

  def filter_by_genre(genre)
    @movies.select { |m| m.genre.downcase == genre.downcase }
  end

  def filter_by_year(year)
    @movies.select { |m| m.year == year }
  end

  def stats
    total = @movies.size
    watched_count = filter_by_watched(true).size
    unwatched = total - watched_count
    watched_ratings = filter_by_watched(true).map(&:rating)
    avg_rating = watched_ratings.empty? ? 0 : watched_ratings.sum.to_f / watched_ratings.size
    genres = Hash.new(0)
    @movies.each { |m| genres[m.genre] += 1 }
    { total: total, watched: watched_count, unwatched: unwatched, avg_rating: avg_rating, genres: genres }
  end

  def save_to_file(filename = "movies_data.json")
    data = { movies: @movies.map(&:to_h) }
    File.write(filename, JSON.pretty_generate(data))
  end

  def load_from_file(filename = "movies_data.json")
    return unless File.exist?(filename)
    data = JSON.parse(File.read(filename), symbolize_names: true)
    @movies.clear
    data[:movies].each do |item|
      movie = Movie.from_h(item)
      @movies << movie
      @next_id = movie.id + 1 if movie.id >= @next_id
    end
  rescue JSON::ParserError
    puts "Ошибка чтения файла."
  end
end

def print_movie(movie)
  status = movie.watched ? "✅ Просмотрен" : "⏳ Не просмотрен"
  puts "##{movie.id} - #{movie.title} (#{movie.year})"
  puts "   Режиссёр: #{movie.director}, Жанр: #{movie.genre}"
  puts "   #{status}, Оценка: #{movie.rating}/10"
  puts "   Добавлен: #{movie.added_date}"
end

def main
  catalog = MovieCatalog.new
  catalog.load_from_file

  loop do
    puts "\n===== КАТАЛОГ ФИЛЬМОВ (Ruby) ====="
    puts "1. Добавить фильм"
    puts "2. Показать все фильмы"
    puts "3. Показать просмотренные фильмы"
    puts "4. Показать непросмотренные фильмы"
    puts "5. Найти фильмы по названию или режиссёру"
    puts "6. Отметить фильм как просмотренный"
    puts "7. Редактировать фильм"
    puts "8. Удалить фильм"
    puts "9. Показать статистику"
    puts "10. Сохранить в файл"
    puts "11. Загрузить из файла"
    puts "0. Выход"
    print "Выберите действие: "
    choice = gets.chomp

    case choice
    when "0"
      break
    when "1"
      print "Название: "
      title = gets.chomp
      next if title.empty?
      print "Режиссёр: "
      director = gets.chomp
      next if director.empty?
      print "Год выпуска: "
      year = gets.chomp.to_i
      print "Жанр: "
      genre = gets.chomp
      print "Статус (1-просмотрен, 0-не просмотрен): "
      watched = gets.chomp == "1"
      print "Оценка (1-10): "
      rating = gets.chomp.to_i
      begin
        movie = catalog.add_movie(title, director, year, genre, watched, rating)
        puts "Фильм добавлен с ID #{movie.id}"
      rescue => e
        puts "Ошибка: #{e.message}"
      end
    when "2"
      if catalog.movies.empty?
        puts "Нет фильмов."
      else
        catalog.movies.each { |m| print_movie(m) }
      end
    when "3"
      movies = catalog.filter_by_watched(true)
      if movies.empty?
        puts "Нет просмотренных фильмов."
      else
        movies.each { |m| print_movie(m) }
      end
    when "4"
      movies = catalog.filter_by_watched(false)
      if movies.empty?
        puts "Нет непросмотренных фильмов."
      else
        movies.each { |m| print_movie(m) }
      end
    when "5"
      print "Введите часть названия или режиссёра: "
      query = gets.chomp
      results = catalog.search_movies(query)
      if results.empty?
        puts "Фильмы не найдены."
      else
        results.each { |m| print_movie(m) }
      end
    when "6"
      print "Введите ID фильма: "
      id = gets.chomp.to_i
      movie = catalog.find_movie(id)
      unless movie
        puts "Фильм не найден."
        next
      end
      movie.watched = true
      puts "Фильм отмечен как просмотренный."
    when "7"
      print "Введите ID фильма для редактирования: "
      id = gets.chomp.to_i
      movie = catalog.find_movie(id)
      unless movie
        puts "Фильм не найден."
        next
      end
      puts "Оставьте поле пустым, чтобы не менять."
      print "Название (#{movie.title}): "
      new_title = gets.chomp
      print "Режиссёр (#{movie.director}): "
      new_director = gets.chomp
      print "Год (#{movie.year}): "
      new_year = gets.chomp
      print "Жанр (#{movie.genre}): "
      new_genre = gets.chomp
      print "Статус (1-просмотрен, 0-не просмотрен) сейчас: #{movie.watched ? '1' : '0'}: "
      new_watched = gets.chomp
      print "Оценка (#{movie.rating}): "
      new_rating = gets.chomp
      updates = {}
      updates[:title] = new_title unless new_title.empty?
      updates[:director] = new_director unless new_director.empty?
      unless new_year.empty?
        updates[:year] = new_year.to_i
      end
      updates[:genre] = new_genre unless new_genre.empty?
      unless new_watched.empty?
        updates[:watched] = new_watched == "1"
      end
      unless new_rating.empty?
        updates[:rating] = new_rating.to_i
      end
      if catalog.edit_movie(id, **updates)
        puts "Фильм обновлён."
      else
        puts "Ошибка обновления."
      end
    when "8"
      print "Введите ID фильма для удаления: "
      id = gets.chomp.to_i
      if catalog.delete_movie(id)
        puts "Фильм удалён."
      else
        puts "Фильм не найден."
      end
    when "9"
      stats = catalog.stats
      puts "\n=== СТАТИСТИКА ==="
      puts "Всего фильмов: #{stats[:total]}"
      puts "Просмотрено: #{stats[:watched]}"
      puts "Не просмотрено: #{stats[:unwatched]}"
      puts "Средняя оценка (только просмотренные): #{'%.2f' % stats[:avg_rating]}"
      puts "По жанрам:"
      stats[:genres].each { |g, c| puts "  #{g}: #{c}" }
    when "10"
      catalog.save_to_file
      puts "Сохранено."
    when "11"
      catalog.load_from_file
      puts "Загружено."
    else
      puts "Неизвестная команда."
    end
  end
end

main if __FILE__ == $0
