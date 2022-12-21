package main

import (
	"context"
	"embed"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/SemmiDev/kanbanapp/internal/client"
	"github.com/SemmiDev/kanbanapp/internal/handler/api"
	"github.com/SemmiDev/kanbanapp/internal/handler/gapi"
	"github.com/SemmiDev/kanbanapp/internal/handler/web"
	"github.com/SemmiDev/kanbanapp/internal/middleware"
	"github.com/SemmiDev/kanbanapp/internal/pb"
	"github.com/SemmiDev/kanbanapp/internal/repository"
	"github.com/SemmiDev/kanbanapp/internal/service"
	"github.com/SemmiDev/kanbanapp/internal/utils"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler     api.UserAPI
	TaskAPIHandler     api.TaskAPI
	CategoryAPIHandler api.CategoryAPI
}

type ClientHandler struct {
	AuthWeb      web.AuthWeb
	DashboardWeb web.DashboardWeb
	ModifyWeb    web.ModifyWeb
	HomeWeb      web.HomeWeb
}

//go:embed views/*
var Resources embed.FS

//go:embed static/*
var staticFiles embed.FS

func FlyURL() string {
	return "http://localhost:8080"
}

func main() {
	//TODO: hapus jika sudah di deploy di fly.io
	os.Setenv("DATABASE_URL", "postgres://root:secret@localhost:5432/kampusmerdeka")

	err := utils.ConnectDB()
	if err != nil {
		panic(err)
	}

	db := utils.GetDBConnection()

	go RunGrpcServer(db)
	go runGatewayServer(db)

	mux := http.NewServeMux()
	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	// Serve static files
	mux.Handle("/static/", fs)

	mux = RunHttpServer(db, mux)
	mux = RunClient(mux, Resources)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func RunGrpcServer(db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo, categoryRepo)

	server, err := gapi.NewServer(userService)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterKanbanServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runGatewayServer(db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userService := service.NewUserService(userRepo, categoryRepo)

	server, err := gapi.NewServer(userService)
	if err != nil {
		log.Fatal().Err(err).Msg("scannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterKanbanHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

func RunHttpServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	userService := service.NewUserService(userRepo, categoryRepo)
	taskService := service.NewTaskService(taskRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo, taskRepo)

	userAPIHandler := api.NewUserAPI(userService)
	taskAPIHandler := api.NewTaskAPI(taskService)
	categoryAPIHandler := api.NewCategoryAPI(categoryService)

	apiHandler := APIHandler{
		UserAPIHandler:     userAPIHandler,
		TaskAPIHandler:     taskAPIHandler,
		CategoryAPIHandler: categoryAPIHandler,
	}

	MuxRoute(mux, "POST", "/api/v1/users/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Login)))
	MuxRoute(mux, "POST", "/api/v1/users/register", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Register)))
	MuxRoute(mux, "POST", "/api/v1/users/logout", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Logout)))
	MuxRoute(mux, "DELETE", "/api/v1/users/delete", middleware.Delete(http.HandlerFunc(apiHandler.UserAPIHandler.Delete)), "?user_id=")

	MuxRoute(mux, "GET", "/api/v1/tasks/get", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.GetTask))), "?task_id=")
	MuxRoute(mux, "POST", "/api/v1/tasks/create", middleware.Post(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.CreateNewTask))))
	MuxRoute(mux, "PUT", "/api/v1/tasks/update", middleware.Put(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.UpdateTask))), "?task_id=")
	MuxRoute(mux, "PUT", "/api/v1/tasks/update/category", middleware.Put(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.UpdateTaskCategory))), "?task_id=")
	MuxRoute(mux, "DELETE", "/api/v1/tasks/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.DeleteTask))), "?task_id=")

	MuxRoute(mux, "GET", "/api/v1/categories/get", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.GetCategory))))
	MuxRoute(mux, "GET", "/api/v1/categories/dashboard", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.GetCategoryWithTasks))))
	MuxRoute(mux, "POST", "/api/v1/categories/create", middleware.Post(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.CreateNewCategory))))
	MuxRoute(mux, "DELETE", "/api/v1/categories/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.DeleteCategory))), "?category_id=")

	return mux
}

func RunClient(mux *http.ServeMux, embed embed.FS) *http.ServeMux {
	userClient := client.NewUserClient()
	categoryClient := client.NewCategoryClient()
	taskClient := client.NewTaskClient()

	authWeb := web.NewAuthWeb(userClient, embed)
	dashboardWeb := web.NewDashboardWeb(categoryClient, embed)
	modifyWeb := web.NewModifyWeb(taskClient, categoryClient, embed)
	homeWeb := web.NewHomeWeb(embed)

	client := ClientHandler{
		authWeb, dashboardWeb, modifyWeb, homeWeb,
	}

	mux.HandleFunc("/login", client.AuthWeb.Login)
	mux.HandleFunc("/login/process", client.AuthWeb.LoginProcess)

	mux.HandleFunc("/register", client.AuthWeb.Register)
	mux.HandleFunc("/register/process", client.AuthWeb.RegisterProcess)

	mux.HandleFunc("/logout", client.AuthWeb.Logout)

	mux.Handle("/dashboard", middleware.Auth(http.HandlerFunc(client.DashboardWeb.Dashboard)))

	mux.Handle("/category/add", middleware.Auth(http.HandlerFunc(client.ModifyWeb.AddCategory)))
	mux.Handle("/category/create", middleware.Auth(http.HandlerFunc(client.ModifyWeb.AddCategoryProcess)))

	mux.Handle("/task/add", middleware.Auth(http.HandlerFunc(client.ModifyWeb.AddTask)))
	mux.Handle("/task/create", middleware.Auth(http.HandlerFunc(client.ModifyWeb.AddTaskProcess)))

	mux.Handle("/task/update", middleware.Auth(http.HandlerFunc(client.ModifyWeb.UpdateTask)))
	mux.Handle("/task/update/process", middleware.Auth(http.HandlerFunc(client.ModifyWeb.UpdateTaskProcess)))

	mux.Handle("/task/delete", middleware.Auth(http.HandlerFunc(client.ModifyWeb.DeleteTask)))
	mux.Handle("/category/delete", middleware.Auth(http.HandlerFunc(client.ModifyWeb.DeleteCategory)))

	mux.HandleFunc("/", client.HomeWeb.Index)

	return mux
}

func MuxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}
