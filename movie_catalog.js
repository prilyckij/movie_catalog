// movie_catalog.js
const fs = require('fs').promises;
const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const question = (prompt) => new Promise(resolve => rl.question(prompt, resolve));

class Movie {
    constructor(id, title, director, year, genre, watched, rating, addedDate) {
        this.id = id;
        this.title = title;
        this.director = director;
        this.year = year;
        this.genre = genre;
        this.watched = watched;
        this.rating = rating;
        this.addedDate = addedDate || new Date().toISOString().slice(0, 10);
    }
}

class MovieCatalog {
    constructor() {
        this.movies = [];
        this.nextId = 1;
    }

    addMovie(title, director, year, genre, watched, rating) {
        if (rating < 1 || rating > 10) throw new Error('Оценка должна быть от 1 до 10');
        const currentYear = new Date().getFullYear();
        if (year < 1900 || year > currentYear) throw new Error(`Год должен быть от 1900 до ${currentYear}`);
        if (!title.trim() || !director.trim()) throw new Error('Название и режиссёр не могут быть пустыми');
        if (!genre.trim()) genre = 'Неизвестно';
        const movie = new Movie(this.nextId, title, director, year, genre, watched, rating);
        this.movies.push(movie);
        this.nextId++;
        return movie;
    }

    findMovie(id) {
        return this.movies.find(m => m.id === id);
    }

    editMovie(id, updates) {
        const movie = this.findMovie(id);
        if (!movie) return false;
        Object.assign(movie, updates);
        return true;
    }

    deleteMovie(id) {
        const index = this.movies.findIndex(m => m.id === id);
        if (index === -1) return false;
        this.movies.splice(index, 1);
        return true;
    }

    searchMovies(query) {
        const q = query.toLowerCase();
        return this.movies.filter(m => m.title.toLowerCase().includes(q) || m.director.toLowerCase().includes(q));
    }

    filterByWatched(watched) {
        return this.movies.filter(m => m.watched === watched);
    }

    filterByGenre(genre) {
        return this.movies.filter(m => m.genre.toLowerCase() === genre.toLowerCase());
    }

    filterByYear(year) {
        return this.movies.filter(m => m.year === year);
    }

    getStats() {
        const total = this.movies.length;
        const watchedCount = this.filterByWatched(true).length;
        const unwatched = total - watchedCount;
        const ratings = this.filterByWatched(true).map(m => m.rating);
        const avgRating = ratings.length ? ratings.reduce((a, b) => a + b, 0) / ratings.length : 0;
        const genres = {};
        this.movies.forEach(m => { genres[m.genre] = (genres[m.genre] || 0) + 1; });
        return { total, watched: watchedCount, unwatched, avgRating, genres };
    }

    async saveToFile(filename = 'movies_data.json') {
        const data = { movies: this.movies };
        await fs.writeFile(filename, JSON.stringify(data, null, 2));
    }

    async loadFromFile(filename = 'movies_data.json') {
        try {
            const data = await fs.readFile(filename, 'utf8');
            const parsed = JSON.parse(data);
            this.movies = parsed.movies.map(m => Object.assign(new Movie(0), m));
            this.nextId = this.movies.reduce((max, m) => Math.max(max, m.id), 0) + 1;
        } catch (err) {
            if (err.code !== 'ENOENT') throw err;
        }
    }
}

function printMovie(movie) {
    const status = movie.watched ? '✅ Просмотрен' : '⏳ Не просмотрен';
    console.log(`#${movie.id} - ${movie.title} (${movie.year})`);
    console.log(`   Режиссёр: ${movie.director}, Жанр: ${movie.genre}`);
    console.log(`   ${status}, Оценка: ${movie.rating}/10`);
    console.log(`   Добавлен: ${movie.addedDate}`);
}

async function main() {
    const catalog = new MovieCatalog();
    await catalog.loadFromFile();

    while (true) {
        console.log('\n===== КАТАЛОГ ФИЛЬМОВ (JavaScript) =====');
        console.log('1. Добавить фильм');
        console.log('2. Показать все фильмы');
        console.log('3. Показать просмотренные фильмы');
        console.log('4. Показать непросмотренные фильмы');
        console.log('5. Найти фильмы по названию или режиссёру');
        console.log('6. Отметить фильм как просмотренный');
        console.log('7. Редактировать фильм');
        console.log('8. Удалить фильм');
        console.log('9. Показать статистику');
        console.log('10. Сохранить в файл');
        console.log('11. Загрузить из файла');
        console.log('0. Выход');
        const choice = await question('Выберите действие: ');

        if (choice === '0') break;

        switch (choice) {
            case '1': {
                const title = await question('Название: ');
                if (!title.trim()) { console.log('Название не может быть пустым.'); continue; }
                const director = await question('Режиссёр: ');
                if (!director.trim()) { console.log('Режиссёр не может быть пустым.'); continue; }
                const year = parseInt(await question('Год выпуска: '));
                let genre = await question('Жанр: ');
                if (!genre.trim()) genre = 'Неизвестно';
                const watched = (await question('Статус (1-просмотрен, 0-не просмотрен): ')) === '1';
                const rating = parseInt(await question('Оценка (1-10): '));
                try {
                    const movie = catalog.addMovie(title, director, year, genre, watched, rating);
                    console.log(`Фильм добавлен с ID ${movie.id}`);
                } catch (err) {
                    console.log('Ошибка:', err.message);
                }
                break;
            }
            case '2':
                if (catalog.movies.length === 0) console.log('Нет фильмов.');
                else catalog.movies.forEach(printMovie);
                break;
            case '3': {
                const movies = catalog.filterByWatched(true);
                if (movies.length === 0) console.log('Нет просмотренных фильмов.');
                else movies.forEach(printMovie);
                break;
            }
            case '4': {
                const movies = catalog.filterByWatched(false);
                if (movies.length === 0) console.log('Нет непросмотренных фильмов.');
                else movies.forEach(printMovie);
                break;
            }
            case '5': {
                const query = await question('Введите часть названия или режиссёра: ');
                const results = catalog.searchMovies(query);
                if (results.length === 0) console.log('Фильмы не найдены.');
                else results.forEach(printMovie);
                break;
            }
            case '6': {
                const id = parseInt(await question('Введите ID фильма: '));
                const movie = catalog.findMovie(id);
                if (!movie) { console.log('Фильм не найден.'); continue; }
                catalog.editMovie(id, { watched: true });
                console.log('Фильм отмечен как просмотренный.');
                break;
            }
            case '7': {
                const id = parseInt(await question('Введите ID фильма для редактирования: '));
                const movie = catalog.findMovie(id);
                if (!movie) { console.log('Фильм не найден.'); continue; }
                console.log('Оставьте поле пустым, чтобы не менять.');
                const newTitle = await question(`Название (${movie.title}): `);
                const newDirector = await question(`Режиссёр (${movie.director}): `);
                const newYear = await question(`Год (${movie.year}): `);
                const newGenre = await question(`Жанр (${movie.genre}): `);
                const newWatched = await question(`Статус (1-просмотрен, 0-не просмотрен) сейчас: ${movie.watched ? '1' : '0'}: `);
                const newRating = await question(`Оценка (${movie.rating}): `);
                const updates = {};
                if (newTitle.trim()) updates.title = newTitle;
                if (newDirector.trim()) updates.director = newDirector;
                if (newYear.trim()) {
                    const y = parseInt(newYear);
                    if (!isNaN(y)) updates.year = y;
                    else console.log('Год должен быть числом, пропускаем.');
                }
                if (newGenre.trim()) updates.genre = newGenre;
                if (newWatched.trim()) updates.watched = newWatched === '1';
                if (newRating.trim()) {
                    const r = parseInt(newRating);
                    if (!isNaN(r)) updates.rating = r;
                    else console.log('Оценка должна быть числом, пропускаем.');
                }
                if (catalog.editMovie(id, updates)) console.log('Фильм обновлён.');
                else console.log('Ошибка обновления.');
                break;
            }
            case '8': {
                const id = parseInt(await question('Введите ID фильма для удаления: '));
                if (catalog.deleteMovie(id)) console.log('Фильм удалён.');
                else console.log('Фильм не найден.');
                break;
            }
            case '9': {
                const stats = catalog.getStats();
                console.log('\n=== СТАТИСТИКА ===');
                console.log(`Всего фильмов: ${stats.total}`);
                console.log(`Просмотрено: ${stats.watched}`);
                console.log(`Не просмотрено: ${stats.unwatched}`);
                console.log(`Средняя оценка (только просмотренные): ${stats.avgRating.toFixed(2)}`);
                console.log('По жанрам:');
                for (const [genre, count] of Object.entries(stats.genres)) {
                    console.log(`  ${genre}: ${count}`);
                }
                break;
            }
            case '10':
                try {
                    await catalog.saveToFile();
                    console.log('Сохранено.');
                } catch (err) {
                    console.log('Ошибка сохранения:', err.message);
                }
                break;
            case '11':
                try {
                    await catalog.loadFromFile();
                    console.log('Загружено.');
                } catch (err) {
                    console.log('Ошибка загрузки:', err.message);
                }
                break;
            default:
                console.log('Неизвестная команда.');
        }
    }
    rl.close();
}

main().catch(console.error);
