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
	projectRepo := reposiotry.NewProjectRepository(database.Pool)
	projectService := service.NewProjectService(projectRepo)
	projectController := controller.NewProjectController(projectService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello")
	})
	mux.HandleFunc("POST /blog/postblog", blogController.PostBlog_Controller)
	mux.HandleFunc("GET /blog/getblog/{slug}", blogController.GetBlogBySlug_Controller)
	mux.HandleFunc("GET /blog/getblog/{id}", blogController.GetBlogById_Controller)
	mux.HandleFunc("DELETE /blog/deleteblog/{id}", blogController.DeleteBlogById_Controller)
	mux.HandleFunc("PATCH /blog/updateblog/{id}", blogController.UpdateBlogById_Controller)
	mux.HandleFunc("POST /project/postproject", projectController.PostProject_Controller)
	
	return mux
}
