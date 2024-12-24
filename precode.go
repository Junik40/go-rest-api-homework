package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...
func GetTasks(w http.ResponseWriter, r *http.Request) {
	marshOut,err := json.Marshal(tasks)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshOut)
}

func GetTasksId (w http.ResponseWriter, r *http.Request){
	urrl:=r.URL.String()
	str:=strings.Split(urrl, "/")
	id:=str[len(str)-1]
	task,ok:=tasks[id]
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	marshOut,err := json.Marshal(task)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshOut)
}

func DeleteTaskId (w http.ResponseWriter, r *http.Request){
	urrl:=r.URL.String()
	str:=strings.Split(urrl, "/")
	id:=str[len(str)-1]
	_,ok:=tasks[id]
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	delete(tasks,id)
	w.WriteHeader(http.StatusOK)
}


func PostTasks(w http.ResponseWriter, r *http.Request){
	var task Task
    var buf bytes.Buffer

    _, err := buf.ReadFrom(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    tasks[task.ID] = task

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}

func main() {
	r := chi.NewRouter()
	
	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", GetTasks)
	r.Post("/tasks", PostTasks)
	r.Get("/tasks/{id}", GetTasksId)
	r.Delete("/tasks/{id}", DeleteTaskId)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
