package thor

import (
	"github.com/go-board/x-go/metadata"
)

type Option func(o *Options)

func Namespace(s string) Option {
	return func(o *Options) {
		o.Namespace = s
	}
}

func ServiceName(s string) Option {
	return func(o *Options) {
		o.ServiceName = s
	}
}

func ServiceID(s string) Option {
	return func(o *Options) {
		o.ServiceID = s
	}
}

func ServiceVersion(s string) Option {
	return func(o *Options) {
		o.ServiceVersion = s
	}
}

func Metadata(meta metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = meta
	}
}

func Listeners(listeners ...ListenerOption) Option {
	return func(o *Options) {
		o.Listeners = listeners
	}
}

func LogDir(s string) Option {
	return func(o *Options) {
		o.Logger.Dir = s
	}
}

func LogLevelFilter(l string) Option {
	return func(o *Options) {
		o.Logger.LevelFilter = l
	}
}
