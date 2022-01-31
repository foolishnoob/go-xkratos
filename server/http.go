package server

import (
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(serverConf *config.Server, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			validate.Validator(),
		),
	}
	if serverConf.GetHttp().GetNetwork() != "" {
		opts = append(opts, http.Network(serverConf.GetHttp().GetNetwork()))
	}
	if serverConf.GetHttp().GetAddr() != "" {
		opts = append(opts, http.Address(serverConf.GetHttp().GetAddr()))
	}
	if serverConf.GetHttp().GetTimeout() != nil {
		opts = append(opts, http.Timeout(serverConf.GetHttp().GetTimeout().AsDuration()))
	}
	return http.NewServer(opts...)
}
