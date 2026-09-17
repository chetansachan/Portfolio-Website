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
	mux.Handle("POST /blog/createblog", middelware.AdminMiddleware(http.HandlerFunc(blogController.CreateBlog_Controller)))
	mux.HandleFunc("GET /blog/getallblogs", blogController.GetAllBlogs_Controller)
	mux.HandleFunc("GET /blog/getblog/slug/{slug}", blogController.GetBlogBySlug_Controller)
	mux.HandleFunc("GET /blog/getblog/id/{id}", blogController.GetBlogById_Controller)
	mux.Handle("DELETE /blog/deleteblog/{id}", middelware.AdminMiddleware(http.HandlerFunc(blogController.DeleteBlogById_Controller)))
	mux.Handle("PATCH /blog/updateblog/{id}", middelware.AdminMiddleware(http.HandlerFunc(blogController.UpdateBlogById_Controller)))
	mux.Handle("POST /project/createproject", middelware.AdminMiddleware(http.HandlerFunc(projectController.CreateProject_Controller)))
	mux.HandleFunc("GET /project/getproject/id/{id}", projectController.GetProjectById_Controller)
	mux.HandleFunc("GET /project/getproject/slug/{slug}", projectController.GetProjectBySlug_Controller)
	mux.Handle("DELETE /project/deleteproject/{id}", middelware.AdminMiddleware(http.HandlerFunc(projectController.DeleteProjectById_Controller)))
	mux.Handle("PATCH /project/updateproject/{id}", middelware.AdminMiddleware(http.HandlerFunc(projectController.UpdateProjectById_Controller)))
	mux.HandleFunc("GET /project/getallprojects", projectController.GetAllProjects_Controller)
	mux.Handle("POST /casestudy/createcasestudy", middelware.AdminMiddleware(http.HandlerFunc(caseStudyController.CreateCaseStudy_Controller)))
	mux.HandleFunc("GET /casestudy/getallcasestudies", caseStudyController.GetAllCaseStudies_Controller)
	mux.HandleFunc("GET /casestudy/getcasestudy/slug/{slug}", caseStudyController.GetCaseStudyBySlug_Controller)
	mux.HandleFunc("GET /casestudy/getcasestudy/id/{id}", caseStudyController.GetCaseStudyById_Controller)
	mux.Handle("DELETE /casestudy/deletecasestudy/{id}", middelware.AdminMiddleware(http.HandlerFunc(caseStudyController.DeleteCaseStudyById_Controller)))
	mux.Handle("PATCH /casestudy/updatecasestudy/{id}", middelware.AdminMiddleware(http.HandlerFunc(caseStudyController.UpdateCaseStudyById_Controller)))
	mux.HandleFunc("POST /subscribe", newsletterController.CreateSubscriber_Controller)
	mux.HandleFunc("POST /user/registeruser", authController.RegisterUser_Controller)
	mux.HandleFunc("POST /user/login", authController.Login_Controller)
	mux.Handle("POST /user/logout", middelware.AuthMiddleware(http.HandlerFunc(authController.Logout_Controller)))
	mux.HandleFunc("POST /user/refresh", authController.RefreshToken_Controller)
	mux.HandleFunc("POST /user/makeadmin",authController.MakeAdmin_Controller)
	mux.Handle("DELETE /user/deleteuser/{id}", middelware.AdminMiddleware(http.HandlerFunc(authController.DeleteUserByID_Controller)))
	mux.Handle("POST /upload", middelware.AdminMiddleware(http.HandlerFunc(uploadController.Upload_Controller)))

	return mux
}
