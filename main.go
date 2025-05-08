package main

import (
	"Go-HTTP-Balancer/rate_limiter"
	"log"
	"net/http"

	"Go-HTTP-Balancer/backend"
	"Go-HTTP-Balancer/config"
	"Go-HTTP-Balancer/lb"
	"Go-HTTP-Balancer/utils"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	backends, err := backend.InitBackends(cfg.Backends)
	if err != nil {
		log.Fatalf("Ошибка инициализации бэкендов: %v", err)
	}

	limiter := rate_limiter.NewLimiter(func(key string) (int, float64) {
		// Можно заменить на чтение из БД или конфиг-файла
		return 10, 1 // 10 токенов, 1 токен/сек
	})

	lb := lb.NewLoadBalancer(backends)

	go utils.HealthCheck(lb, 5)

	log.Printf("Балансировщик запущен на порту :%s", cfg.ListenPort)
	err = http.ListenAndServe(":"+cfg.ListenPort, limiter.Middleware(lb))
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
