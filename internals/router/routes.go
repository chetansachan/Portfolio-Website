package router

import (
	"fmt"
	"net/http"
	"portfolio-website/internals/controller"
	"portfolio-website/internals/database"
	"portfolio-website/internals/reposiotry"
	"portfolio-website/internals/service"
)

func Router() http.Handler {
	mux := http.NewServeMux()
	
	blogRepo := reposiotry.NewBlogRepository(database.Pool)
	blogService := service.NewBlogService(blogRepo)
	blogController := controller.NewBlogController(blogService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello")
	})
	mux.HandleFunc("POST /portfolio/blog", blogController.PostBlog_Controller)

	return mux
}
