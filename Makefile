# Команда по умолчанию
default: build

# Имя основного бинарного файла
BINARY_NAME=TaskManager.exe

# Переменные для путей
SRC_DIR=.

# Сборка проекта
build:
	go build -o $(BINARY_NAME) $(SRC_DIR)

# Очистка
clean:
	go clean

# Запуск приложения
run:
	$(BINARY_NAME)

# Очистка, сборка и запуск приложения
build-run: clean build run

# Установка зависимостей проекта
install:
	go mod init TaskManager
	go get fyne.io/fyne/v2
	go get fyne.io/fyne/v2/cmd/fyne
	go get -d ./...
	go install github.com/akavel/rsrc@latest
	rsrc -ico TaskManager.ico -o rsrc.syso
	go mod tidy

# Обозначение фальшивых целей
.PHONY: build clean run install build-run