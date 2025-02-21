package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"httpServer/internal/app"
	"httpServer/internal/database"
	"net/http"
)

type Task struct {
	ID         int
	Name       string `json:"name"`
	Details    string `json:"detail"`
	Completion bool   `json:"completion"`
}

func AddTask(appState *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	task := Task{}
	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	_, err = appState.DB.CreateTask(r.Context(), database.CreateTaskParams{
		Title:  task.Name,
		Detail: sql.NullString{String: task.Details, Valid: task.Details != ""},
	})
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
