package rewriteip

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// 初始化插件
func init() {
	plugin.Register("rewriteip", setup)
}

// setup 配置插件
func setup(c *caddy.Controller) error {
	c.Next() // 读取 "rewriteip" 关键字
	if c.NextArg() {
		return plugin.Error("rewriteip", c.ArgErr()) // 不接受额外参数
	}

	// 在 CoreDNS 处理链中注册插件
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return RewriteIP{Next: next}
	})

	return nil
}
