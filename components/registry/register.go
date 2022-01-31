package registry

import (
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/foolishnoob/go-xkratos/util/xdebug"
	etcdRegistry "github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/v2/registry"
	etcd "go.etcd.io/etcd/client/v3"
	googleGrpc "google.golang.org/grpc"
	"time"
)

type Register interface {
	GetInterface() registry.Registrar
}

type registrar struct {
	registrar registry.Registrar
}

func (r *registrar) GetInterface() registry.Registrar {
	return r.registrar
}

func NewRegistrar(conf *config.BootConfig) Register {
	//@todo close etcd connection
	if 0 >= len(conf.GetEtcd().GetEndpoints()) {
		return nil
	}
	client, err := etcd.New(etcd.Config{
		Endpoints:   conf.GetEtcd().GetEndpoints(),
		DialTimeout: 2 * time.Second, //@todo read by config
		DialOptions: []googleGrpc.DialOption{googleGrpc.WithBlock()},
	})
	xdebug.IfPanic(err)
	return &registrar{
		registrar: etcdRegistry.New(client),
	}
}
