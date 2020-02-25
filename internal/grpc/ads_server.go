// Copyright © 2020 VMware
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"google.golang.org/grpc"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	server "github.com/envoyproxy/go-control-plane/pkg/server"
)

// NewServer returns a *grpc.Server which responds to the Envoy v2 xDS gRPC API.
func NewServer(s server.Server, registry *prometheus.Registry, opts ...grpc.ServerOption) *grpc.Server {
	metrics := grpc_prometheus.NewServerMetrics()
	registry.MustRegister(metrics)
	opts = append(opts, grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()))
	g := grpc.NewServer(opts...)
	discovery.RegisterAggregatedDiscoveryServiceServer(g, s)
	discovery.RegisterSecretDiscoveryServiceServer(g, s)
	v2.RegisterClusterDiscoveryServiceServer(g, s)
	v2.RegisterListenerDiscoveryServiceServer(g, s)
	v2.RegisterRouteDiscoveryServiceServer(g, s)
	// v2.RegisterEndpointDiscoveryServiceServer(g, s)
	metrics.InitializeMetrics(g)
	return g
}
