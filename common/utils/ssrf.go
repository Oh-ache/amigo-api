package utils

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ErrSSRFBlocked 表示 URL 被 SSRF 防护拦截
var ErrSSRFBlocked = fmt.Errorf("URL blocked by SSRF protection")

// 私有/保留 IP 网段列表
var privateNets = []net.IPNet{
	// IPv4 私有地址
	parseCIDR("10.0.0.0/8"),
	parseCIDR("172.16.0.0/12"),
	parseCIDR("192.168.0.0/16"),
	// IPv4 回环地址
	parseCIDR("127.0.0.0/8"),
	// IPv4 链路本地
	parseCIDR("169.254.0.0/16"),
	// IPv4 组播
	parseCIDR("224.0.0.0/4"),
	// IPv4 保留
	parseCIDR("240.0.0.0/4"),
	// IPv6 回环
	parseCIDR("::1/128"),
	// IPv6 链路本地
	parseCIDR("fe80::/10"),
	// IPv6 唯一本地
	parseCIDR("fc00::/7"),
}

func parseCIDR(s string) net.IPNet {
	_, n, _ := net.ParseCIDR(s)
	return *n
}

// ValidateURL 检查 URL 是否为安全的公网地址，防止 SSRF 攻击
// 1. 只允许 http/https 协议
// 2. 禁止访问内网/私有 IP
// 3. 禁止使用 IP 地址访问（可选，默认允许公网 IP）
// 4. 禁止 DNS 重绑定到内网地址
func ValidateURL(rawURL string) error {
	// 解析 URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// 只允许 http 和 https 协议
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("%w: scheme %q is not allowed, only http/https", ErrSSRFBlocked, u.Scheme)
	}

	// 提取主机名（去掉端口）
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("%w: empty host", ErrSSRFBlocked)
	}

	// 解析主机名为 IP 地址
	ip := net.ParseIP(host)
	if ip != nil {
		// 直接使用 IP 地址访问，检查是否为私有地址
		if isPrivateIP(ip) {
			return fmt.Errorf("%w: private IP address %s is not allowed", ErrSSRFBlocked, ip)
		}
		return nil
	}

	// 域名情况：DNS 解析后检查所有 IP 地址
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("failed to resolve host %s: %w", host, err)
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("%w: host %s resolves to private IP %s", ErrSSRFBlocked, host, ip)
		}
	}

	return nil
}

// isPrivateIP 检查 IP 是否属于私有/保留网段
func isPrivateIP(ip net.IP) bool {
	// 检查是否为未指定地址 (0.0.0.0)
	if ip.IsUnspecified() {
		return true
	}
	// 检查是否为回环地址
	if ip.IsLoopback() {
		return true
	}
	// 检查是否为链路本地地址
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	// 检查是否在私有网段中
	for _, n := range privateNets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
