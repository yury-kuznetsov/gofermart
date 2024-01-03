package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yury-kuznetsov/gofermart/cmd/gophermart/config"
	"github.com/yury-kuznetsov/gofermart/middleware"
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

	return r
}
