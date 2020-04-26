package lb

import (
	"google.golang.org/grpc/resolver"
)

type etcdResolver struct{}

func (c *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
	panic("implement me")
}

func (c *etcdResolver) Close() {
	panic("implement me")
}

type consulResolver struct{}

func (c *consulResolver) ResolveNow(options resolver.ResolveNowOptions) {
	panic("implement me")
}

func (c *consulResolver) Close() {
	panic("implement me")
}
