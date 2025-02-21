package handlers

import (
	"httpServer/internal/app"
	"net/http"
)

// Just testing function
func StartPage(appState *app.AppState, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
