package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"time"

	"amigo-api/app/gateway/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxBodySize = 10 << 20 // 10MB

type proxyEntry struct {
	prefix string
	proxy  *httputil.ReverseProxy
}

type Manager struct {
	entries []proxyEntry
}

func NewManager(rules []config.RouteRule) *Manager {
	// 按前缀长度降序排列，确保最长前缀优先匹配
	sort.Slice(rules, func(i, j int) bool {
		return len(rules[i].Prefix) > len(rules[j].Prefix)
	})

	entries := make([]proxyEntry, 0, len(rules))
	for _, rule := range rules {
		u, err := url.Parse(rule.Upstream)
		if err != nil {
			logx.Errorf("[gateway] invalid upstream for prefix %s: %v", rule.Prefix, err)
			continue
		}
		proxy := buildReverseProxy(u)
		entries = append(entries, proxyEntry{prefix: rule.Prefix, proxy: proxy})
		logx.Infof("[gateway] route: %s -> %s", rule.Prefix, rule.Upstream)
	}

	return &Manager{entries: entries}
}

func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 忽略 favicon 等非 API 请求
	if !strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// 限制请求体大小
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	// 找到匹配的上游
	for _, entry := range m.entries {
		if strings.HasPrefix(r.URL.Path, entry.prefix) {
			serveProxy(w, r, entry.proxy, entry.prefix)
			return
		}
	}

	http.NotFound(w, r)
}

func (m *Manager) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func serveProxy(w http.ResponseWriter, r *http.Request, proxy *httputil.ReverseProxy, upstream string) {
	// 注入 Request ID
	reqID := r.Header.Get("X-Request-Id")
	if reqID == "" {
		reqID = fmt.Sprintf("%d", time.Now().UnixNano())
		r.Header.Set("X-Request-Id", reqID)
	}

	start := time.Now()
	sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

	proxy.ServeHTTP(sw, r)

	duration := time.Since(start)
	logx.Infof("[gateway] %s %s -> %s %d %s", r.Method, r.URL.Path, upstream, sw.status, duration.Round(time.Millisecond))
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("hijack not supported")
}

func buildReverseProxy(u *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logx.Errorf("[gateway] proxy error for %s: %v", r.URL.Path, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"code":1,"msg":"upstream unavailable"}`))
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		// 确保上游返回的 Content-Type 不被覆盖
		return nil
	}

	return proxy
}
