package adminpanel

import (
	"context"
	"errors"
	"log"
	"logger/internal/database"
	"logger/internal/models"
	"net/http"
	"path/filepath"
	"strconv"
)

var (
	ExecAdminErr = errors.New("Не удалось заполнить шаблон админ-панели")
	pageConvErr  = errors.New("Ошибка конвертации параметра \"page\"")
	AuthErr      = errors.New("Аутентификация не пройдена!")
)

const limit = 10

type templateData struct {
	LogsTable   []models.LogsTable
	CurrentPage int
	NextPage    int
	PrevPage    int
}

func panelHandler(w http.ResponseWriter, r *http.Request) {
	err := requierAuth(w, r)
	if err != nil {
		return
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	ctx := r.Context()

	logs, err := database.GetLogs(ctx, limit, (page-1)*limit)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("Запрос отменен пользователем!\n")
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := templateData{
		LogsTable:   logs,
		CurrentPage: page,
		NextPage:    page + 1,
		PrevPage:    page - 1,
	}

	err = adminTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, ExecAdminErr.Error(), http.StatusInternalServerError)
		return
	}
}

func requierAuth(w http.ResponseWriter, r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if !ok || username != user_admin || password != pass_admin {
		w.Header().Set("WWW-Authenticate", `Basic realm="Gologit login"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return AuthErr
	}

	return nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Gologit logout"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Вы успешно вышли из системы!"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(rootDir, "templates", "404.html"))
}

func redirectToPanel(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8081/admin-panel", http.StatusFound)
}
