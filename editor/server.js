import express from 'express';
import session from 'express-session';
import passport from 'passport';
import LocalStrategy from 'passport-local';
import bcrypt from 'bcrypt';
import dotenv from 'dotenv';
import crypto from 'crypto'
import fs from 'fs';
import path, {dirname} from 'path';
import {exec} from 'child_process';
import OpenAI from "openai";
import Mustache from 'mustache';
import {fileURLToPath} from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
dotenv.config()

const PORT = process.env.PORT || 3100;
const BIN_ROOT = process.env.BIN_ROOT || path.join(__dirname, 'bin')
const CONTENT_TOOL = path.join(BIN_ROOT, 'content')
const TMP_ROOT = process.env.TMP_ROOT || path.join(__dirname, 'tmp')
const FILES_ROOT = process.env.FILES_ROOT || path.join(process.env.HOME, 'src/archive/app/content');
const PROMPTS_DIR = process.env.PROMPTS_DIR || path.join(__dirname, 'prompts');
const TEMPLATES_DIR = process.env.TEMPLATES_DIR || path.join(__dirname, 'templates');

const AUTH_SESSION_SECRET = process.env.AUTH_SESSION_SECRET || crypto.randomBytes(32).toString('hex');
const AUTH_USER = process.env.AUTH_USER || 'admin'
const AUTH_PASSWORD = process.env.AUTH_PASSWORD || crypto.randomBytes(32).toString('hex');

const NO_PREFIX = process.env.NO_PREFIX || false;

!!!process.env.AUTH_PASSWORD && console.log(`User: ${AUTH_USER}\nPassword: ${AUTH_PASSWORD}`)

const openai = new OpenAI({
    organization: process.env.OPENAI_ORGANIZATION,
    project: process.env.OPENAI_PROJECT,
    apiKey: process.env.OPENAI_API_KEY
});

const app = express();

app.use(express.json({limit: '50mb'}));
app.use(express.urlencoded({limit: '50mb', extended: true}));
NO_PREFIX && app.use((req, res, next) => {
    if (req.path.startsWith('/_editor/')) {
        const newPath = req.path.replace('/_editor', '');
        const fullUrl = newPath + (req.url.includes('?') ? req.url.slice(req.url.indexOf('?')) : '');
        res.redirect(307, fullUrl); // Сохраняет ?tmp=1
    } else {
        next();
    }
});
app.use(express.static('public'));

app.use(session({
    secret: AUTH_SESSION_SECRET, // Замените на свой секретный ключ
    resave: false,
    saveUninitialized: false
}));
app.use(passport.initialize());
app.use(passport.session());

const saltRounds = 10;

const hashedPassword = bcrypt.hashSync(AUTH_PASSWORD, saltRounds);
const users = [{id: 1, username: AUTH_USER, password: hashedPassword}];

passport.use(new LocalStrategy(
    (username, password, done) => {
        const user = users.find(u => u.username === username);
        if (!user || !bcrypt.compareSync(password, user.password)) {
            return done(null, false, {message: 'Неверный логин или пароль'});
        }
        return done(null, user);
    }
));

passport.serializeUser((user, done) => {
    done(null, user.id);
});

passport.deserializeUser((id, done) => {
    const user = users.find(u => u.id === id);
    done(null, user);
});

// Middleware для проверки аутентификации
function isAuthenticated(req, res, next) {
    if (req.isAuthenticated()) {
        return next();
    }
    res.redirect('/login');
}

// Middleware для API
function protectAPI(req, res, next) {
    if (req.isAuthenticated()) {
        return next();
    }
    res.status(401).json({error: 'Unauthorized'});
}

// Маршруты

// Защищенная главная страница
app.get('/', isAuthenticated, (req, res) => {
    res.sendFile(path.join(__dirname, 'protected', 'index.html'));
});

app.get('/index.html', isAuthenticated, (req, res) => {
    res.sendFile(path.join(__dirname, 'protected', 'index.html'));
});

// Страница логина
app.get('/login', (req, res) => {
    if (req.isAuthenticated()) {
        return res.redirect('/'); // Если уже авторизован, редирект на главную
    }
    res.sendFile(path.join(__dirname, 'public', 'login.html'));
});

// Обработка логина
app.post('/login',
    passport.authenticate('local', {
        successRedirect: '/', // Успешный логин -> главная страница
        failureRedirect: '/login' // Неудача -> обратно на логин
    })
);

// Выход
app.get('/logout', (req, res) => {
    req.logout((err) => {
        if (err) return next(err);
        res.redirect('/login');
    });
});

function sortFilesNumerically(files) {
    return files.sort((a, b) => {
        const numA = parseInt(a.name.match(/^\d+/)?.[0]) || 0;
        const numB = parseInt(b.name.match(/^\d+/)?.[0]) || 0;
        return numA - numB;
    });
}

// Рекурсивная функция для получения структуры файлов
function getFileTree(dir = ".") {
    const items = fs.readdirSync(path.join(FILES_ROOT, dir), {withFileTypes: true});
    return sortFilesNumerically(items.map(item => {
        const fullPath = path.join(dir, item.name);
        return item.isDirectory()
            ? {name: item.name, type: 'directory', path: fullPath, children: getFileTree(fullPath)}
            : {name: item.name, type: 'file', path: fullPath};
    }));
}

// Получение структуры файлового дерева
app.get('/api/files', protectAPI, (req, res) => {
    try {
        const fileTree = getFileTree();
        res.json(fileTree);
    } catch (err) {
        console.error(err.message ?? err);
        res.status(500).json({error: 'Ошибка при чтении файлов'});
    }
});

// Запуск скрипта упаковки файла
app.post('/api/pack-file', protectAPI, (req, res) => {
    const {filename, markdown, json, tmp} = req.body;
    if (!filename) {
        return res.status(400).json({error: 'Не указано имя файла'});
    }

    const filePath = tmp ? path.join(TMP_ROOT, filename.replace(/\.json$/, '.preview.json')) : path.join(FILES_ROOT, filename);
    const jsonFile = path.join(TMP_ROOT, filename)
    const markdownFile = jsonFile.replace(/\.json$/, '.md');

    fs.writeFileSync(markdownFile, markdown)
    fs.writeFileSync(jsonFile, json)

    exec(`${CONTENT_TOOL} --pack --item-file="${filePath}" --fields-file="${jsonFile}" --text-file="${markdownFile}"`, (err, stdout, stderr) => {
        if (err) {
            console.error(err.message ?? err);
            return res.status(500).json({error: stderr});
        }
        res.json({message: 'Файл запакован успешно', output: stdout});
    });
});

// Запуск скрипта разделения файла и отправка результата в редактор
app.post('/api/unpack-file', protectAPI, (req, res) => {
    const {filename} = req.body;
    if (!filename) {
        return res.status(400).json({error: 'Не указано имя файла'});
    }

    const filePath = path.join(FILES_ROOT, filename);
    const jsonFile = path.join(TMP_ROOT, filename)
    const markdownFile = jsonFile.replace(/\.json$/, '.md');

    exec(`${CONTENT_TOOL} --unpack --item-file="${filePath}" --fields-file="${jsonFile}" --text-file="${markdownFile}"`, (err, stdout, stderr) => {
        if (err) {
            console.error(err.message ?? err);
            return res.status(500).json({error: stderr});
        }

        const markdownContent = fs.existsSync(markdownFile) ? fs.readFileSync(markdownFile, 'utf8') : '';
        const jsonContent = fs.existsSync(jsonFile) ? fs.readFileSync(jsonFile, 'utf8') : '';

        res.json({markdown: markdownContent, json: jsonContent});
    });
});

async function processCommand(command, text) {
    const promptFile = path.join(PROMPTS_DIR, `${command}.txt`);
    const templateFile = path.join(TEMPLATES_DIR, `${command}.mustache`);

    if (!fs.existsSync(promptFile)) {
        return `Ошибка: команда '${command}' не найдена`;
    }
    if (!fs.existsSync(templateFile)) {
        return `Ошибка: шаблон '${command}' не найден`;
    }

    const promptTemplate = fs.readFileSync(promptFile, 'utf8');
    const fullPrompt = `${promptTemplate}\n\n${text}`;

    try {
        const completion = await openai.chat.completions.create({
            messages: [{role: "user", content: fullPrompt}],
            model: "gpt-4o",
            store: true,
        });

        if (completion.choices && completion.choices.length > 0) {
            const response = completion.choices[0].message.content.trim()
            console.log(response)
            let responseJson = JSON.parse(response);
            if (responseJson.table?.length) {
                responseJson.has_table = true
            }
            const templateContent = fs.readFileSync(templateFile, 'utf8');
            return Mustache.render(templateContent, responseJson);
            ;
        } else {
            return 'Ошибка обработки ответа OpenAI';
        }
    } catch (err) {
        console.error(err.message ?? err);
        return 'Ошибка при запросе к OpenAI';
    }
}

async function processFragments(text) {
    const regex = /<%%(.*?)%%(.*?)%%>/gs;
    let matches = [...text.matchAll(regex)];

    if (matches.length === 0) {
        return text; // Если нет команд, возвращаем оригинальный текст
    }

    let tasks = matches.map(async ([fullMatch, command, fragment]) => {
        return {fullMatch, replacement: await processCommand(command, fragment)};
    });

    let results = await Promise.all(tasks);
    results.forEach(({fullMatch, replacement}) => {
        text = text.replace(fullMatch, replacement);
    });

    return text;
}

// API-метод обработки текста с параллельными запросами
app.post('/api/process-all', protectAPI, async (req, res) => {
    const {text} = req.body;
    if (!text) {
        return res.status(400).json({error: 'Текст не передан'});
    }

    try {
        const processedText = await processFragments(text);
        res.json({result: processedText});
    } catch (err) {
        console.error(err.message ?? err);
        res.status(500).json({error: 'Ошибка при обработке текста'});
    }
});

app.post('/api/process', protectAPI, async (req, res) => {
    const {command, text} = req.body;
    if (!command || !text) {
        return res.status(400).json({error: 'Не указана команда или текст'});
    }
    res.json({result: await processCommand(command, text)});
});

// Эндпойнт для предпросмотра HTML из JSON файла
app.get('/preview/*', protectAPI, (req, res) => {
    // Получаем путь к файлу из URL
    const filename = req.params[0];
    const tmp = req.query.tmp;

    if (!filename) {
        return res.status(400).json({error: 'Имя файла не указано'});
    }

    const filePath = tmp ? path.join(TMP_ROOT, filename.replace(/\.json$/, '.preview.json')) : path.join(FILES_ROOT, filename);
    const templatePath = path.join(TEMPLATES_DIR, 'preview.mustache');

    try {
        // Проверяем существование файлов
        if (!fs.existsSync(filePath)) {
            return res.status(404).json({error: 'Файл не найден'});
        }
        if (!fs.existsSync(templatePath)) {
            return res.status(500).json({error: 'Шаблон не найден'});
        }

        // Читаем содержимое файла и шаблона
        const fileContent = fs.readFileSync(filePath, 'utf8');
        const template = fs.readFileSync(templatePath, 'utf8');

        // Парсим JSON
        const jsonData = JSON.parse(fileContent);

        // Проверяем наличие поля body
        if (!jsonData.body) {
            return res.status(400).json({error: 'Поле body не найдено в JSON'});
        }

        // Рендерим шаблон с данными
        const htmlPage = Mustache.render(template, {
            body: jsonData.body,
            title: jsonData.title,
        });

        // Отправляем HTML
        res.setHeader('Content-Type', 'text/html');
        res.send(htmlPage);
    } catch (err) {
        console.error(err.message ?? err);
        res.status(500).json({error: 'Ошибка при обработке файла', details: error.message});
    }
});

app.listen(PORT, () => {
    console.log(`Сервер запущен на http://localhost:${PORT}`);
});
