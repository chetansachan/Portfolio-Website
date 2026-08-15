package router

import (
	"fmt"
	"net/http"
	"portfolio-website/internals/controller"
)

func Router() {
	mux := http.NewServeMux()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello")
	})
	mux.HandleFunc("POST /portfolio/blog", controller.PostBlog_Controller)
}
