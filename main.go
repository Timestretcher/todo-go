package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

func saveTasks(tasks []Task, filename string) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	// ⬇️ Вот здесь создаём файл
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() // Файл закрывается после записи

	_, err = file.Write(data) // Пишем JSON в файл
	return err
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	var tasks []Task // список задач

	idCounter := 1 // счётчик ID

	for {
		fmt.Print("Введите новую задачу(или напишите 'выход'): ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)

		//проверка на выход
		if strings.ToLower(title) == "выход" {
			break
		}

		//создаётся задача
		task := Task{
			ID:    idCounter,
			Title: title,
			Done:  false,
		}

		//добавление в список задач
		tasks = append(tasks, task)
		idCounter++
	}

	//вывод задач
	printTasks := func(tasks []Task) {
		fmt.Println("\nСписок задач:")
		for _, t := range tasks {
			status := "✗"
			if t.Done {
				status = "✓"
			}
			fmt.Printf("%d. [%-1s] %s\n", t.ID, status, t.Title)
		}
	}

	//Сохранение задач в JSON file

	//первоначальный вывод
	printTasks(tasks)

	err := saveTasks(tasks, "tasks.json")
	if err != nil {
		fmt.Println("Ошибка при сохранении задач:", err)
	} else {
		fmt.Println("Задачи успешно сохранены в файл.")
	}

	//отметка выполненных задач
	for {
		fmt.Print("\nВведите номер выполненной задачи (или '0' для выхода без изменений)")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Ошибка: введите номер существующей задачи")
			continue

		}

		if num == 0 {
			fmt.Println("Выход из режима отметки задач")
			break
		}

		//поиск и обновление
		found := false

		for i, t := range tasks {
			if t.ID == num {
				tasks[i].Done = true
				fmt.Printf("Задача %d отмечена как выполеннная.\n", num)

				saveErr := saveTasks(tasks, "tasks.json")
				if saveErr != nil {
					fmt.Println("Ошибка при сохранении задач:", saveErr)
				} else {
					fmt.Println("Изменения сохранены в файл.")
				}
				found = true
				break
			}

		}

		if found {
			printTasks(tasks)
		} else {
			fmt.Println("Ошибка: введите номер существующей задачи")
		}
	}
}
