package registry

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/go-board/x-go/xnet"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-hclog"
)

type consulRegistry struct{ client *api.Client }

type consulWatcher struct {
	plan        *watch.Plan
	registry    *consulRegistry
	watchResult chan *WatchResult
}

func (r *consulRegistry) watch() (*consulWatcher, error) {
	wp, err := watch.Parse(nil)
	if err != nil {
		return nil, err
	}
	wp.Handler = r.watchHandler
	go func() { err = wp.RunWithClientAndHclog(r.client, hclog.L()) }()
	return &consulWatcher{
		plan:        wp,
		registry:    r,
		watchResult: make(chan *WatchResult, 1),
	}, err
}

func (r *consulRegistry) watchHandler(id uint64, data interface{}) {
	services, ok := data.([]*api.ServiceEntry)
	if !ok {
		return
	}
	for _, service := range services {
		print(service)
	}
}

func (r *consulRegistry) GetService(ctx context.Context, name string) ([]*Service, error) {
	return nil, nil
}

func (r *consulRegistry) Register(ctx context.Context, service Service) error {
	address, err := xnet.PrivateAddress()
	if err != nil {
		return err
	}
	_, port, err := net.SplitHostPort(service.ServiceAddr)
	if err != nil {
		return err
	}
	intPort, _ := strconv.ParseInt(port, 10, 64)
	srv := &api.AgentServiceRegistration{
		Address:   address,
		ID:        service.ServiceID,
		Name:      service.ServiceName,
		Namespace: service.Namespace,
		Meta:      service.Metadata,
		Port:      int(intPort),
		Check: &api.AgentServiceCheck{
			TTL: strconv.FormatInt(int64(time.Second)*10, 10),
			// DeregisterCriticalServiceAfter: strconv.FormatInt(int64(time.Second)*5, 10),
		},
	}
	return r.client.Agent().ServiceRegister(srv)
}

func (r *consulRegistry) Deregister(ctx context.Context, service Service) error {
	return r.client.Agent().ServiceDeregister(service.ServiceID)
}

func (r *consulRegistry) Watch(ctx context.Context, name string) (Watcher, error) { return r.watch() }

func (c *consulWatcher) Next() (*WatchResult, error) {
	wr, ok := <-c.watchResult
	if !ok {
		return nil, errors.New("err: retrieve watch result from closed channel")
	}
	return wr, nil
}

func (c *consulWatcher) Shutdown() { c.plan.Stop() }
