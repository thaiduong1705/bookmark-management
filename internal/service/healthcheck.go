package service

//go:generate mockery --name=HealthCheck --filename healthcheck.go --outpkg healthcheck_svc_mocks
type HealthCheck interface {
	Check() (string, string)
}

type healthCheck struct {
	serviceName string
	instanceID  string
}

func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheck{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (h *healthCheck) Check() (string, string) {
	return h.serviceName, h.instanceID
}
