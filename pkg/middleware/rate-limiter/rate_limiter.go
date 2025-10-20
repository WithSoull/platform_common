package rateLimiterInterceptor

import (
	"context"

	rateLimiter "github.com/WithSoull/platform_common/pkg/rate-limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Error(codes.ResourceExhausted, "to many requests")
	}

	return handler(ctx, req)
}
