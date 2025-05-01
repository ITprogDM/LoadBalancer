package proxy

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string, log *logrus.Logger) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Printf("Ошибка при парсинге данного URL: %s", target)
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	origDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		origDirector(req)
		log.Printf("Проксируем запрос на: %s", target)
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Ошибка проксирования на сервер: %s ", err)
		http.Error(w, "Сервер недоступен", http.StatusServiceUnavailable)
	}

	return proxy, nil
}
