package ephorsale

// main file for start microservice
import (
	config "ephorservices/config"
	manager_service "ephorservices/ephorsale/manager"
	logger "ephorservices/pkg/logger"
	"flag"
	"log"

	"github.com/kardianos/service"
)

// Program structures.
// Define Start and Stop methods.
type SaleService struct {
	Manager    *manager_service.Manager
	State      int
	Logger     *logger.Logger
	Config     service.Config
	ConfigFile *config.Config
}

func (se *SaleService) Start(s service.Service) error {
	if service.Interactive() {
		log.Print("Running in terminal.")
	} else {
		log.Print("Running under service manager.")
	}
	go se.run(s)
	return nil
}

func (se *SaleService) run(s service.Service) error {
	se.Manager.InitService()
	return nil
}

func (se *SaleService) Stop(s service.Service) error {
	se.Manager.StopService()
	log.Print("I'm Stopping!")
	return nil
}

func (se *SaleService) Status(s service.Service) error {
	log.Print("I'm status!")
	return nil
}

func (se *SaleService) InitLogger() {
	var err error
	se.Logger, err = logger.New(se.ConfigFile.Log, se.ConfigFile.PrefixLog, se.ConfigFile.LogAvalable, se.ConfigFile.LogEnable)
	if err != nil {
		log.Println(err.Error())
	}
	se.Logger.Print("Load Logger...")
}

func (se *SaleService) ReadConfig() {
	config := &config.Config{}
	config.Load()
	se.ConfigFile = config
	log.Println("Load config...")
}

func (se *SaleService) InitControlConsole(s service.Service) {
	svcFlag := flag.String("start", "", "start programm.")
	flag.Parse()
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
}

func Init() {
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := service.Config{
		Name:         "EphorSaleService",
		DisplayName:  "Ephor microservice for sales",
		Description:  "Microservice of Ephor company",
		Dependencies: []string{},
		Option:       options,
	}
	prg := &SaleService{}
	prg.ReadConfig()
	prg.InitLogger()
	prg.Manager = manager_service.New(prg.ConfigFile)
	prg.State = 1
	s, err := service.New(prg, &svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.InitControlConsole(s)
	err = s.Run()
	if err != nil {
		logger.Log.Error(err)
	}
}

/*
 Функция Run запускает микросервис и ставит состояние Idle
*/
func Run() {
	Init()
}
