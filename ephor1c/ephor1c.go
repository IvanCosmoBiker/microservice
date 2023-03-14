package ephor1c

import (
	manager "ephorservices/ephor1c/manager"
	"flag"
	"log"

	"github.com/kardianos/service"
)

type Service_1c struct {
	State   service.Status
	Config  service.Config
	Manager *manager.Manager
}

func (sc *Service_1c) setStatus(status byte) {
	sc.State = status
}

func (sc *Service_1c) Start(s service.Service) error {
	if service.Interactive() {
		log.Print("Running in terminal.")
	} else {
		log.Print("Running under service manager.")
	}
	go sc.run(s)
	return nil
}

func (sc *Service_1c) run(s service.Service) error {
	manager.InitService(sc)
	return nil
}

func (sc *Service_1c) Stop(s service.Service) error {
	manager.StopService()
	log.Print("I'm Stopping!")
	return nil
}

func (sc *Service_1c) Restart(s service.Service) error {
	return nil
}

func (sc *Service_1c) Status(s service.Service) error {
	log.Print("I'm status!")
	return nil
}

func (sc *Service_1c) SetConfig() error {
	return nil
}

func (sc *Service_1c) SetDependencies() error {
	return nil
}

func Init() {
	prg := &Service_1c{}
	prg.setStatus(service.StateRunning)
	svcFlag := flag.String("start", "", "start programm.")
	flag.Parse()
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := service.Config{
		Name:         "EphorService1c",
		DisplayName:  "Ephor microservice for integartion 1c",
		Description:  "Microservice of Ephor company",
		Dependencies: []string{},
		Option:       options,
	}
	s, err := service.New(prg, &svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 10)
	prg.ReadConfig()
	prg.initLogger()
	log.Println("Load Logger...")
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	log.Println("Start Service_1c")
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}

/*
 Функция Run запускает микросервис и ставит состояние Idle
*/
func Run() {
	Init()
}
