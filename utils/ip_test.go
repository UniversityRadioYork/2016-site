package utils

import (
	"github.com/UniversityRadioYork/2016-site/structs"
	"net/http"
	"testing"
)

var Config = structs.Config{TrustedProxies: []string{"127.0.0.1"}}

func TestRemoteIP(t *testing.T) {
	req := http.Request{RemoteAddr: "127.127.0.0:2134"}
	ip, err := GetRequesterIP(&Config, &req)
	if err != nil {
		t.Error(err)
	}
	if ip.String() != "127.127.0.0" {
		t.Errorf("wrong IP: got %v", ip)
	}
}

func TestXffInvalid(t *testing.T) {
	req := http.Request{
		RemoteAddr: "127.127.0.0:2134",
		Header: http.Header{
			"X-Forwarded-For": []string{"127.0.255.1"},
		},
	}
	ip, err := GetRequesterIP(&Config, &req)
	if err != nil {
		t.Error(err)
	}
	if ip.String() != "127.127.0.0" {
		t.Errorf("wrong IP: got %v", ip)
	}
}

func TestXffUntrusted(t *testing.T) {
	req := http.Request{
		RemoteAddr: "127.127.0.0:2134",
		Header: http.Header{
			"X-Forwarded-For": []string{"127.0.255.1, 127.0.255.2"},
		},
	}
	ip, err := GetRequesterIP(&Config, &req)
	if err != nil {
		t.Error(err)
	}
	if ip.String() != "127.127.0.0" {
		t.Errorf("wrong IP: got %v", ip)
	}
}

func TestXffValid(t *testing.T) {
	req := http.Request{
		RemoteAddr: "127.0.0.1:2134",
		Header: http.Header{
			"X-Forwarded-For": []string{"127.0.100.1, 127.0.0.1"},
		},
	}
	ip, err := GetRequesterIP(&Config, &req)
	if err != nil {
		t.Error(err)
	}
	if ip.String() != "127.0.100.1" {
		t.Errorf("wrong IP: got %v", ip)
	}
}

// This is really a pathological case.
func TestInvalidRemoteAddr(t *testing.T) {
	req := http.Request{
		RemoteAddr: "127.0.0.1",
	}
	_, err := GetRequesterIP(&Config, &req)
	if err == nil {
		t.Fatal("Expected error, got none")
	}
}
