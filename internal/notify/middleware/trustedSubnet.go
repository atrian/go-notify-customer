package middlewares

import (
	"net"
	"net/http"
)

func TrustedSubnetMW(trustedSubnet string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// При пустом значении trustedSubnet обрабатываем запрос без ограничений
			if trustedSubnet == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Если парсинг конфигурации безопасной сети выдает ошибку, передаем запрос дальше
			_, ipNet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Если адрес входит в список доверенных, передаем запрос дальше
			agentIp := net.ParseIP(r.Header.Get("X-Real-IP"))
			if agentIp != nil && ipNet.Contains(agentIp) {
				next.ServeHTTP(w, r)
				return
			}

			// иначе сбрасываем соединение
			dropConnection(w)
		})
	}
}

func dropConnection(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}
