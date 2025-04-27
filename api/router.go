package api

import "net/http"

func GetApiMux() *http.ServeMux {
	var ApiMux = http.NewServeMux()

	ApiMux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"users": ["John", "Jane", "Bob"]}`))
	})

	return ApiMux
}
