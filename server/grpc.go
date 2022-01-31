package server

import (
	"github.com/foolishnoob/go-xkratos/components/tracer"
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(serverConf *config.Server, tracer tracer.Tracer, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(tracing.WithTracerProvider(tracer.GetInstance())),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
		),
	}
	if serverConf.GetGrpc().GetNetwork() != "" {
		opts = append(opts, grpc.Network(serverConf.GetGrpc().GetNetwork()))
	}
	if serverConf.GetGrpc().GetAddr() != "" {
		opts = append(opts, grpc.Address(serverConf.GetGrpc().GetAddr()))
	}
	if serverConf.GetGrpc().GetTimeout() != nil {
		opts = append(opts, grpc.Timeout(serverConf.GetGrpc().GetTimeout().AsDuration()))
	}
	return grpc.NewServer(opts...)
}
