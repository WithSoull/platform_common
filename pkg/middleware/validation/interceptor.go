package validationInterceptor

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
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		res, err = handler(ctx, req)
		if nil == err {
			return res, nil
		}

		badErr := false

		switch {
		case sys.IsCommonError(err):
			commEr := sys.GetCommonError(err)
			code := toGRPCCode(commEr.Code())

			err = status.Error(code, commEr.Error())

		case validate.IsValidationError(err):
			err = status.Error(grpcCodes.InvalidArgument, err.Error())

		default:
			badErr = true
			var se GRPCStatusInterface
			if errors.As(err, &se) {
				return nil, se.GRPCStatus().Err()
			} else {
				if errors.Is(err, context.DeadlineExceeded) {
					err = status.Error(grpcCodes.DeadlineExceeded, err.Error())
				} else if errors.Is(err, context.Canceled) {
					err = status.Error(grpcCodes.Canceled, err.Error())
				} else {
					err = status.Error(grpcCodes.Internal, "internal error")
				}
			}
		}

		if badErr {
			logger.Error(ctx, "error interceptor hanlde error", zap.Error(err))
		} else {
			logger.Info(ctx, "error interceptor hanlde error", zap.Error(err))
		}

		return res, err
	}
}
