package main

import (
	"database/sql"
	"fmt"
	"httpServer/internal/app"
	"httpServer/internal/config"
	"httpServer/internal/database"
	"httpServer/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Function for initializing the Database
func statDB(DBURL string) (*sql.DB, error) {
	DB, err := sql.Open("postgres", DBURL)
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func main() {
	router := mux.NewRouter()

	configMain, err := config.ConfigRead()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	DB, err := statDB(configMain.DbUrl)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer DB.Close()

	queries := database.New(DB)
	appState := &app.AppState{
		AppConfig: configMain,
		DB:        queries,
	}
	router.HandleFunc("/", wrapHandler(appState, handlers.StartPage)).Methods("GET")
	router.HandleFunc("/addTask", wrapHandler(appState, handlers.AddTask)).Methods("POST")
	router.HandleFunc("/listTasks", wrapHandler(appState, handlers.ListTasks)).Methods("GET")
	router.HandleFunc("/deleteTask/{task_id}", wrapHandler(appState, handlers.DeleteTask)).Methods("DELETE")
	router.HandleFunc("/completeTask/{task_id}", wrapHandler(appState, handlers.UpdateTask)).Methods("PUT")
	fmt.Printf("%s\n", configMain.CurrentUserName)
	fmt.Println("Server is running on PORT :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("%s", err)
	}

}

// Wrap function to pass appstate along with r and w
func wrapHandler(appState *app.AppState, handlerFunc func(*app.AppState, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(appState, w, r)
	}
}
