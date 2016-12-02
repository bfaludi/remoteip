package remoteip

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

const X_FORWARDED_FOR = "X-Forwarded-For"
const X_REAL_IP = "X-Real-Ip"

type PrivateIPv4AddressRange struct {
	Start net.IP
	End   net.IP
}

var PrivateIPv4AddressRanges = []*PrivateIPv4AddressRange{
	{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
	{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
	{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
}

func (p *PrivateIPv4AddressRange) Contains(IP net.IP) bool {
	return bytes.Compare(IP, p.Start) >= 0 && bytes.Compare(IP, p.End) <= 0
}

func IsPrivateIPv4Address(IP net.IP) bool {
	for _, r := range PrivateIPv4AddressRanges {
		if r.Contains(IP) {
			return true
		}
	}
	return false
}

func IsIPv4Address(IP net.IP) bool {
	return IP.To4() != nil
}

func GetFirstIPv4Address(addresses string) string {
	for _, address := range strings.Split(addresses, ",") {
		IP := net.ParseIP(strings.TrimSpace(address))
		if IP == nil || IsIPv4Address(IP) == false || IsPrivateIPv4Address(IP) == true {
			continue
		}
		return IP.String()
	}
	return ""
}

func GetIPv4Address(r *http.Request) string {
	for _, addresses := range []string{
		r.Header.Get(X_FORWARDED_FOR),
		r.Header.Get(X_REAL_IP),
		r.RemoteAddr,
	} {
		if IP := GetFirstIPv4Address(addresses); IP != "" {
			return IP
		}
	}
	return ""
}
