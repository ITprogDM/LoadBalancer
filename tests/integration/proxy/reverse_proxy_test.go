package proxy

import (
	"LoadBalancer/internal/proxy"
	"LoadBalancer/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewReverseProxy_ServeHTTP(t *testing.T) {
	log := logger.NewLogger()
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Привет от сервера"))
	}))
	defer backend.Close()

	proxy, err := proxy.NewReverseProxy(backend.URL, log)
	if err != nil {
		t.Fatalf("Ошибка создания прокси: %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	proxy.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Ожидали 200 OK, got %d", recorder.Code)
	}

	if body := recorder.Body.String(); body != "Привет от сервера" {
		t.Errorf("Неожиданное поведение: %s", body)
	}
}
