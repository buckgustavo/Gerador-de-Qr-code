package handler

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	contagem map[string]int
	limite   int
}

func NovoLimiter(limite int) *Limiter {
	l := &Limiter{
		contagem: make(map[string]int),
		limite:   limite,
	}
	go func() {
		for range time.Tick(time.Minute) {
			l.mu.Lock()
			l.contagem = make(map[string]int)
			l.mu.Unlock()
		}
	}()
	return l
}

func (l *Limiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := ipReal(r)
		l.mu.Lock()
		l.contagem[ip]++
		excedeu := l.contagem[ip] > l.limite
		l.mu.Unlock()

		if excedeu {
			escreverJSON(w, http.StatusTooManyRequests, respostaGerar{
				Erro: "muitas tentativas, aguarde um minuto",
			})
			return
		}
		next(w, r)
	}
}

func ipReal(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
