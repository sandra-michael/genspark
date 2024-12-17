package handlers

import (
	"Assignment1/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am working and healthy"))
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	//id := insertUser(context.Background(), db, "Jane", "<EMAIL>", 24)

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request")
		return
	}
	var newTask models.NewTask

	err = json.Unmarshal(reqBody, &newTask)
	if err != nil {
		fmt.Println("unmarshall error", err)
		http.Error(w, "Error while unmarshal", http.StatusExpectationFailed)
		return
	}

	err = h.validate.Struct(newTask)
	if err != nil {
		fmt.Println("validation failed error", err)
		http.Error(w, "Error while validation"+err.Error(), http.StatusExpectationFailed)
		return
	}

	task, err := h.c.CreateTask(r.Context(), newTask)
	if err != nil {
		fmt.Println("Error while recieving", err)
		http.Error(w, "Error while recieving task", http.StatusExpectationFailed)
		return
	}

	res, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}

	w.Write(res)
}

func (h *Handler) fetchTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println("Error fetching id from path", err)
		http.Error(w, "Error fetching id from path", http.StatusExpectationFailed)
		return
	}

	task, err := h.c.FetchTask(r.Context(), id)
	if err != nil {
		fmt.Println("Error fetching task", err)
		http.Error(w, "Error while recieving task", http.StatusExpectationFailed)
		return
	}

	jsondata, err := json.Marshal(task)
	if err != nil {
		fmt.Println("could not marshal", err)
		http.Error(w, "Error while rmasrshalling", http.StatusExpectationFailed)

		return
	}
	w.Write(jsondata)

}

func (h *Handler) fetchTasks(w http.ResponseWriter, r *http.Request) {
	tasksList, err := h.c.FetchTasks(r.Context())
	if err != nil {
		fmt.Println("Error fetching tasks", err)
		http.Error(w, "Error while recieving tasks", http.StatusExpectationFailed)
		return
	}
	jsondata, err := json.Marshal(tasksList)
	if err != nil {
		fmt.Println("could not marshal", err)
		http.Error(w, "Error while rmasrshalling", http.StatusExpectationFailed)

		return
	}
	w.Write(jsondata)
}

func (h *Handler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		fmt.Println("Error fetching id from path", err)
		http.Error(w, "Error fetching id from path", http.StatusExpectationFailed)
		return
	}

	err = h.c.UpdateTaskStatus(r.Context(), id)
	if err != nil {
		fmt.Println("Error fetching tasks", err)
		http.Error(w, err.Error(), http.StatusExpectationFailed)
		return
	}

	w.Write([]byte("Updated Task "))
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		fmt.Println("Error fetching id from path", err)
		http.Error(w, "Error fetching id from path", http.StatusExpectationFailed)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request")
		return
	}
	var updateTask models.UpdateTask

	err = json.Unmarshal(reqBody, &updateTask)
	if err != nil {
		fmt.Println("unmarshall error", err)
		http.Error(w, "Error while unmarshal", http.StatusExpectationFailed)
		return
	}

	err = h.validate.Struct(updateTask)
	if err != nil {
		fmt.Println("validation failed error", err)
		http.Error(w, "Error while validation", http.StatusExpectationFailed)
		return
	}

	err = h.c.UpdateTask(r.Context(), id, updateTask)
	if err != nil {
		fmt.Println("Error fetching tasks", err)
		http.Error(w, err.Error(), http.StatusExpectationFailed)
		return
	}

	w.Write([]byte("Updated Task "))
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println("Error fetching id from path", err)
		http.Error(w, "Error fetching id from path", http.StatusExpectationFailed)
		return
	}

	err = h.c.DeleteTask(r.Context(), id)
	if err != nil {
		fmt.Println("Error deleting task", err)
		http.Error(w, "Error while deleting task", http.StatusExpectationFailed)
		return
	}

	w.Write([]byte("Deleted data"))

}
