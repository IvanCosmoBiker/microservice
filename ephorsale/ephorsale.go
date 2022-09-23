package ephorsale

// main file for start microservice
import (
	config "ephorservices/config"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
)

// Program structures.
// Define Start and Stop methods.
type SaleService struct {
	State      int
	Logger     *log.Logger
	Config     service.Config
	ConfigFile *config.Config
}

var manager SaleServiceManager

func (se *SaleService) Start(s service.Service) error {
	if service.Interactive() {
		log.Print("Running in terminal.")
	} else {
		log.Print("Running under service manager.")
	}
	// Start should not block. Do the actual work async.
	go se.run(s)
	return nil
}

func (se *SaleService) run(s service.Service) error {
	log.Printf("%v", se.ConfigFile.Db.Login)
	manager.InitService(se)
	return nil
}

func (se *SaleService) Stop(s service.Service) error {
	manager.StopService()
	log.Print("I'm Stopping!")
	return nil
}

func (se *SaleService) Status(s service.Service) error {
	log.Print("I'm status!")
	return nil
}

func ReadConfig() (*config.Config, bool) {
	var config = config.Config{}
	config.Load()
	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(config.Db.Login)
		log.SetOutput(file)
		return &config, true
	}
	return &config, true
}

func Init() {
	var logger service.Logger
	svcFlag := flag.String("start", "", "start programm.")
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	config, _ := ReadConfig()
	log.Println("Load Config...")
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
	prg.ConfigFile = config
	s, err := service.New(prg, &svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.State = 1
	errs := make(chan error, 10)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Load Logger...")
	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	log.Println("Start SaleService")
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
