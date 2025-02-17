package routes

import (
	"net/http"
)

// ルーティングを設定
func SetupRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
}
