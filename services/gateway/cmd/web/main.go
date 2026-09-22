package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load("services/gateway/.env")

	router := chi.NewRouter()

	// todo service
	todo := os.Getenv("TODO_URL")
	todoURL, err := url.Parse(todo)
	if err != nil {
		log.Printf("[Gateway] todo service error: %v", err)
	} else {
		serviceCheck(ctx, todo, "todo")

		todoProxy := httputil.NewSingleHostReverseProxy(todoURL)
		proxyRoute(router, "/v1/todo", todoProxy)
	}

	// history service
	history := os.Getenv("HISTORY_URL")
	historyURL, err := url.Parse(history)
	if err != nil {
		log.Printf("[Gateway] history service error: %v", err)
	} else {
		serviceCheck(ctx, history, "history") // check service it running
		historyProxy := httputil.NewSingleHostReverseProxy(historyURL)
		proxyRoute(router, "/v1/history", historyProxy)
	}

	log.Print("[Gateway] up and running...")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router); err != nil {
		log.Fatalf("[Gateway] %v", err)
	}
}

// Proxying the services
func proxyRoute(router chi.Router, prefix string, proxy *httputil.ReverseProxy) {
	router.HandleFunc(prefix+"/*", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		proxy.ServeHTTP(w, r)
	})
}

// Call health endpoint of the service
func isServiceHealthy(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return false, err
	}
	return true, nil
}

// Check if service is up
func serviceCheck(ctx context.Context, url, service string) {

	timeout := 30 * time.Second
	retryDelay := 5 * time.Second

	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		isUp, err := isServiceHealthy(fmt.Sprintf("%s/health", url))
		if err == nil && isUp {
			log.Printf("[Gateway] ✅ %s service is up and ready\n", service)
			return
		}

		select {
		case <-checkCtx.Done():
			log.Printf("[Gateway] 🛑 %s service is not up\n", service)
			return
		case <-time.After(retryDelay):
		}
	}
}
