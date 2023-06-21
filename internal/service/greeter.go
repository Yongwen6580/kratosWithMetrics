package service

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	v1 "helloworld/api/helloworld/v1"
	"helloworld/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc           *biz.GreeterUsecase
	requestCount prometheus.Counter
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	requestCount := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "greeter_requests_total",
		Help: "Total number of requests handled by the greeter service.",
	})

	prometheus.MustRegister(requestCount)
	return &GreeterService{uc: uc, requestCount: requestCount}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.requestCount.Inc()
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
}
