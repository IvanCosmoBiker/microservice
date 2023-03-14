package status

import (
	"net/http"
	"net/http/pprof"
)

type ServiceStatus struct {
	Address string
}

func New() *ServiceStatus {
	serviceApi := &ServiceStatus{
		Address: "188.225.18.140:8020",
	}
	return serviceApi
}

func (ss *ServiceStatus) InitApi() {
	r := http.NewServeMux()
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	http.ListenAndServe(ss.Address, r)
}
