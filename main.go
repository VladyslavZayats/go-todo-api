package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Task struct {
	id   int
	Name string
}

var tasks []Task

func handler(write http.ResponseWriter, read *http.Request) {
	fmt.Println(*read)
	fmt.Fprint(write, "Hello from Task API\n")
}

func listTasks(write http.ResponseWriter, read *http.Request) {
	write.Header().Set("Content-Type", "application/json")
	json.NewEncoder(write).Encode(tasks)
	write.WriteHeader(http.StatusOK)
}

func addTask(write http.ResponseWriter, read *http.Request) {
	/*
		if read.Method != http.MethodPost {
			http.Error(write, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
	*/

	var task Task

	/*
		err := json.NewDecoder(read.Body).Decode(&task)

		if err != nil {
			http.Error(write, "Error decoding JSON", http.StatusBadRequest)
			return
		}
	*/

	var taskName = read.URL.Query().Get("name")

	if taskName == "" {
		http.Error(write, "Invalid task name", http.StatusBadRequest)
		return
	}

	task.Name = taskName

	if len(tasks) == 0 {
		task.id = 0
	} else {
		task.id = tasks[len(tasks)-1].id + 1
	}

	tasks = append(tasks, task)
	fmt.Println(tasks)
	write.WriteHeader(http.StatusOK)
	write.Write([]byte("Task added sucessfully\n"))
}

func deleteTask(write http.ResponseWriter, read *http.Request) {
	if len(tasks) == 0 {
		http.Error(write, "Nothing to delete", http.StatusBadRequest)
		return
	}

	var taskId, err = strconv.Atoi(read.URL.Query().Get("id"))

	if err != nil {
		http.Error(write, "Invalid task id", http.StatusBadRequest)
		return
	}

	var found bool = false

	fmt.Printf("taskId: %v, type: %T\n", taskId, taskId)

	for i, task := range tasks {
		if task.id == taskId {
			fmt.Println(i)
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(write, "Task not found", http.StatusNotFound)
		return
	}

	write.WriteHeader(http.StatusOK)
	write.Write([]byte("Task deleted sucessfully\n"))
}

func main() {
	fmt.Println("Listening...")

	http.HandleFunc("/", handler)
	http.HandleFunc("/list", listTasks)
	http.HandleFunc("/add", addTask)
	http.HandleFunc("/delete", deleteTask)

	http.ListenAndServe(":8080", nil)

}
