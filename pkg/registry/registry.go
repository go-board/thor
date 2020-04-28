package registry

import (
	"context"
	"errors"
)

type Service struct {
	Namespace   string
	ServiceName string
	ServiceID   string
	ServiceAddr string
	Metadata    map[string]string
}

type WatchResult struct {
	Services []*Service
	Action   Action
	Error    error
}

type Action int

const (
	ActionAdd Action = iota
	ActionUpdate
	ActionRemove
)

type Registry interface {
	GetService(ctx context.Context, name string) ([]*Service, error)
	Register(ctx context.Context, service Service) error
	Deregister(ctx context.Context, service Service) error

	Watch(ctx context.Context, name string) (Watcher, error)
}

type Watcher interface {
	Next() (*WatchResult, error)
	Shutdown()
}

type etcdRegistry struct{}

func (r *etcdRegistry) GetService(ctx context.Context, name string) ([]*Service, error) {
	return nil, nil
}

func (r *etcdRegistry) Register(ctx context.Context, service Service) error {
	return nil
}

func (r *etcdRegistry) Deregister(ctx context.Context, service Service) error {
	return nil
}

func (r *etcdRegistry) Watch(ctx context.Context, name string) (Watcher, error) {
	return nil, errors.New("not implemented")
}
