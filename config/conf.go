package config

import (
	"flag"
	"github.com/foolishnoob/go-xkratos/util/xdebug"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func Inject(container *dig.Container) {
	_ = container.Provide(bootConfig)
	_ = container.Provide(serviceConfig)
	_ = container.Provide(serverConfig)
	_ = container.Provide(databaseConfig)
	_ = container.Provide(redisConfig)
	_ = container.Provide(etcdConfig)
	_ = container.Provide(traceConfig)

	var env string
	var err = container.Invoke(func(bootConfig *BootConfig) {
		env = bootConfig.GetEnvironment()
	})
	if nil == err && "" == env {
		err = errors.New("environment must be set!")
	}
	xdebug.IfPanic(err)
}

func bootConfig() *BootConfig {
	//@todo 目前使用yaml配置，后续接入apollo、etcd配置
	return formFlag()
}

func formFlag() *BootConfig {
	var flagConf string
	flag.StringVar(&flagConf, "conf", "", "config path, eg: -conf config.yaml")
	flag.Parse()
	conf := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	if err := conf.Load(); err != nil {
		panic(err)
	}
	var bootConfig *BootConfig
	if err := conf.Scan(&bootConfig); err != nil {
		panic(err)
	}
	return bootConfig
}

func serviceConfig(conf *BootConfig) *Service {
	var serviceConf = conf.GetService()
	if "" == serviceConf.GetId() {
		if id, err := uuid.NewUUID(); err == nil {
			serviceConf.Id = id.String()
		}
	}
	return serviceConf
}
func serverConfig(conf *BootConfig) *Server {
	return conf.GetServer()
}
func databaseConfig(conf *BootConfig) *Database {
	return conf.GetDatabase()
}
func redisConfig(conf *BootConfig) *Redis {
	return conf.GetRedis()
}
func etcdConfig(conf *BootConfig) *Etcd {
	return conf.GetEtcd()
}
func traceConfig(conf *BootConfig) *Trace {
	return conf.GetTrace()
}
