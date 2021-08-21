package utils

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"net"
	"net/http"
	"strings"
)

func GetRequesterIP(c *structs.Config, r *http.Request) (net.IP, error) {
	var ip net.IP
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	ip = net.ParseIP(host)

	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ", ")
		if len(ips) >= 2 {
			lastProxy := ips[len(ips)-1]
			for _, trusted := range c.TrustedProxies {
				if strings.TrimSpace(lastProxy) == trusted {
					ip = net.ParseIP(ips[0])
				}
			}
		}
	}

	return ip, nil
}
