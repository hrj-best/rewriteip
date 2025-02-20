package rewriteip

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// RewriteIP 插件结构
type RewriteIP struct {
	Next plugin.Handler
}

// ServeDNS 处理 DNS 查询
func (ri RewriteIP) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// 记录 Unbound 返回的响应
	rec := dnstest.NewRecorder(w)
	rcode, err := ri.Next.ServeDNS(ctx, rec, r)
	if err != nil {
		return rcode, err
	}

	// 解析客户端 IP
	clientIP := net.ParseIP(state.IP())
	if clientIP == nil {
		return rcode, nil // 如果获取失败，直接返回
	}

	// 修改 A 记录
	for _, ans := range rec.Msg.Answer {
		if aRecord, ok := ans.(*dns.A); ok {
			aRecord.A = clientIP.To4() // 替换 A 记录
		}
	}

	// 发送修改后的响应
	w.WriteMsg(rec.Msg)
	return rcode, nil
}

// Name 返回插件名称
func (ri RewriteIP) Name() string {
	return "rewriteip"
}
