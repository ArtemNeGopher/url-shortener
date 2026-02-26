package logger

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GetLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	log = log.With(slog.String("context", "grpc-middleware"))
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		log.Debug(
			"request started",
			slog.Time("time", start),
			slog.String("method", info.FullMethod),
		)

		// Вызов обработчика
		resp, err := handler(ctx, req)

		// Логируем результат
		duration := time.Since(start)
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(err.Error(),
				slog.String("method", info.FullMethod),
				slog.String("code", fmt.Sprintf("%i", st.Code())),
				slog.String("duration", fmt.Sprintf("%vms", duration.Milliseconds())),
			)
		} else {
			log.Info(
				"Success",
				slog.String("method", info.FullMethod),
				slog.String("duration", fmt.Sprintf("%vms", duration.Milliseconds())),
			)
		}

		return resp, err
	}
}
