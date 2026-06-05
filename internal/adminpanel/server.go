package adminpanel

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	//Paths
	rootDir string

	//Templates
	adminTmpl *template.Template

	//Errors
	RootDirErr   = errors.New("Не удалось получить корневой каталог")
	AdminTmplErr = errors.New("Не удалось обработать шаблон admin.html")

	//Admin profile
	user_admin string
	pass_admin string
)

func StartServer() {
	err := getRootDir()
	if err != nil {
		log.Fatalf("%v\n", RootDirErr.Error())
	}

	getAdminProfile()
	log.Printf("%s\n%s\n", user_admin, pass_admin)

	err = getAdminTemplate()
	if err != nil {
		log.Fatalf("%v\n", AdminTmplErr.Error())
	}

	mux := http.NewServeMux()

	mux.Handle(
		"/static/css/",
		http.StripPrefix(
			"/static/css/",
			http.FileServer(http.Dir(filepath.Join(rootDir, "templates", "static", "css")))))

	mux.HandleFunc("/", notFoundHandler)
	mux.HandleFunc("/{$}", redirectToPanel)
	mux.HandleFunc("/admin-panel", panelHandler)
	mux.HandleFunc("/admin-panel/logout", logoutHandler)

	http.ListenAndServe(":8081", mux)
}

func getRootDir() error {
	var err error
	rootDir, err = os.Getwd()
	if err != nil {
		return RootDirErr
	}
	log.Printf("ROOTDIR: %s", rootDir)

	return nil
}

func getAdminProfile() {
	user_admin = os.Getenv("USER_ADMIN")
	pass_admin = os.Getenv("PASSWORD_ADMIN")
}

func getAdminTemplate() error {
	var err error
	adminTmpl, err = template.ParseFiles(filepath.Join(
		rootDir,
		"templates",
		"admin.html"),
	)
	if err != nil {
		return AdminTmplErr
	}

	log.Printf("Шаблон админ-панели обработан!\n")
	return nil
}
