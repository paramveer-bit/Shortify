package helper

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	// Try to get the IP address from the X-Forwarded-For header first
	ip := r.Header.Get("X-FORWARDED-FOR")
	if ip != "" {
		// The X-Forwarded-For header can contain multiple IPs, so we take the first one
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// If no X-Forwarded-For header, use the RemoteAddr (client IP from the connection)
	ip = r.RemoteAddr

	// Use net.SplitHostPort to separate the IP from the port
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return ip // Return the original IP if there's an error
	}
	return host
}
