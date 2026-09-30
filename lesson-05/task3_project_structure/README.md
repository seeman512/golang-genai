# Завдання 3 — `task3_project_structure`

Завдання показує, як **монолітний `main.go` перетворити на стандартну багатошарову
структуру Go-проєкту** на прикладі маленької CLI-програми «todo list».

Програма вміє три дії (обидві версії виконують ті самі дії):

- **add item** — додати завдання: `add <title>`
- **delete item** — видалити завдання за ID: `delete <id>`
- **complete item** — позначити завдання виконаним: `complete <id>`

Плюс службові команди `list`, `help`, `quit` (або `exit`). Дані зберігаються
лише в пам'яті й зникають після завершення процесу.

---

## Структура каталогу

```
task3_project_structure/
├── README.md                       ← цей файл
├── cmd/
│   ├── bad_structure/
│   │   └── main.go                 ← «погана» версія: вся програма в одному файлі
│   └── good_structure/
│       └── main.go                 ← «добра» версія: тонкий main.go (лише wiring)
└── internal/                       ← шари доброї версії
    ├── handler/
    │   ├── handler.go              ← CLI-інтерфейс: читання команд, вивід результатів
    │   └── handler_test.go         ← тест команди add/complete/list/delete
    ├── service/
    │   ├── todo_service.go         ← бізнес-логіка та валідація вхідних даних
    │   └── todo_service_test.go    ← тести помилок сервісу
    ├── repository/
    │   ├── todo_repository.go      ← інтерфейс сховища + sentinel ErrNotFound
    │   └── memory.go               ← in-memory реалізація сховища
    └── model/
        └── todo.go                 ← доменна структура Todo
```

---

## Погана версія: все в одному файлі

`cmd/bad_structure/main.go` — типовий початковий варіант навчальної програми.
У ньому одночасно живуть **усі шари**:

- модель `todoItem` (доменні дані);
- «сховище» — зріз `todos` і лічильник `nextID` як локальні змінні в `run()`;
- бізнес-логіка — пошук `findTodo`, видалення через `append(todos[:i], todos[i+1:]...)`;
- валідація вводу — `readID` (парсинг ID через `strconv.Atoi`);
- presentation — `switch` по командах, `fmt.Print("> ")`, `printHelp()`;
- wiring — усе це запускається з `main()`.

### Чому це погано

| Проблема | Наслідок |
|---|---|
| Усе в одному пакунку `main` | Жодну частину логіки неможливо імпортувати, а отже — **неможливо протестувати** окремо від CLI |
| Стан (`todos`, `nextID`) — локальні змінні функції | Логіку неможливо підмінити фейком; тести потребують запуску цілої програми та парсингу stdout |
| Presentation змішана з логікою | Зміна виводу (наприклад, JSON замість тексту) вимагає правити бізнес-код |
| Немає меж відповідальності | Додавання нового сховища (файл, БД) означає переписування тієї самої функції |
| Помилки перевіряються «на місці» й одразу друкуються | Неможливо відрізнити «не знайдено» від «невалідний ID» програмно |

Це не «неправильний» код — він працює. Але він **не масштабується**: вже на
третьому-четвертому типі сутності такий файл перетворюється на кілька сотень
рядків, де важко щось знайти.

---

## Добра версія: багатошарова структура

Той самий функціонал розкладено за областю відповідальності (by domain) у
чотири пакунки під `internal/`, плюс тонкий `main.go`.

### Напрям залежностей

```
cmd/good_structure/main.go  (wiring)
        │
        ▼
  internal/handler      ← CLI: читає команди, друкує результат
        │
        ▼
  internal/service      ← бізнес-логіка, валідація, обгортання помилок
        │
        ▼
  internal/repository   ← інтерфейс TodoRepository + in-memory реалізація
        │
        ▼
  internal/model        ← структура Todo (не імпортує нікого)
```

Правило одностороннє: **handler → service → repository → model**.
`model` не залежить ні від кого, `handler` нічого не знає про сховище.
Зворотних стрілок (model → service, repository → handler) у коді немає.

### Відповідальність шарів

- **`model`** — доменна структура `Todo{ID, Title, Completed}`. Жодної логіки.
- **`repository`** — інтерфейс `TodoRepository` (`Add`, `Delete`, `Complete`,
  `List`) і sentinel-помилка `ErrNotFound`. `MemoryTodoRepository` — єдина
  реалізація; саме тут живе стан (зріз елементів і `nextID`).
- **`service`** — прикладні сценарії: обрізає/валідує заголовок
  (`ErrEmptyTitle`), делегує операції сховищу й **обгортає помилки з `%w`**
  (`fmt.Errorf("delete todo %d: %w", id, err)`), щоб виклик міг перевірити
  причину через `errors.Is`.
- **`handler`** — лише presentation: `bufio.Scanner` по командах, `switch`,
  форматування виводу. Приймає `io.Reader`/`io.Writer`, тому **тестується без
  консолі** (див. `handler_test.go` з `strings.Reader` і `bytes.Buffer`).
- **`cmd/good_structure/main.go`** — тонкий: створює сховище, передає його в
  сервіс, сервіс — у handler, запускає `Run()` і звітує про помилку в
  `os.Stderr`. У ньому немає бізнес-логіки взагалі.

### Що це дає

- Кожен шар тестується окремо; заміна in-memory сховища на файл чи БД не
  зачіпає `service` і `handler` (достатньо нової реалізації `TodoRepository`).
- Помилки мають контекст і ланцюг причин (`errors.Is(err, repository.ErrNotFound)`).
- Тести не залежать від stdin/stdout і працюють швидко.
- Код легко читати: кожен файл має одну причину змінюватися.

---

## Як запустити

Усі команди виконуються з каталогу `lesson-05/` (там лежить `go.mod`
з модулем `lesson05`).

### Погана (однофайлова) версія

```bash
go run ./task3_project_structure/cmd/bad_structure
```

### Добра (багатошарова) версія

```bash
go run ./task3_project_structure/cmd/good_structure
```

### Приклад сесії

```text
Todo list (items are kept in memory for this run).
Commands: add <title>, delete <id>, complete <id>, list, help, quit
> add Buy groceries
Added #1: Buy groceries
> add Read a book
Added #2: Read a book
> complete 1
Completed #1
> list
[x] #1 Buy groceries
[ ] #2 Read a book
> delete 2
Deleted #2
> list
[x] #1 Buy groceries
> quit
```

Неінтерактивна перевірка (зручно для скриптів і демо):

```bash
printf 'add Buy groceries\ncomplete 1\nlist\ndelete 1\nquit\n' \
  | go run ./task3_project_structure/cmd/good_structure
```

### Тести та перевірки

```bash
# тести task3
go test ./task3_project_structure/... -v

# перевірка форматування
gofmt -l ./task3_project_structure
```

---

## Порада для рефакторингу

Порядок розкладання моноліта на шари:

1. Винести **доменні типи** в `model` — вони ні від чого не залежать.
2. Винести **роботу зі станом** у `repository`; одразу оголосити інтерфейс за
   потребами споживача, а не за можливостями реалізації.
3. Винести **правила та валідацію** в `service`; там само обгортати помилки
   контекстом операції через `%w`.
4. Залишити в `handler` тільки **ввід/вивід**; приймати `io.Reader`/`io.Writer`,
   щоб код був тестованим.
5. У `main.go` лишити **wiring** — створення залежностей і їх передачу.

Файл `cmd/bad_structure/main.go` навмисно залишено без змін — це еталон
«до», з яким порівнюється результат рефакторингу.
