package tokens

import (
	"context"
)

type TokenGenerator interface {
	GenerateAccessToken(context.Context, UserInfo) (string, error)
	GenerateRefreshToken(context.Context, UserInfo) (string, error)
	VerifyAccessToken(context.Context, string) (*UserClaims, error)
	VerifyRefreshToken(context.Context, string) (*UserClaims, error)
}
