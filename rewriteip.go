package rewriteip

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
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

	// 调用下游插件（Unbound）
	rcode, err := ri.Next.ServeDNS(ctx, w, r)
	if err != nil {
		return rcode, err
	}

	// 获取客户端 IP
	clientIP := net.ParseIP(state.IP())

	// 遍历 DNS 响应，修改 A 记录
	for _, ans := range r.Answer {
		if aRecord, ok := ans.(*dns.A); ok {
			aRecord.A = clientIP.To4() // 替换为客户端 IP
		}
	}

	// 发送修改后的响应
	w.WriteMsg(r)
	return rcode, nil
}

// Name 返回插件名称
func (ri RewriteIP) Name() string {
	return "rewriteip"
}
