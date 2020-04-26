package registry

import (
	"context"
)

type Service struct {
	ServiceName string
	ServiceID   string
	ServiceAddr string
	Metadata    map[string]string
}

type Registry interface {
	GetService(ctx context.Context, name string) ([]*Service, error)
	Register(ctx context.Context, service Service) error
	Deregister(ctx context.Context, service Service) error
}

type etcdRegistry struct {
}

func (r *etcdRegistry) GetService(ctx context.Context, name string) ([]*Service, error) {
	return nil, nil
}

func (r *etcdRegistry) Register(ctx context.Context, service Service) error {
	return nil
}

func (r *etcdRegistry) Deregister(ctx context.Context, service Service) error {
	return nil
}

type consulRegistry struct{}

func (c *consulRegistry) GetService(ctx context.Context, name string) ([]*Service, error) {
	return nil, nil
}

func (c *consulRegistry) Register(ctx context.Context, service Service) error {
	return nil
}

func (c *consulRegistry) Deregister(ctx context.Context, service Service) error {
	return nil
}
