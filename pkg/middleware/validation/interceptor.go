package validation

import (
	"context"

	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCStatusInterface interface {
	GRPCStatus() *status.Status
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

func ErrorCodesInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		res, err = handler(ctx, req)
		if nil == err {
			return res, nil
		}

		switch {
		case sys.IsCommonError(err):
			commEr := sys.GetCommonError(err)
			code := toGRPCCode(commEr.Code())

			logger.Info(ctx, "error interceptor handle common error", zap.Error(err))
			return nil, status.Error(code, commEr.Error())

		case validate.IsValidationError(err):
			logger.Info(ctx, "error interceptor handle validation error", zap.Error(err))
			return nil, status.Error(grpcCodes.InvalidArgument, err.Error())
		default:
			var se GRPCStatusInterface
			if errors.As(err, &se) {
				return nil, se.GRPCStatus().Err()
			} else {
				logger.Error(ctx, "error interceptor handle validation error", zap.Error(err))
				if errors.Is(err, context.DeadlineExceeded) {
					return nil, status.Error(grpcCodes.DeadlineExceeded, err.Error())
				} else if errors.Is(err, context.Canceled) {
					return nil, status.Error(grpcCodes.Canceled, err.Error())
				} else {
					return nil, status.Error(grpcCodes.Internal, "internal error")
				}
			}
		}
	}
}
