package netutil

import (
	"os"
	"strings"
	//	"path/filepath"
	"context"
	"errors"
	"fmt"
	. "github.com/gemsnote/gemsnote/app/lea"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// net的util

// toPath 文件保存的目录
// 默认是/tmp
// 返回文件的完整目录
func WriteUrl(url string, toPath string) (length int64, newFilename, path string, ok bool) {
	if url == "" {
		return
	}
	content, err := GetContent(url)
	if err != nil {
		return
	}

	length = int64(len(content))

	// a.html?a=a11&xxx
	ext, valid := RasterImageExtension(content)
	if !valid {
		return
	}
	if toPath == "" {
		toPath = "/tmp"
	}
	//	dir := filepath.Dir(toPath)
	newFilename = NewGuid() + ext
	fullPath := toPath + "/" + newFilename

	// 写到文件中
	file, err := os.Create(fullPath)
	if err != nil {
		return
	}
	_, err = file.Write(content)
	closeErr := file.Close()
	if err != nil || closeErr != nil {
		os.Remove(fullPath)
		return
	}

	path = fullPath
	ok = true
	return
}

// 得到内容
func GetContent(rawURL string) ([]byte, error) {
	if err := validateRemoteURL(rawURL); err != nil {
		return nil, err
	}
	transport := &http.Transport{DialContext: publicDialContext, TLSHandshakeTimeout: 10 * time.Second, ResponseHeaderTimeout: 15 * time.Second}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: 30 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("too many redirects")
		}
		return validateRemoteURL(req.URL.String())
	}}
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download status %d", resp.StatusCode)
	}
	if resp.ContentLength > maxRemoteImageBytes {
		return nil, errors.New("image too large")
	}
	content, err := io.ReadAll(io.LimitReader(resp.Body, maxRemoteImageBytes+1))
	if err != nil {
		return nil, err
	}
	if len(content) > maxRemoteImageBytes {
		return nil, errors.New("image too large")
	}
	return content, nil
}

const maxRemoteImageBytes = 32 << 20

func validateRemoteURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		return errors.New("invalid remote image URL")
	}
	return nil
}

func publicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip.To4() == nil && ip[0]&0xe0 != 0x20 {
		return false
	}
	// Non-public ranges not covered by net.IP.IsPrivate, including shared
	// address space and IPv6 translation/tunnelling prefixes.
	for _, cidr := range []string{"0.0.0.0/8", "100.64.0.0/10", "192.0.0.0/24", "192.0.2.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "240.0.0.0/4", "64:ff9b::/96", "64:ff9b:1::/48", "100::/64", "2001::/23", "2001:db8::/32", "fec0::/10", "2002::/16", "3fff::/20"} {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return false
		}
	}
	return true
}

func publicDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, errors.New("no remote address")
	}
	for _, address := range ips {
		if !publicIP(address.IP) || address.Zone != "" {
			return nil, errors.New("non-public remote address")
		}
	}
	var last error
	for _, ip := range ips {
		// Dial the validated address, never resolve the hostname a second time.
		conn, err := (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ip.IP.String(), port))
		if err == nil {
			return conn, nil
		}
		last = err
	}
	return nil, last
}

// 将url ?, #后面的字符串去掉
func trimQueryParams(url string) string {
	pos := strings.Index(url, "?")
	if pos != -1 {
		url = Substr(url, 0, pos)
	}

	pos = strings.Index(url, "#")
	if pos != -1 {
		url = Substr(url, 0, pos)
	}

	pos = strings.Index(url, "!")
	if pos != -1 {
		url = Substr(url, 0, pos)
	}
	return url
}

// 通过domain得到ip
func GetIpFromDomain(domain string) string {
	ip, _ := net.LookupIP(domain)
	if ip != nil && len(ip) > 0 {
		return ip[0].String()
	}
	return ""
}
