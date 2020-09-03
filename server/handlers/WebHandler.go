package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type WebOptions struct {
	PublicPath   string
	IndexFile    string
	NotFoundFile string
}

type WebHandler struct {
	WebOptions
	FileServer http.Handler
}

var _ RouteHandler = (*WebHandler)(nil)

func NewWebHandler(opts WebOptions) *WebHandler {
	fs := http.FileServer(http.Dir(opts.PublicPath))
	return &WebHandler{
		WebOptions: opts,
		FileServer: fs,
	}
}

func (h *WebHandler) RegisterRoutes(r *mux.Router) {
	r.PathPrefix("/s").HandlerFunc(h.HandleFile("s/[id].html"))
	r.PathPrefix("/").HandlerFunc(h.HandleRoot)
}

func (h *WebHandler) HandleRoot(w http.ResponseWriter, req *http.Request) {
	path, err := filepath.Abs(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = os.Stat(filepath.Join(h.PublicPath, path))
	if os.IsNotExist(err) {
		http.ServeFile(w, req, filepath.Join(h.PublicPath, h.NotFoundFile))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	h.FileServer.ServeHTTP(w, req)
}

func (h *WebHandler) HandleFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath.Join(h.PublicPath, path))
	}
}
