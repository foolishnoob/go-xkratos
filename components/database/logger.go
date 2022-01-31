package database

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
	"time"
)

type xlog struct {
	logger.Writer
	l logger.Interface
}

func newXlog() *xlog {
	return &xlog{
		l: logger.New(new(xlog), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      true,
		}),
	}
}

func (xl *xlog) getLogger() logger.Interface {
	return xl.l
}

//gorm使用kratos的日志（而不是使用gorm自带的日志）
func (xl *xlog) Printf(format string, v ...interface{}) {
	var level klog.Level
	_ = klog.GetLogger().Log(level, v)
}
