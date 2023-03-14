package control

import (
	service_status "ephorservices/ephorsale/control/api/state/status"
	service_transaction_state "ephorservices/ephorsale/control/api/state/transaction/service_http"
	"ephorservices/ephorsale/transport"
)

/*
Control module. Performs such functions as:
1) Monitors and reports the status of each service module
2) Receives commands from the control system
*/

const (
	State_Idle uint8 = iota
	State_Work
	State_Warning
	State_Error
)

type Control struct {
	Name  string
	State uint8
}

func New() (controlModule *Control) {
	controlModule = &Control{
		Name:  "Control",
		State: State_Idle,
	}
	return controlModule
}

func (c *Control) InitApi() {
	serviceHttpHandler := service_transaction_state.New()
	serviceHttpHandler.InitApi(transport.TransportManager.RequestManager)
	serviceHttpStatus := service_status.New()
	go serviceHttpStatus.InitApi()
}
