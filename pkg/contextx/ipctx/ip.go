package ipctx

import (
	"context"
	"net"

	"github.com/WithSoull/AuthService/internal/contextx"
	"google.golang.org/grpc/peer"
)

const IpKey contextx.CtxKey = "ip"

func InjectIp(ctx context.Context) context.Context {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ctx
	}

	ip, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return ctx
	}

	return context.WithValue(ctx, IpKey, ip)
}

func ExtractIP(ctx context.Context) (string, bool) {
	ip, ok := ctx.Value(IpKey).(string)
	if !ok {
		return "unknown", false
	}
	return ip, true
}
