package database

import (
	"github.com/foolishnoob/go-xkratos/config"
	_ "github.com/pkg/errors"
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

func NewDatabase(conf *config.BootConfig) Database {
	if "" == conf.GetDatabase().GetDsn() {
		return nil
	}
	var conn, err = gorm.Open(mysql.Open(conf.GetDatabase().GetDsn()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // 表名前缀，统一在业务程序里面设置
			SingularTable: true, // 使用单数表名，统一在业务程序里面设置
		},
		Logger: newXlog().getLogger().LogMode(logger.Info), //@todo 日志等级
	})
	if nil != err {
		panic(err)
	}
	return &database{db: conn}
}
