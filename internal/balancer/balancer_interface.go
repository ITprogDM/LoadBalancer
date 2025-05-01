package balancer

type LoadBalancer interface {
	Next() string
}
