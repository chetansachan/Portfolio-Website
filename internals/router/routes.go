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
	mux.HandleFunc("POST /blog/createblog", blogController.CreateBlog_Controller)
	mux.HandleFunc("GET /blog/getallblogs", blogController.GetAllBlogs_Controller)
	mux.HandleFunc("GET /blog/featured", blogController.GetFeaturedBlogs_Controller)
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
	mux.HandleFunc("GET /project/featured", projectController.GetFeaturedProjects_Controller)
	mux.HandleFunc("POST /casestudy/createcasestudy", caseStudyController.CreateCaseStudy_Controller)
	mux.HandleFunc("GET /casestudy/getallcasestudies", caseStudyController.GetAllCaseStudies_Controller)
	mux.HandleFunc("GET /casestudy/featured", caseStudyController.GetFeaturedCaseStudies_Controller)
	mux.HandleFunc("GET /casestudy/getcasestudy/slug/{slug}", caseStudyController.GetCaseStudyBySlug_Controller)
	mux.HandleFunc("GET /casestudy/getcasestudy/id/{id}", caseStudyController.GetCaseStudyById_Controller)
	mux.HandleFunc("DELETE /casestudy/deletecasestudy/{id}", caseStudyController.DeleteCaseStudyById_Controller)
	mux.HandleFunc("PATCH /casestudy/updatecasestudy/{id}", caseStudyController.UpdateCaseStudyById_Controller)
	mux.HandleFunc("POST /subscribe", newsletterController.CreateSubscriber_Controller)
	mux.HandleFunc("POST /user/registeruser", authController.RegisterUser_Controller)
	mux.HandleFunc("GET /user/login", authController.Login_Controller)
	mux.HandleFunc("POST /user/refresh", authController.RefreshToken_Controller)
	mux.HandleFunc("POST /upload", uploadController.Upload_Controller)
	mux.HandleFunc("GET /files/{file}", uploadController.FindFile_Controller)
	mux.HandleFunc("DELETE /deletefiles/{file}", uploadController.DeleteFile_Controller)


	return mux
}
