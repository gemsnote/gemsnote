package netutil

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestRemoteImagesCannotReachLocalNetwork(t *testing.T) {
	var calls atomic.Int32
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.Write([]byte("private")) }))
	defer s.Close()
	if _, err := GetContent(s.URL); err == nil {
		t.Fatal("loopback request accepted")
	}
	if calls.Load() != 0 {
		t.Fatal("request reached private service")
	}
	for _, ip := range []string{"127.0.0.1", "10.1.2.3", "172.16.0.1", "192.168.0.1", "169.254.169.254", "::1", "::ffff:127.0.0.1", "fc00::1", "fe80::1", "100.64.0.1", "0.0.0.0", "224.0.0.1", "2002:7f00:1::", "2001:db8::1", "fec0::1", "::127.0.0.1"} {
		if publicIP(net.ParseIP(ip)) {
			t.Errorf("non-public IP accepted: %s", ip)
		}
	}
	for _, ip := range []string{"8.8.8.8", "2606:4700:4700::1111"} {
		if !publicIP(net.ParseIP(ip)) {
			t.Errorf("public address rejected: %s", ip)
		}
	}
	for _, raw := range []string{"file:///etc/passwd", "ftp://example.test/a", "http://user:secret@example.test/", "http:///missing"} {
		if validateRemoteURL(raw) == nil {
			t.Errorf("unsafe URL accepted: %s", raw)
		}
	}
}
