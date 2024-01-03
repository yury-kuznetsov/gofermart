package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yury-kuznetsov/gofermart/cmd/gophermart/config"
	"github.com/yury-kuznetsov/gofermart/cmd/gophermart/handlers"
	userRepository "github.com/yury-kuznetsov/gofermart/internal/user/repository"
	userService "github.com/yury-kuznetsov/gofermart/internal/user/service"
	"github.com/yury-kuznetsov/gofermart/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()

	// создаем сервер
	server := &http.Server{Addr: config.Options.HostAddr, Handler: service()}

	// готовим канал для прослушивания системных сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// запускаем сервера в отдельной горутине
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// ожидаем сигнала остановки из канала `stop`
	<-stop

	// даем серверу 5 секунд на завершение обработки текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// завершаем "мягко" работу сервера
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server Shutdown: %v", err)
	}
}

func service() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)

	db, err := sql.Open("pgx", config.Options.DatabaseAddr)
	if err != nil {
		log.Fatal(err)
	}

	// сервисы аутентификации
	userRepo := userRepository.NewUserRepository(db)
	userSvc := userService.NewUserService(userRepo)
	jwtSvc := userService.NewTokenService()

	r.Post("/api/user/register", handlers.RegisterHandler(userSvc, jwtSvc))
	r.Post("/api/user/login", handlers.LoginHandler(userSvc, jwtSvc))

	return r
}
