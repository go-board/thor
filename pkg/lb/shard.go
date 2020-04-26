package lb

import (
	"math/rand"

	"github.com/go-board/x-go/xhash/ring"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

func init() {
	balancer.Register(base.NewBalancerBuilderV2("shard", shardPickerBuilder{}, base.Config{HealthCheck: true}))
}

type shardPickerBuilder struct{}

func (shardPickerBuilder) Build(info base.PickerBuildInfo) balancer.V2Picker {
	picker := &shardPicker{subConns: map[string]balancer.SubConn{}, m: ring.New(50, nil)}
	for sc, scInfo := range info.ReadySCs {
		picker.m.Add(scInfo.Address.Addr)
		picker.subConns[scInfo.Address.Addr] = sc
		picker.remotes = append(picker.remotes, scInfo.Address.Addr)
	}
	return picker
}

type shardPicker struct {
	key      string
	subConns map[string]balancer.SubConn
	remotes  []string
	m        *ring.Map
}

func (s *shardPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	md, ok := metadata.FromOutgoingContext(info.Ctx)
	if !ok {
		return balancer.PickResult{SubConn: s.subConns[s.remotes[rand.Intn(len(s.remotes))]]}, nil
	}
	values := md.Get(s.key)
	if len(values) == 0 {
		return balancer.PickResult{SubConn: s.subConns[s.remotes[rand.Intn(len(s.remotes))]]}, nil
	}
	key := s.m.Get(values[0])
	return balancer.PickResult{
		SubConn: s.subConns[key],
	}, nil
}
