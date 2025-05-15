package server

import (
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"github.com/hobbyGG/kmall/review-service/internal/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistar)

// 服务启动时创建grpc服务器
func NewRegistar(config *conf.Registry) registry.Registrar {
	// new consul client
	c := api.DefaultConfig()
	c.Address = config.Consul.Addr
	c.Scheme = config.Consul.Scheme
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// new reg with consul client
	reg := consul.New(client, consul.WithHealthCheck(true))
	return reg
}
