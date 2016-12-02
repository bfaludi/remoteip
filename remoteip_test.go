package remoteip

import (
	"net"
	"net/http"
	"testing"
)

func TestFunctionIsIPv4Address(t *testing.T) {
	cases := []struct {
		IP      string
		IsValid bool
	}{
		{"this is a text", false},
		{"192.160.0.1", true},
		{"2001:0db8:0a0b:12f0:0000:0000:0000:0001", false},
		{"", false}}

	for _, c := range cases {
		if exp := IsIPv4Address(net.ParseIP(c.IP)); exp != c.IsValid {
			t.Errorf("IP address `%s` should be %s instead of %s", c.IP, c.IsValid, exp)
		}
	}
}

func TestFunctionIsPrivateIPv4Address(t *testing.T) {
	cases := []struct {
		IP        string
		IsPrivate bool
	}{
		{"this is a tex", false},
		{"10.0.0.0", true},
		{"10.255.255.255", true},
		{"192.160.43.1", false},
		{"192.168.0.1", true},
		{"172.16.4.1", true},
		{"172.13.0.255", false},
		{"10.65.4.32", true},
		{"10.255.43.1", true},
		{"2001:0db8:0a0b:12f0:0000:0000:0000:0001", false},
		{"", false}}

	for _, c := range cases {
		if exp := IsPrivateIPv4Address(net.ParseIP(c.IP)); exp != c.IsPrivate {
			t.Errorf("IP address `%s` should be %s instead of %s", c.IP, c.IsPrivate, exp)
		}
	}
}

func TestFunctionGetFirstIPv4Address(t *testing.T) {
	cases := []struct {
		IPs    string
		Result string
	}{
		{" 3324.53434.654.32", ""},
		{"this is a text", ""},
		{" 192.168.0.1 ", ""},                         // In the private range, it'll be empty
		{" 192.142.0.1 ", "192.142.0.1"},              // Use the only one
		{" 192.168.0.1, 192.142.0.1 ", "192.142.0.1"}, // Use the first good one
		{" 192.142.0.1, 192.168.0.1 ", "192.142.0.1"}, // Use the first good one
		{" 2001:0db8:0a0b:12f0:0000:0000:0000:0001, 192.168.0.1, 192.142.0.1 ", "192.142.0.1"},
		{" 10.38.135.210 ,  68.111.195.103,  ", "68.111.195.103"},
	}

	for _, c := range cases {
		if exp := GetFirstIPv4Address(c.IPs); exp != c.Result {
			t.Errorf("Wrong IP address classification, it should be `%s` instead of `%s`", c.Result, exp)
		}
	}
}

func TestFunctionGetIPv4Address(t *testing.T) {
	cases := []struct {
		Header map[string]string
		Result string
	}{
		{map[string]string{}, ""},
		{map[string]string{X_FORWARDED_FOR: " 192.168.0.1 "}, ""},
		{map[string]string{X_FORWARDED_FOR: "this is a text"}, ""},
		{map[string]string{X_FORWARDED_FOR: "2001:0db8:0a0b:12f0:0000:0000:0000:0001"}, ""},
		{map[string]string{X_FORWARDED_FOR: " 192.142.0.1 "}, "192.142.0.1"},
		{map[string]string{X_FORWARDED_FOR: " 192.168.0.1, 192.142.0.1 "}, "192.142.0.1"},
		{map[string]string{X_FORWARDED_FOR: " 192.142.0.1, 192.168.0.1 "}, "192.142.0.1"},
		{map[string]string{X_FORWARDED_FOR: " 192.168.0.1 ", X_REAL_IP: " 192.142.0.1 "}, "192.142.0.1"},
		{map[string]string{X_FORWARDED_FOR: " 192.142.0.1 ", X_REAL_IP: " 192.168.0.1 "}, "192.142.0.1"},
	}

	for _, c := range cases {
		r, _ := http.NewRequest("POST", "/api/sth", nil)
		for k, v := range c.Header {
			r.Header.Set(k, v)
		}
		if exp := GetIPv4Address(r); exp != c.Result {
			t.Errorf("Wrong IP address classification, it should be `%s` instead of `%s`", c.Result, exp)
		}
	}
}
