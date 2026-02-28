package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/moabdelazem/gitops-sample-app/internal/config"
	"github.com/moabdelazem/gitops-sample-app/internal/model"
	"github.com/moabdelazem/gitops-sample-app/pkg/version"
)

type Handler struct {
	cfg          *config.Config
	tmpl         *template.Template
	startTime    time.Time
	requestCount atomic.Int64
}

func New(cfg *config.Config, tmpl *template.Template) *Handler {
	return &Handler{
		cfg:       cfg,
		tmpl:      tmpl,
		startTime: time.Now(),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("GET /healthz", h.Health)
	mux.HandleFunc("GET /api/info", h.APIInfo)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.requestCount.Add(1)
	info := h.buildInfo()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.Execute(w, info); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("template execution error: %v", err)
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) APIInfo(w http.ResponseWriter, r *http.Request) {
	info := h.buildInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (h *Handler) buildInfo() model.AppInfo {
	return model.AppInfo{
		Version:      version.Version,
		GitCommit:    version.GitCommit,
		BuildTime:    version.BuildTime,
		Environment:  h.cfg.Environment,
		PodName:      h.cfg.PodName,
		NodeName:     h.cfg.NodeName,
		HostIP:       getHostIP(),
		GoVersion:    runtime.Version(),
		Uptime:       formatUptime(time.Since(h.startTime)),
		RequestCount: h.requestCount.Load(),
	}
}

func getHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "unknown"
}

func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, mins, secs)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, mins, secs)
	}
	if mins > 0 {
		return fmt.Sprintf("%dm %ds", mins, secs)
	}
	return fmt.Sprintf("%ds", secs)
}
