package lb

import (
	"log"
	"net/http"
	"sync/atomic"

	"Go-HTTP-Balancer/backend"
)

type LoadBalancer struct {
	backends []*backend.Backend
	index    uint32
}

// NewLoadBalancer создает новый балансировщик на основе списка URL бэкендов.
func NewLoadBalancer(b []*backend.Backend) *LoadBalancer {
	return &LoadBalancer{
		backends: b,
	}
}

// ServeHTTP реализует интерфейс http.Handler для LoadBalancer.
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b := lb.NextBackend()
	if b == nil {
		http.Error(w, "Нет доступных серверов", http.StatusServiceUnavailable)
		return
	}
	log.Printf("Запрос перенаправлен на %s%s", b.URL, r.URL.Path)
	b.ReverseProxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) NextBackend() *backend.Backend {
	total := uint32(len(lb.backends))
	for i := uint32(0); i < total; i++ {
		idx := atomic.AddUint32(&lb.index, 1) % total
		b := lb.backends[idx]
		if b.Alive {
			return b
		}
	}
	return nil
}

func (lb *LoadBalancer) GetBackends() []*backend.Backend {
	return lb.backends
}
