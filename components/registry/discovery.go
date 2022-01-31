package registry

import (
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/foolishnoob/go-xkratos/util/xdebug"
	kEtcd "github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/v2/registry"
	etcd "go.etcd.io/etcd/client/v3"
	googleGrpc "google.golang.org/grpc"
	"time"
)

type Discovery interface {
	GetInterface() registry.Discovery
}

type discovery struct {
	discoveryClient registry.Discovery
}

func (d *discovery) GetInterface() registry.Discovery {
	return d.discoveryClient
}

func NewDiscovery(conf *config.BootConfig) Discovery {
	if 0 >= len(conf.GetEtcd().GetEndpoints()) {
		return nil
	}
	client, err := etcd.New(
		etcd.Config{
			Endpoints:   conf.GetEtcd().GetEndpoints(),
			DialTimeout: time.Second,
			DialOptions: []googleGrpc.DialOption{googleGrpc.WithBlock()},
		})
	xdebug.IfPanic(err)
	r := kEtcd.New(client)
	return &discovery{
		discoveryClient: r,
	}
}
