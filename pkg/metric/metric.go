package metric

import (
	"fmt"
	"sync"

	"github.com/go-board/x-go/xos"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	processCollector   prometheus.Collector
	goCollector        prometheus.Collector
	buildInfoCollector prometheus.Collector
	once               sync.Once
)

const duplicatedCollector = "duplicate metrics collector registration attempted"

func Initialize(namespace, serviceName, serviceId string, version string) {
	once.Do(func() {
		processCollector = prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})
		goCollector = prometheus.NewGoCollector()
		buildInfoCollector = prometheus.NewBuildInfoCollector()

		gatherer := prometheus.NewRegistry()
		registerer := prometheus.WrapRegistererWith(
			prometheus.Labels{
				"service_instance": xos.Hostname(),
				"service_id":       serviceId,
				"service_version":  version,
			}, gatherer,
		)
		registerer = prometheus.WrapRegistererWithPrefix(fmt.Sprintf("%s_%s_", namespace, serviceName), registerer)
		registerer.MustRegister(processCollector, goCollector, buildInfoCollector)

		prometheus.DefaultGatherer = gatherer
		prometheus.DefaultRegisterer = registerer
	})
}

func EnableSystemMetric() error {
	if err := prometheus.Register(processCollector); err != nil && err.Error() != duplicatedCollector {
		return err
	}
	if err := prometheus.Register(goCollector); err != nil && err.Error() != duplicatedCollector {
		return err
	}
	if err := prometheus.Register(buildInfoCollector); err != nil && err.Error() != duplicatedCollector {
		return err
	}
	return nil
}

func DisableSystemMetric() {
	prometheus.Unregister(processCollector)
	prometheus.Unregister(goCollector)
	prometheus.Unregister(buildInfoCollector)
}
