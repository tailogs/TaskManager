# TaskManager

![image](https://github.com/tailogs/TaskManager/assets/69743960/0ad7d17f-0587-434d-ac06-6cad49cc9c90)

*Рисунок 1. Интерфейс моей программы*

---

Вы можете скачать бинарную версию приложения из репозитория релизов по этой [ссылке](https://github.com/tailogs/TaskManager/releases) или скомпилировать самому из исходных кодов.

## Возможности

- Добавление, удаление и отметка задач как выполненных.
- Визуализация списка задач с указанием статуса выполнения.
- Сохранение задач в файл и загрузка из файла.

## Установка и запуск

### Предварительные требования

- Установлен Go (версия 1.16 или выше). Скачайте и установите с [официального сайта Go](https://golang.org/dl/).
- Установлен компилятор C (для библиотеки Fyne требуется CGo).

### Используя нативные средства `Go`

1. Склонируйте проект:

    ```sh
    git clone https://github.com/tailogs/TaskManager.git
    ```

    ```sh
    cd TaskManager
    ```

2. Инициализируйте модуль Go:

    ```sh
    go mod init TaskManager
    ```

3. Установите библиотеки для создания графического интерфейса и другие зависимости:

    ```sh
    go get fyne.io/fyne/v2
    ```
    
    ```sh
    go get github.com/Knetic/govaluate
    ```

4. Установите инструмент командной строки для сборки иконок и ресурсов:

    ```sh
    go get fyne.io/fyne/v2/cmd/fyne
    ```

5. Установите инструмент для внедрения ресурсов:

    ```sh
    go install github.com/akavel/rsrc@latest
    ```

6. Создайте файл ресурсов с иконкой:

    ```sh
    rsrc -ico TaskManager.ico -o rsrc.syso
    ```

7. Очистите зависимости:

    ```sh
    go mod tidy
    ```

8. Постройте проект:

    ```sh
    go build -ldflags="-H=windowsgui" -o TaskManager.exe .
    ```

9. Запустите проект:

    ```sh
    TaskManager.exe
    ```

### Использование

- Запустите приложение. Откроется окно менеджера задач.
- Добавляйте новые задачи, удаляйте или отмечайте их как выполненные.
- Воспользуйтесь кнопкой "Показать список задач", чтобы просмотреть текущий список задач.

---

> Работает только в системе `WINDOWS`, так как это моя основная система, но вы можете сделать форк этого проекта и помочь мне сделать поддержку `Linux` и других ОС.

## Лицензия

Этот проект лицензирован под лицензией MIT. Подробности смотрите в файле LICENSE.
