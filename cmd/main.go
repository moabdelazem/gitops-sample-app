package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/moabdelazem/gitops-sample-app/internal/config"
	"github.com/moabdelazem/gitops-sample-app/internal/handler"
)

func main() {
	cfg := config.Load()

	templatePath := filepath.Join("web", "templates", "index.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("failed to parse template %s: %v", templatePath, err)
	}

	h := handler.New(cfg, tmpl)
	mux := http.NewServeMux()
	h.Register(mux)

	log.Printf("server starting on :%s [env=%s]", cfg.Port, cfg.Environment)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
