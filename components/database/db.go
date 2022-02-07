package database

import (
	"github.com/foolishnoob/go-xkratos/components/tracer"
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/foolishnoob/go-xkratos/util/xdebug"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database interface {
	GetInstance() *gorm.DB
}

type database struct {
	Database
	db *gorm.DB
}

func (d *database) GetInstance() *gorm.DB {
	return d.db
}

func NewDatabase(conf *config.BootConfig, tr tracer.Tracer) Database {
	if "" == conf.GetDatabase().GetDsn() {
		return nil
	}
	var l logger.Interface
	if "local" == conf.GetEnvironment() {
		l = logger.Default.LogMode(logger.Info)
	} else {
		l = newXlog().getLogger().LogMode(logger.Info)
	}
	var conn, err = gorm.Open(mysql.Open(conf.GetDatabase().GetDsn()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // 表名前缀，统一在业务程序里面设置
			SingularTable: true, // 使用单数表名，统一在业务程序里面设置
		},
		Logger: l, //@todo 日志等级
	})
	xdebug.IfPanic(err)
	_ = conn.Use(otelgorm.NewPlugin(otelgorm.WithTracerProvider(tr.GetInstance())))
	return &database{db: conn}
}
