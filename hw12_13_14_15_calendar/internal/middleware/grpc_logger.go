package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/kpechenenko/hw12_13_14_15_calendar/calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GrpcGatewayLoggingRequest логирует выполнение запросов для grpc-gateway.
func GrpcGatewayLoggingRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	md, ok := metadata.FromIncomingContext(ctx)
	// todo можно напечатать все хэдеры и посмотреть что они в себе несут
	if !ok {
		return resp, err
	}
	logger.Infof(
		"ip: %s, method: %s, error: %d, latency: %s, userAgent: %s",
		getClientIP(md),
		info.FullMethod,
		err,
		time.Since(start),
		getUserAgent(md),
	)
	return resp, err
}

func getClientIP(md metadata.MD) string {
	xForwardFor := md.Get("x-forwarded-for")
	if len(xForwardFor) > 0 && xForwardFor[0] != "" {
		ips := strings.Split(xForwardFor[0], ",")
		if len(ips) > 0 {
			clientIp := ips[0]
			return clientIp
		}
	}
	return ""
}

func getUserAgent(md metadata.MD) string {
	userAgent := md.Get("user-agent")
	if len(userAgent) == 0 {
		return ""
	}
	return userAgent[0]
}
