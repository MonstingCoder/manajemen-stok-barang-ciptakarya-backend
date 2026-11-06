package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/database"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/handler"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/middleware"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/repository"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/service"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/utility"
)

func main() {
	run()
}

func run() {
	cfg, err := config.NewConfig()
	if err != nil {
		mustError(err)
	}
	db, err := database.NewDB(cfg)
	if err != nil {
		mustError(err)
	}

	userR := repository.NewUserRepo(db)
	roleR := repository.NewRoleRepo(db)
	itemR := repository.NewItemRepo(db)
	tagR := repository.NewTagRepo(db)

	userS := service.NewUserService(userR)
	roleS := service.NewRoleService(roleR)
	itemS := service.NewItemService(itemR)
	tagS := service.NewTagService(tagR)

	userH := handler.NewUserHandler(userS, cfg)
	roleH := handler.NewRoleHandler(roleS)
	itemH := handler.NewItemHandler(itemS)
	tagH := handler.NewTagHandler(tagS)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		utility.WriteResponse(w, http.StatusOK, map[string]string{"message": "server healthy"})
	})

	ms := middleware.NewMiddlewares(cfg)
	adminOnly := ms.Autho("admin")

	mux.HandleFunc("POST /reg", userH.Reg)
	mux.HandleFunc("POST /login", userH.Login)

	mux.Handle("GET /users", middleware.Chain(http.HandlerFunc(userH.Finds), ms.Auth, adminOnly))
	mux.Handle("PATCH /users/{id}", middleware.Chain(http.HandlerFunc(userH.Update), ms.Auth))
	mux.Handle("DELETE /users/{id}", middleware.Chain(http.HandlerFunc(userH.Delete), ms.Auth, adminOnly))

	mux.Handle("DELETE /roles/{id}", middleware.Chain(http.HandlerFunc(roleH.Delete), ms.Auth, adminOnly))

	mux.Handle("POST /items", middleware.Chain(http.HandlerFunc(itemH.Create), ms.Auth, adminOnly))
	mux.Handle("GET /items", middleware.Chain(http.HandlerFunc(itemH.Finds), ms.Auth))
	mux.Handle("GET /items/{name}", middleware.Chain(http.HandlerFunc(itemH.FindByName), ms.Auth))
	mux.Handle("PATCH /items/{id}", middleware.Chain(http.HandlerFunc(itemH.Update), ms.Auth, adminOnly))
	mux.Handle("DELETE /items/{id}", middleware.Chain(http.HandlerFunc(itemH.Delete), ms.Auth, adminOnly))

	mux.Handle("POST /tags", middleware.Chain(http.HandlerFunc(tagH.Create), ms.Auth, adminOnly))
	mux.Handle("GET /tags", middleware.Chain(http.HandlerFunc(tagH.Finds), ms.Auth))
	mux.Handle("GET /tags/{name}", middleware.Chain(http.HandlerFunc(tagH.Finds), ms.Auth))
	mux.Handle("DELETE /tags/{id}", middleware.Chain(http.HandlerFunc(tagH.Delete), ms.Auth, adminOnly))

	h := middleware.Chain(
		mux,
		ms.CORS,
	)
	srv := new(http.Server{
		Addr:         *cfg.AppPort,
		Handler:      h,
		ReadTimeout:  *cfg.DefaultDuration / 2,
		WriteTimeout: *cfg.DefaultDuration,
	})
	go func() {
		fmt.Println("server running on port", *cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintln("runing server error:", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down server")
	if err = srv.Close(); err != nil {
		panic(fmt.Sprintln("shutdown server error:", err))
	}
	fmt.Println("shutdown server success")
}

func mustError(err error) {
	fmt.Println("error:", err)
	os.Exit(1)
}
