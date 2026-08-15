package router

import "net/http"

func Router() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /portfolio/blog", )
}