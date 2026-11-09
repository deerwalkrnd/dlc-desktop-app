package server

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

// placeholderHTML is served when the binary was built without the web app.
const placeholderHTML = `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>DLC Desktop App</title></head>
<body style="font-family:system-ui,sans-serif;max-width:40rem;margin:4rem auto;padding:0 1rem">
<h1>The web interface is not built into this binary</h1>
<p>Build the frontend first, then rebuild the app:</p>
<pre>make build</pre>
<p>During development you can also run <code>npm run dev</code> in <code>web/</code>
and open the Vite URL; it proxies API calls to this server.</p>
</body></html>`

type spaHandler struct {
	files http.Handler
	ui    fs.FS
	index []byte
}

func newSPA(ui fs.FS) (http.Handler, error) {
	if ui == nil {
		return placeholder(), nil
	}
	index, err := fs.ReadFile(ui, "index.html")
	if errors.Is(err, fs.ErrNotExist) {
		return placeholder(), nil
	}
	if err != nil {
		return nil, err
	}
	return &spaHandler{files: http.FileServerFS(ui), ui: ui, index: index}, nil
}

func placeholder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(placeholderHTML))
	})
}

// ServeHTTP serves built assets as files and every other path as index.html,
// letting the client-side router handle it. Directory listings are never
// produced: only regular files reach the file server.
func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
	if p != "" && p != "index.html" {
		if fi, err := fs.Stat(h.ui, p); err == nil && fi.Mode().IsRegular() {
			if strings.HasPrefix(p, "_app/immutable/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			h.files.ServeHTTP(w, r)
			return
		}
		if p == "_app" || strings.HasPrefix(p, "_app/") {
			// A missing JS/CSS chunk must not be answered with HTML.
			http.NotFound(w, r)
			return
		}
	}
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(h.index))
}
