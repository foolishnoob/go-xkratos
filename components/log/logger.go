package log

import (
	"github.com/foolishnoob/go-xkratos/config"
	"github.com/foolishnoob/go-xkratos/util/xdebug"
	"github.com/go-kratos/kratos/contrib/log/aliyun/v2"
	kLog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/tidwall/gjson"
	"go.uber.org/dig"
	"os"
)

func Inject(container *dig.Container) {
	_ = container.Provide(NewLogger)
	var err = container.Invoke(func(l kLog.Logger) {
		kLog.SetLogger(l)
	})
	xdebug.IfPanic(err)
}

func NewLogger(conf *config.BootConfig) kLog.Logger {
	//如果配置了阿里云日志，则使用阿里云日志；否则使用默认日志。有需要可以添加其他平台的日志对象。
	var l = kLog.NewStdLogger(os.Stdout)
	switch conf.GetLogger().GetPlatform() {
	case "aliyun":
		l = newAliyunLogger(conf)
	}
	return kLog.With(l,
		"ts", kLog.DefaultTimestamp,
		"caller", kLog.DefaultCaller,
		"service.id", conf.GetService().GetId(),
		"service.name", conf.GetService().GetName(),
		"service.version", conf.GetService().GetVersion(),
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
}

func newAliyunLogger(conf *config.BootConfig) kLog.Logger {
	var setting = conf.GetLogger().GetSetting()
	var options []aliyun.Option
	if endpoint := gjson.GetBytes([]byte(setting), "Endpoint").String(); "" != endpoint {
		options = append(options, aliyun.WithEndpoint(endpoint))
	}
	if accessKey := gjson.GetBytes([]byte(setting), "AccessKey").String(); "" != accessKey {
		options = append(options, aliyun.WithAccessKey(accessKey))
	}
	if accessSecret := gjson.GetBytes([]byte(setting), "AccessSecret").String(); "" != accessSecret {
		options = append(options, aliyun.WithAccessSecret(accessSecret))
	}
	return aliyun.NewAliyunLog(options...)
}
