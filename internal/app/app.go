package app

import (
	"LoadBalancer/internal/balancer"
	"LoadBalancer/internal/config"
	"LoadBalancer/internal/handlers"
	"LoadBalancer/internal/proxy"
	"LoadBalancer/internal/ratelimiter"
	"LoadBalancer/internal/repository"
	"LoadBalancer/pkg/logger"
	"LoadBalancer/pkg/postgres"
	"LoadBalancer/pkg/server"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Run() {
	log := logger.NewLogger()
	log.Info("Запуск балансировщика нагрузки...")

	cfg, err := config.RunConfig(log)
	if err != nil {
		log.Fatalf("Ошибка чтения конфига: %s", err)
	}

	db, err := postgres.ClientPostgres(log)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	repo := repository.NewRepository(db, log)

	// Балансировщик и лимитер
	balanc := balancer.NewBalancer(cfg.Backends, log)
	limiter := ratelimiter.NewRateLimiterWithRepository(repo)

	// HTTP-мультиплексор
	mux := http.NewServeMux()

	// Роуты
	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		handler := handlers.ClientHandler{Repo: repo}
		switch r.Method {
		case http.MethodGet:
			handler.ListClients(w, r)
		case http.MethodPost:
			handler.CreateOrUpdateClient(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/clients/", func(w http.ResponseWriter, r *http.Request) {
		handler := handlers.ClientHandler{Repo: repo}
		id := strings.TrimPrefix(r.URL.Path, "/clients/")
		r = r.WithContext(context.WithValue(r.Context(), "id", id))

		switch r.Method {
		case http.MethodGet:
			handler.GetClient(w, r)
		case http.MethodDelete:
			handler.DeleteClient(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Прокси на все остальные маршруты
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.RemoteAddr

		if !limiter.Allow(clientID) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{
				"code":    429,
				"message": "Rate limit exceeded",
			})
			return
		}

		backend := balanc.Next()
		proxyHandler, err := proxy.NewReverseProxy(backend, log)
		if err != nil {
			log.Errorf("Ошибка создания прокси: %v", err)
			http.Error(w, `{"error": "internal proxy error"}`, http.StatusInternalServerError)
			return
		}

		proxyHandler.ServeHTTP(w, r)
	})

	srv := new(server.Server)

	// Асинхронный запуск
	go func() {
		if err := srv.RunServer(mux, cfg.Port); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()
	log.Infof("Сервер запущен на порту: %s!", cfg.Port)

	// Завершение
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Выключение сервера...")
	if err := srv.ShutdownServer(); err != nil {
		log.Errorf("Ошибка при завершении работы сервера: %v", err)
	}
	log.Info("Сервер успешно выключен!")
}
