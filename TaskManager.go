package main

import (
    "bufio"
    "fmt"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "sync"
    "syscall"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/layout"
)

type Task struct {
    id          int
    description string
    completed   bool
}

type TaskManager struct {
    mu     sync.Mutex
    tasks  []Task
    nextID int
}

func NewTaskManager() *TaskManager {
    return &TaskManager{}
}

func (tm *TaskManager) AddTask(description string) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    tm.nextID++
    task := Task{id: tm.nextID, description: description}
    tm.tasks = append(tm.tasks, task)
    fmt.Println("Задача добавлена.")
    go func() {
        if err := tm.SaveToFile("tasks.txt"); err != nil {
            fmt.Println("Ошибка сохранения задач в файл:", err)
        }
    }()
}

func (tm *TaskManager) DeleteTask(id int) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    for i, task := range tm.tasks {
        if task.id == id {
            tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
            fmt.Println("Задача удалена.")
            go func() {
                if err := tm.SaveToFile("tasks.txt"); err != nil {
                    fmt.Println("Ошибка сохранения задач в файл:", err)
                }
            }()
            return
        }
    }
    fmt.Println("Задача не найдена.")
}

func (tm *TaskManager) CompleteTask(id int) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    for i, task := range tm.tasks {
        if task.id == id {
            tm.tasks[i].completed = true
            fmt.Println("Задача отмечена как выполненная.")
            go func() {
                if err := tm.SaveToFile("tasks.txt"); err != nil {
                    fmt.Println("Ошибка сохранения задач в файл:", err)
                }
            }()
            return
        }
    }
    fmt.Println("Задача не найдена.")
}

func (tm *TaskManager) ShowTasks() string {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    if len(tm.tasks) == 0 {
        return "Список задач пуст."
    }

    var sb strings.Builder
    sb.WriteString("Список задач:\n")
    for _, task := range tm.tasks {
        status := "не выполнена"
        if task.completed {
            status = "выполнена"
        }
        sb.WriteString(fmt.Sprintf("ID: %d, Описание: %s, Статус: %s\n", task.id, task.description, status))
    }
    return sb.String()
}

func (tm *TaskManager) SaveToFile(filename string) error {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    for _, task := range tm.tasks {
        completed := 0
        if task.completed {
            completed = 1
        }
        if _, err := writer.WriteString(fmt.Sprintf("%d;%s;%d\n", task.id, task.description, completed)); err != nil {
            return err
        }
    }

    return nil
}

func (tm *TaskManager) LoadFromFile(filename string) error {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    tm.tasks = nil // Clear existing tasks

    var maxID int

    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            continue // Skip empty lines
        }

        parts := strings.Split(line, ";")
        if len(parts) != 3 {
            fmt.Println("Ошибка разбора строки: некорректное количество полей")
            continue // Skip invalid lines
        }

        id, err := strconv.Atoi(parts[0])
        if err != nil {
            fmt.Println("Ошибка разбора строки: неверный формат ID")
            continue
        }

        completed, err := strconv.Atoi(parts[2])
        if err != nil {
            fmt.Println("Ошибка разбора строки: неверный формат статуса выполнения")
            continue
        }

        if id > maxID {
            maxID = id
        }

        tm.tasks = append(tm.tasks, Task{id: id, description: parts[1], completed: completed == 1})
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Ошибка при сканировании файла:", err)
        return err
    }

    tm.nextID = maxID
    fmt.Printf("Загружено %d задач из файла\n", len(tm.tasks))
    return nil
}

func resourceIcon() fyne.Resource {
    ico, err := loadIconFromFile("TaskManager.ico")
    if err != nil {
        fmt.Println("Ошибка загрузки иконки:", err)
        // Возвращаем стандартную иконку в случае ошибки
        return theme.FyneLogo()
    }
    return ico
}

func loadIconFromFile(filename string) (fyne.Resource, error) {
    icon, err := fyne.LoadResourceFromPath(filename)
    if err != nil {
        fmt.Printf("Ошибка загрузки иконки из %s: %v\n", filename, err)
        return nil, err
    }
    return icon, nil
}

func main() {
    version := "2.0.1"

    taskManager := NewTaskManager()

    // Загрузка задач из файла при запуске
    if err := taskManager.LoadFromFile("tasks.txt"); err != nil && !os.IsNotExist(err) {
        fmt.Println("Ошибка загрузки задач из файла:", err)
        return
    } else if os.IsNotExist(err) {
        fmt.Println("Файл tasks.txt не существует, создаем новый файл.")
        if _, err := os.Create("tasks.txt"); err != nil {
            fmt.Println("Ошибка создания файла:", err)
            return
        }
    }

    myApp := app.New()
    myApp.SetIcon(resourceIcon())
    myWindow := myApp.NewWindow("Task Manager")

    defer func() {
        if err := taskManager.SaveToFile("tasks.txt"); err != nil {
            fmt.Println("Ошибка сохранения задач в файл:", err)
        }
        myApp.Quit()
    }()

    descriptionEntry := widget.NewEntry()
    descriptionEntry.SetPlaceHolder("Описание задачи")

    addButton := widget.NewButton("Добавить задачу", func() {
        description := descriptionEntry.Text
        if description != "" {
            taskManager.AddTask(description)
            descriptionEntry.SetText("")
            dialog.ShowInformation("Успех", "Задача добавлена.", myWindow)
        } else {
            dialog.ShowError(fmt.Errorf("Описание не может быть пустым"), myWindow)
        }
    })

    idEntry := widget.NewEntry()
    idEntry.SetPlaceHolder("ID задачи")

    deleteButton := widget.NewButton("Удалить задачу", func() {
        idStr := idEntry.Text
        id, err := strconv.Atoi(idStr)
        if err != nil {
            dialog.ShowError(fmt.Errorf("Неверный ID"), myWindow)
            return
        }
        taskManager.DeleteTask(id)
        idEntry.SetText("")
        dialog.ShowInformation("Успех", "Задача удалена.", myWindow)
    })

    completeButton := widget.NewButton("Отметить задачу как выполненную", func() {
        idStr := idEntry.Text
        id, err := strconv.Atoi(idStr)
        if err != nil {
            dialog.ShowError(fmt.Errorf("Неверный ID"), myWindow)
            return
        }
        taskManager.CompleteTask(id)
        idEntry.SetText("")
        dialog.ShowInformation("Успех", "Задача отмечена как выполненная.", myWindow)
    })

    showButton := widget.NewButton("Показать список задач", func() {
        tasks := taskManager.ShowTasks()
        dialog.ShowInformation("Список задач", tasks, myWindow)
    })

    versionLabel := widget.NewLabel("Version: " + version) // Здесь вы можете указать вашу версию

    content := container.NewVBox(
        widget.NewLabel("Task Manager"),
        descriptionEntry,
        addButton,
        idEntry,
        deleteButton,
        completeButton,
        showButton,
    )

    footer := container.NewHBox(
        layout.NewSpacer(), // Оставляет пространство слева
        versionLabel,       // Метка версии справа
    )

    myWindow.SetContent(container.NewBorder(
        content, // Center content
        nil,     // Top content
        nil,     // Left content
        nil,     // Right content
        footer,  // Bottom content
    ))

    myWindow.Resize(fyne.NewSize(400, 280))
    myWindow.Show()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-stop
        if err := taskManager.SaveToFile("tasks.txt"); err != nil {
            fmt.Println("Ошибка сохранения задач в файл:", err)
        }
        myApp.Quit()
    }()

    myApp.Run()
}