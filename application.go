package xkratos

import (
	"github.com/asaskevich/EventBus"
	"github.com/foolishnoob/go-xkratos/components/cache"
	"github.com/foolishnoob/go-xkratos/components/database"
	"github.com/foolishnoob/go-xkratos/components/log"
	"github.com/foolishnoob/go-xkratos/components/registry"
	"github.com/foolishnoob/go-xkratos/components/tracer"
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/foolishnoob/go-xkratos/server"
	"github.com/foolishnoob/go-xkratos/util/xdebug"
	"github.com/go-kratos/kratos/v2"
	kLog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/dig"
	"sync"
)

var (
	once sync.Once
	App  *application
)

type application struct {
	container *dig.Container
}

func Application() *application {
	initialize := func() {
		App = new(application)
		App.container = dig.New()

		config.Inject(App.container)
		log.Inject(App.container)
		tracer.Inject(App.container)

		_ = App.container.Provide(database.NewDatabase)
		_ = App.container.Provide(database.NewDatabase)
		_ = App.container.Provide(cache.NewCache)
		_ = App.container.Provide(registry.NewDiscovery)
		_ = App.container.Provide(registry.NewRegistrar)
		_ = App.container.Provide(server.NewHTTPServer)
		_ = App.container.Provide(server.NewGRPCServer)
		_ = App.container.Provide(EventBus.New)
	}
	once.Do(initialize)
	return App
}

func (app *application) Environment() string {
	var env string
	var err = app.container.Invoke(func(bootConfig *config.BootConfig) {
		env = bootConfig.GetEnvironment()
	})
	xdebug.IfPanic(err)
	return env
}

func (app *application) GetContainer() *dig.Container {
	return app.container
}

func (app *application) Logger() kLog.Logger {
	var l kLog.Logger
	var err = app.container.Invoke(func(obj kLog.Logger) {
		l = obj
	})
	xdebug.IfPanic(err)
	return l
}

func (app *application) Database() database.Database {
	var conn database.Database
	var err = app.container.Invoke(func(obj database.Database) {
		conn = obj
	})
	xdebug.IfPanic(err)
	return conn
}

func (app *application) Cache() cache.Cache {
	var conn cache.Cache
	var err = app.container.Invoke(func(obj cache.Cache) {
		conn = obj
	})
	xdebug.IfPanic(err)
	return conn
}

func (app *application) EventBus() EventBus.Bus {
	var eventBus EventBus.Bus
	var err = app.container.Invoke(func(obj EventBus.Bus) {
		eventBus = obj
	})
	xdebug.IfPanic(err)
	return eventBus
}

func (app *application) Run(httpServer *http.Server, rpcServer *grpc.Server) {
	var kApp *kratos.App
	var err = app.container.Invoke(func(rr registry.Register, serviceConf *config.Service, logger kLog.Logger) {
		var servers []transport.Server
		if httpServer != nil {
			servers = append(servers, httpServer)
		}
		if rpcServer != nil {
			servers = append(servers, rpcServer)
		}
		var options = []kratos.Option{
			kratos.ID(serviceConf.GetId()),
			kratos.Name(serviceConf.GetName()),
			kratos.Version(serviceConf.GetVersion()),
			kratos.Metadata(map[string]string{}),
			kratos.Logger(logger),
		}
		if 0 < len(servers) {
			options = append(options, kratos.Server(servers...))
		}
		if nil != rr.GetInterface() {
			options = append(options, kratos.Registrar(rr.GetInterface()))
		}
		kApp = kratos.New(options...)
	})
	xdebug.IfPanic(err)
	if err = kApp.Run(); err != nil {
		panic(err)
	}
}

func (app *application) Stop() {
	_ = app.container.Invoke(func(db database.Database) {
		if nil == db {
			return
		}
		if conn, err := db.GetInstance().DB(); nil != err {
			_ = conn.Close()
		}
	})
}
