package handlers

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/KusoKaihatsuSha/tray_helper/internal/config"
	"github.com/KusoKaihatsuSha/tray_helper/internal/helpers"
	"github.com/KusoKaihatsuSha/tray_helper/internal/tray"

	"fyne.io/systray"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	_ http.HandlerFunc = errorHandler(0)
	_ http.HandlerFunc = listHandler()
	_ http.HandlerFunc = saveHandler()
	_ http.HandlerFunc = loadHandler()
)

//go:embed web/*
var embedDataTemplate embed.FS

// FileServer was get from chi examples AS-IS.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		helpers.ToLog("FileServer does not permit any URL parameters.")
	}
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		customFs := http.StripPrefix(pathPrefix, http.FileServer(root))
		customFs.ServeHTTP(w, r)
	})
}

// routing is used for 'chi' routing presets.
func routing() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/settings", listHandler())
	r.Post("/save", saveHandler())
	r.Post("/load", loadHandler())
	r.Post("/update", updateHandler())

	serverRoot, err := fs.Sub(embedDataTemplate, "web")
	helpers.ToLog(err)
	fileServer(r, "/"+"files"+"/", http.FS(serverRoot))

	runTest := chi.NewRouteContext()
	runTest.Reset()
	if r.Match(runTest, "GET", "/debug") {
		helpers.ToLog("WARNING!!! route to /debug")
	}

	return r
}

// errorHandler wrap around status code.
func errorHandler(code int) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, http.StatusText(code), code)
	}
}

// listHandler used to display the settings.
func listHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "text/html")
		funcs := make(map[string]interface{})
		myT := template.Must(template.New("main").Funcs(funcs).ParseFS(embedDataTemplate, "web/templates/index.html")).Funcs(funcs)
		err := myT.ExecuteTemplate(rw, "index.html", nil)
		helpers.ToLog(err)
	}
}

// saveHandler used to save the settings.
func saveHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if b, err := io.ReadAll(r.Body); err == nil {
			helpers.ToLog(string(b))
			os.WriteFile(config.Get().Config, b, 0775)
			tray.This.Update()
		}
	}
}

// loadHandler used to load the settings.
func loadHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "text/html")
		if b, err := os.ReadFile(config.Get().Config); err == nil {
			fmt.Fprint(rw, string(b))
		}
	}
}

// updateHandler used to update tray menu.
func updateHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		tray.This.Update()
	}
}

// Init used to start the tray and settings server.
func Init(ctx context.Context, cancel context.CancelFunc) {
	var onceFlags sync.Once
	srv := &http.Server{
		Addr:    config.Get().Address,
		Handler: routing(),
	}

	go func() {
		srv.ListenAndServe()
	}()

	onStart := func() {
		tray.This.OnReady()
	}

	onExit := func() {
		onceFlags.Do(
			func() {
				helpers.DeleteTmp(config.Get().Ico)
				srv.Shutdown(ctx)
				cancel()
			})
	}

	go systray.Run(onStart, onExit)
	<-ctx.Done()
	onExit()
}
