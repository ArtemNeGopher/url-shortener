package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GetLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
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
			log.Error("Error",
				slog.String("method", info.FullMethod),
				slog.String("code", fmt.Sprintf("%v", st.Code())),
				slog.Duration("duration", duration),
			)
		} else {
			log.Info(
				"",
				slog.String("method", info.FullMethod),
				slog.Duration("duration", duration),
			)
		}

		return resp, err
	}
}
