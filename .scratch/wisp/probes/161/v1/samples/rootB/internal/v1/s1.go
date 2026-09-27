package v1

type svc struct{}

func (svc) worker() {}

func kick(s svc) {
	go s.worker()
}
