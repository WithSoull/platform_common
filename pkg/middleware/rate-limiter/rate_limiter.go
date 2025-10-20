package rateLimiterInterceptor

import (
	"context"

	rateLimiter "github.com/WithSoull/platform_common/pkg/rate-limiter"
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
	"google.golang.org/grpc"
)

type RateLimiterInterceptor struct {
	rateLimiter *rateLimiter.TokenBucketLimiter
}

func NewRateLimiterInterceptor(ctx context.Context, cfg rateLimiter.RateLimiterConfig) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{
		rateLimiter: rateLimiter.NewTokenBucketLimiter(ctx, cfg),
	}
}

func (r *RateLimiterInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if !r.rateLimiter.Allow() {
		return nil, sys.NewCommonError("too many requests", codes.ResourceExhausted)
	}

	return handler(ctx, req)
}
