package interceptor

import (
	"context"
	"testing"

	"gitee.com/cruvie/kk_kit/go/kk_jwt"
)

func TestAuth(t *testing.T) {
	_ = UnaryAuth(&AuthConfig{
		JWTChecker: jwtAuthFunc(&kk_jwt.ConfigJWT{}),
	})
	_ = StreamAuth(&AuthConfig{
		JWTChecker: jwtAuthFunc(&kk_jwt.ConfigJWT{}),
	})
}

func jwtAuthFunc(cfg *kk_jwt.ConfigJWT) JWTCheckFunc {
	return func(ctx context.Context, token string) (
		newCtx context.Context,
		needRefresh bool,
		err error,
	) {
		_, needRefresh, err = cfg.VerifyToken(token)
		if err != nil {
			return nil, false, err
		}

		return newCtx, needRefresh, nil
	}
}
