package utils

import (
	"log"
	"net/http"
	"time"

	"Go-HTTP-Balancer/backend"
	"Go-HTTP-Balancer/lb"
)

// HealthCheck проверяет доступность всех бэкендов.
func HealthCheck(lb *lb.LoadBalancer, intervalSec int) {
	for {
		for _, b := range lb.GetBackends() {
			alive := isAlive(b)
			if b.Alive != alive {
				status := "восстановлен"
				if !alive {
					status = "недоступен"
				}
				log.Printf("Бэкенд %s стал %s", b.URL, status)
			}
			b.Alive = alive
		}
		time.Sleep(time.Duration(intervalSec) * time.Second)
	}
}

func isAlive(b *backend.Backend) bool {
	resp, err := http.Get(b.URL.String() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
