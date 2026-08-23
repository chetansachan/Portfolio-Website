package router

import (
	"fmt"
	"net/http"
	"portfolio-website/internals/config"
	"portfolio-website/internals/controller"
	"portfolio-website/internals/database"
	"portfolio-website/internals/reposiotry"
	"portfolio-website/internals/service"
)

func Router() http.Handler {
	mux := http.NewServeMux()

	blogRepo := reposiotry.NewBlogRepository(database.Pool)
	newsletterRepo := reposiotry.NewNewsletterRepo(database.Pool, config.NewClient())
	blogService := service.NewBlogService(blogRepo, newsletterRepo)
	newsletterService := service.NewNewsletterService(newsletterRepo)
	newsletterController := controller.NewNewsletterController(newsletterService)
	blogController := controller.NewBlogController(blogService)
	projectRepo := reposiotry.NewProjectRepository(database.Pool)
	projectService := service.NewProjectService(projectRepo)
	projectController := controller.NewProjectController(projectService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello")
	})
	mux.HandleFunc("POST /blog/createblog", blogController.CreateBlog_Controller)
	mux.HandleFunc("GET /blog/getallblogs", blogController.GetAllBlogs_Controller)
	mux.HandleFunc("GET /blog/getblog/slug/{slug}", blogController.GetBlogBySlug_Controller)
	mux.HandleFunc("GET /blog/getblog/id/{id}", blogController.GetBlogById_Controller)
	mux.HandleFunc("DELETE /blog/deleteblog/{id}", blogController.DeleteBlogById_Controller)
	mux.HandleFunc("PATCH /blog/updateblog/{id}", blogController.UpdateBlogById_Controller)
	mux.HandleFunc("POST /project/createproject", projectController.CreateProject_Controller)
	mux.HandleFunc("GET /project/getproject/id/{id}", projectController.GetProjectById_Controller)
	mux.HandleFunc("GET /project/getproject/slug/{slug}", projectController.GetProjectBySlug_Controller)
	mux.HandleFunc("DELETE /project/deleteproject/{id}", projectController.DeleteProjectById_Controller)
	mux.HandleFunc("PATCH /project/updateproject/{id}", projectController.UpdateProjectById_Controller)
	mux.HandleFunc("GET /project/getallprojects", projectController.GetAllProjects_Controller)
	mux.HandleFunc("POST /blog/subscribe", newsletterController.CreateSubscriber_Controller)

	return mux
}
