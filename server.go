package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type staticHandler struct {
	root             string
	notFoundPath     string
	notFoundRedirect bool
	filemap          map[string]bool
	auth             bool
	authUser         string
	authPass         string
}

func newStaticHandler(root string) (*staticHandler, error) {
	filemap := map[string]bool{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			filemap[strings.Replace(path, root, "", 1)] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &staticHandler{root: root, filemap: filemap}, nil
}

func (h *staticHandler) useAuth(user, password string) {
	h.auth = true
	h.authUser = user
	h.authPass = password
}

func (h *staticHandler) serveFile(status int, w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(h.root + r.URL.Path)
	if err != nil {
		log.Println("server error:", err)
		http.Error(w, "file read error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Println("server error:", err)
		http.Error(w, "file read error", http.StatusInternalServerError)
		return
	}

	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	http.ServeContent(w, r, r.URL.Path, info.ModTime(), file)
}

func (h *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.auth {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user != h.authUser || pass != h.authPass {
			http.Error(w, "permission denied", http.StatusUnauthorized)
			return
		}
	}

	path := r.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(path, "/") {
		path = path + "index.html"
	}

	if !h.filemap[path] {
		path = path + ".html"
	}
	if !h.filemap[path] {
		if h.notFoundPath != "" && h.filemap[h.notFoundPath] {
			r.URL.Path = h.notFoundPath
			h.serveFile(http.StatusNotFound, w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	r.URL.Path = path
	h.serveFile(http.StatusOK, w, r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "build"
	}

	handler, err := newStaticHandler(staticDir)
	if err != nil {
		log.Fatal("cant setup handler:", err)
		return
	}

	if os.Getenv("AUTH_USER") != "" && os.Getenv("AUTH_PASSWORD") != "" {
		handler.useAuth(os.Getenv("AUTH_USER"), os.Getenv("AUTH_PASSWORD"))
	}

	if notFoundPath := os.Getenv("NOT_FOUND_PATH"); notFoundPath != "" {
		handler.notFoundPath = notFoundPath
	}

	log.Println("starting server on port", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("server error:", err)
	}
}
