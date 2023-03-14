package response

type Event struct {
	Id string
}

type OutCome struct {
	Imei string
	Data struct {
		Message, Status  string
		Events           []Event
		Code, StatusCode int
		Fields           struct {
			Fp, Fd, Fn, DateFisal string
		}
	}
}

func (out *OutCome) SetEventId(ev []string) {
	out.Data.Events = make([]Event, 0, 1)
	for _, s := range ev {
		out.Data.Events = append(out.Data.Events, Event{Id: s})
	}
}
