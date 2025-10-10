package claimsctx

import (
	"context"

	"github.com/WithSoull/platform_common/pkg/contextx"
)

const (
	UserEmailKey contextx.CtxKey = "user_email"
	UserIDKey    contextx.CtxKey = "user_id"
)

func InjectUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, UserEmailKey, email)
}

func InjectUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func ExtractUserEmail(ctx context.Context) (string, bool) {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email, true
	}
	return "unknown", false
}

func ExtractUserID(ctx context.Context) (int64, bool) {
	if userID, ok := ctx.Value(UserIDKey).(int64); ok {
		return userID, true
	}
	return 0, false
}
