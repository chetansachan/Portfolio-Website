package router

import (
	"fmt"
	"net/http"
	"portfolio-website/internals/config"
	"portfolio-website/internals/controller"
	"portfolio-website/internals/database"
	"portfolio-website/internals/middelware"
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
	projectService := service.NewProjectService(projectRepo, newsletterRepo)
	projectController := controller.NewProjectController(projectService)
	authRepo := reposiotry.NewAuthRepo(database.Pool)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)
	caseStudyRepo := reposiotry.NewCaseStudyRepository(database.Pool)
	caseStudyService := service.NewCaseStudyService(caseStudyRepo, newsletterRepo)
	caseStudyController := controller.NewCaseStudyController(caseStudyService)
	uploadRepo := reposiotry.NewUploadRepo(database.Pool)
	uploadService := service.NewUploadService(uploadRepo)
	uploadController := controller.NewUploadController(uploadService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello")
	})
	mux.Handle("POST /blog/createblog", middelware.AuthMiddleware(http.HandlerFunc(blogController.CreateBlog_Controller)))
	mux.HandleFunc("GET /blog/getallblogs", blogController.GetAllBlogs_Controller)
	mux.HandleFunc("GET /blog/getblog/slug/{slug}", blogController.GetBlogBySlug_Controller)
	mux.HandleFunc("GET /blog/getblog/id/{id}", blogController.GetBlogById_Controller)
	mux.Handle("DELETE /blog/deleteblog/{id}", middelware.AuthMiddleware(http.HandlerFunc(blogController.DeleteBlogById_Controller)))
	mux.Handle("PATCH /blog/updateblog/{id}", middelware.AuthMiddleware(http.HandlerFunc(blogController.UpdateBlogById_Controller)))
	mux.Handle("POST /project/createproject", middelware.AuthMiddleware(http.HandlerFunc(projectController.CreateProject_Controller)))
	mux.HandleFunc("GET /project/getproject/id/{id}", projectController.GetProjectById_Controller)
	mux.HandleFunc("GET /project/getproject/slug/{slug}", projectController.GetProjectBySlug_Controller)
	mux.Handle("DELETE /project/deleteproject/{id}", middelware.AuthMiddleware(http.HandlerFunc(projectController.DeleteProjectById_Controller)))
	mux.Handle("PATCH /project/updateproject/{id}", middelware.AuthMiddleware(http.HandlerFunc(projectController.UpdateProjectById_Controller)))
	mux.HandleFunc("GET /project/getallprojects", projectController.GetAllProjects_Controller)
	mux.Handle("POST /casestudy/createcasestudy", middelware.AuthMiddleware(http.HandlerFunc(caseStudyController.CreateCaseStudy_Controller)))
	mux.HandleFunc("GET /casestudy/getallcasestudies", caseStudyController.GetAllCaseStudies_Controller)
	mux.HandleFunc("GET /casestudy/getcasestudy/slug/{slug}", caseStudyController.GetCaseStudyBySlug_Controller)
	mux.HandleFunc("GET /casestudy/getcasestudy/id/{id}", caseStudyController.GetCaseStudyById_Controller)
	mux.Handle("DELETE /casestudy/deletecasestudy/{id}", middelware.AuthMiddleware(http.HandlerFunc(caseStudyController.DeleteCaseStudyById_Controller)))
	mux.Handle("PATCH /casestudy/updatecasestudy/{id}", middelware.AuthMiddleware(http.HandlerFunc(caseStudyController.UpdateCaseStudyById_Controller)))
	mux.HandleFunc("POST /subscribe", newsletterController.CreateSubscriber_Controller)
	mux.HandleFunc("POST /user/registeruser", authController.RegisterUser_Controller)
	mux.HandleFunc("POST /user/login", authController.Login_Controller)
	mux.Handle("POST /user/refresh", middelware.AuthMiddleware(http.HandlerFunc(authController.RefreshToken_Controller)))
	mux.Handle("POST /upload", middelware.AuthMiddleware(http.HandlerFunc(uploadController.Upload_Controller)))
	mux.HandleFunc("GET /files/{file}", uploadController.FindFile_Controller)
	mux.Handle("DELETE /deletefiles/{file}", middelware.AuthMiddleware(http.HandlerFunc(uploadController.DeleteFile_Controller)))

	return mux
}
