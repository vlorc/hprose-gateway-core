package option

import (
	"context"
	"github.com/vlorc/hprose-gateway-core/invoker"
	"github.com/vlorc/hprose-gateway-core/plugin"
	"github.com/vlorc/hprose-gateway-core/router"
	"github.com/vlorc/hprose-gateway-core/source"
	"github.com/vlorc/hprose-gateway-core/watcher"
	"github.com/vlorc/hprose-gateway-types"
	"go.uber.org/zap"
)

type GatewayOption struct {
	Router   types.NamedRouter
	Context  context.Context
	Resolver types.NamedResolver
	Watcher  types.NamedWatcher
	Manager  types.SourceManger
	Named    types.NamedMode
	Balancer string
	Error    error
	Log      func() *zap.Logger
	Plugins  invoker.Invoker
	Prefix   string
}

func Resolver(resolver types.NamedResolver) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Resolver = resolver
	}
}

func Manager(manager types.SourceManger) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Manager = manager
	}
}

func Name(name string) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Prefix = name
	}
}

func Prefix(prefix string) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Prefix = prefix
	}
}

func WithValue(key string, val interface{}) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Context = context.WithValue(opt.Context, key, val)
	}
}

func Context(ctx context.Context, env ...interface{}) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		ctx = context.WithValue(ctx, "option", opt)
		for i, l := 0, len(env)/2; i < l; i++ {
			ctx = context.WithValue(ctx, env[i*2+0], env[i*2+1])
		}
		opt.Context = ctx
	}
}

func Error(err error) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Error = err
	}
}

func Logger(debug bool) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		var log *zap.Logger
		var err error
		if debug {
			log, err = zap.NewDevelopment()
		} else {
			log, err = zap.NewProduction()
		}
		if nil != err {
			panic(err)
		}
		opt.Log = func() *zap.Logger {
			return log
		}
	}
}

func Balancer(name string) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Balancer = name
	}
}

func Named(mode types.NamedMode) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		opt.Named = mode
	}
}

func Router(r ...types.NamedRouter) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		if len(r) > 0 {
			opt.Router = r[0]
		} else {
			opt.Router = router.NewNamedRouter(opt.Balancer, opt.Manager, opt.Named)
		}
	}
}

func Watcher(w ...types.NamedWatcher) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		var ww types.NamedWatcher
		if len(w) > 0 {
			ww = watcher.NewMultiWatcher(watcher.NewChannelWatcher(100),w...)
		} else {
			ww = watcher.EmptyWatcher{}
		}
		opt.Watcher = watcher.NewSnapshotWatcher(
			watcher.NewChannelWatcher(100),
			opt.Router,
			opt.Manager,
			source.NewErrorSource(opt.Error),
			ww)
	}
}

func pluginAppend(opt *GatewayOption, info types.Describe) {
	factory := plugin.Query(info.Name)
	opt.Log().Debug("Plugin", zap.String("name", info.Name), zap.Bool("query", nil != factory))
	if nil == factory {
		return
	}
	bean := factory.Instance(opt.Context, info.Param)
	opt.Log().Debug("Plugin", zap.String("name", info.Name), zap.Bool("instance", nil != bean))
	if nil == bean {
		return
	}
	opt.Plugins = append(opt.Plugins, bean)
}

func Plugin(name string, param map[string]string) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		pluginAppend(opt, types.Describe{Name: name, Param: param})
	}
}

func Plugins(args ...types.Describe) func(*GatewayOption) {
	return func(opt *GatewayOption) {
		for i := range args {
			pluginAppend(opt, args[i])
		}
	}
}
