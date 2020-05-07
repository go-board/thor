package thor

import (
	"log"
	"os"
	"time"

	"github.com/go-board/x-go/metadata"
	"gopkg.in/yaml.v2"

	"github.com/go-board/thor/pkg/logger"
	"github.com/go-board/thor/pkg/metric"
	"github.com/go-board/thor/pkg/trace"
)

// Options is the thor application global configuration.
type Options struct {
	Namespace      string            `yaml:"namespace"`
	ServiceName    string            `yaml:"service_name"`
	ServiceID      string            `yaml:"service_id"`
	ServiceVersion string            `yaml:"service_version"`
	Metadata       metadata.Metadata `yaml:"metadata"`
	Listeners      []ListenerOption  `yaml:"listeners"`
	Logger         LoggerOption      `yaml:"logger"`
	Trace          TraceOption       `yaml:"trace"`
	Registry       RegistryOption    `yaml:"registry"`
	Resilience     ResilienceOption  `yaml:"resilience"`
}

type ListenerOption struct {
	Type    ListenerType `yaml:"listener_type"`
	NetType string       `yaml:"listener_net_type"`
	Addr    string       `yaml:"listener_addr"`
}

type ListenerType int

func (l ListenerType) String() string {
	switch l {
	case ListenerTypeGRPC:
		return "GRPC"
	case ListenerTypeHTTP:
		return "HTTP"
	default:
		return "TCP"
	}
}

const (
	ListenerTypeTcp ListenerType = iota
	ListenerTypeGRPC
	ListenerTypeHTTP
)

type RegistryOption struct {
	RegistryType string        `yaml:"registry_type"`
	RegistryAddr string        `yaml:"registry_addr"`
	RegistryTTL  time.Duration `yaml:"registry_ttl"`
}

type RegistryType int

func (r RegistryType) String() string {
	switch r {
	case RegistryTypeEtcd:
		return "etcd"
	case RegistryTypeConsul:
		return "consul"
	case RegistryTypeK8s:
		return "k8s"
	case RegistryMdns:
		return "mdns"
	}
	return ""
}

const (
	RegistryTypeEtcd RegistryType = iota
	RegistryTypeConsul
	RegistryTypeK8s
	RegistryMdns
)

type ResilienceOption struct {
	RetryOnIdempotent bool `yaml:"retry_on_idempotent"`
}

type LoggerOption struct {
	Dir         string `yaml:"log_dir"`
	LevelFilter string `yaml:"log_level_filter"`
}

type TraceOption struct {
	TraceSampleType  string  `yaml:"trace_sample_type"`
	TraceSampleParam float64 `yaml:"trace_sample_param"`
}

var globalOptions = new(Options)

func OptionReader() Options { return *globalOptions }

// Initialize create the whole world of the current application.
func Initialize(options ...Option) {
	f, err := os.Open("env")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("env file exist, but open failed, %s\n", err)
		}
	}
	err = yaml.NewDecoder(f).Decode(globalOptions)
	if err != nil {
		log.Fatalf("parse env file failed, %s\n", err)
	}

	for _, option := range options {
		option(globalOptions)
	}

	logger.Initialize(globalOptions.Logger.Dir, globalOptions.Namespace, globalOptions.ServiceName, globalOptions.ServiceID)
	trace.Initialize(globalOptions.ServiceName, globalOptions.Trace.TraceSampleType, globalOptions.Trace.TraceSampleParam)
	metric.Initialize(globalOptions.Namespace, globalOptions.ServiceName, globalOptions.ServiceID, globalOptions.ServiceVersion)
}
