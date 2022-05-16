package main
// main file for start microservice
import (
     "flag"
     "log"
     "os"
     "time"
     "github.com/kardianos/service"
     transactionDispetcher "transactionDispetcher"
     connectionPostgresql "connectionDB/connect"
	 ConnectionRabbitMQ "lib-rabbitmq"
     configEphor "configEphor"
     ephorpay "servicesEphor/ephorpay"
     ephorcommand "servicesEphor/ephorcommand"
     ephorfiscal "servicesEphor/ephorfiscal"
)
const (
    Start string = "start"
    Stop string = "stop"
)
var arrayTransaction map[int]interface{}
var logger service.Logger
    
// Program structures.
// Define Start and Stop methods.
type program struct {
    exit chan bool
    connectRabbit chan bool
}

func (p *program) ConnectionDb(conf *configEphor.Config) {
	ConnDb.NewConn(conf.Db.PgConnectionPool, conf.Db.Login, conf.Db.Password, conf.Db.Address, conf.Db.DatabaseName, conf.Db.Port)
	_,err := ConnDb.GetConn()
    if err != nil {
        logger.Infof("%s",err)
    }
}

func (p *program) ConnectRabbit(conf *configEphor.Config) {
    err := ConnectionRabbit.ConnectionToRabbit(conf.RabbitMq.Login, conf.RabbitMq.Password, conf.RabbitMq.Address, conf.RabbitMq.Port,conf.RabbitMq.MaxAttempts)
    if err != nil {
         logger.Infof("%v",err)
    }
    ConnectionRabbit.ConnectQueue()
}

func (p *program) Start(s service.Service) error {
    if service.Interactive() {
        logger.Info("Running in terminal.")
    } else {
        logger.Info("Running under service manager.")
    }
    p.exit = make(chan bool)
    p.connectRabbit = make(chan bool)

    // Start should not block. Do the actual work async.
    go p.run()
    return nil
}

func (p *program) checkStatusRabbit(){
     go func(){
        for {
            select {
                case <-p.connectRabbit:
                log.Println("Reconect Queues")
                go ephorpay.ReconnectQueue()
                go ephorcommand.ReconnectQueue()
                go ephorfiscal.ReconnectQueue()
            }
            time.Sleep(10 * time.Second)
        }
    }()
}

func (p *program) run() error {
    transactions := transactionDispetcher.New()
    logger.Infof("I'm running %v.", service.Platform())
    p.ConnectionDb(&cfg)
    p.ConnectRabbit(&cfg)
    p.checkStatusRabbit()
    go ConnectionRabbit.Reconnect(p.connectRabbit)
    go ephorpay.Start(&cfg,&ConnectionRabbit,ConnDb,p.exit,p.connectRabbit,*transactions)
    go ephorcommand.Start(&cfg,&ConnectionRabbit,ConnDb,p.exit,p.connectRabbit,*transactions)
    go ephorfiscal.Start(&cfg,&ConnectionRabbit,ConnDb,p.exit,p.connectRabbit,*transactions)
    return nil
}
func (p *program) Stop(s service.Service) error {
    // Any work in Stop should be quick, usually a few seconds at most.
    p.exit <- true
    ConnDb.CloseConnectionDb()
    ConnectionRabbit.CloseConnectRabbit()
    logger.Info("I'm Stopping!")
    close(p.exit)
    return nil
}

func (p *program) Status(s service.Service) error {
    logger.Info("I'm status!")
    return nil
}

var cfg configEphor.Config
var ConnectionRabbit ConnectionRabbitMQ.ChannelMQ
var ConnDb connectionPostgresql.DatabaseInstance

func main() {
    svcFlag := flag.String("start","", "start programm.") 
    log.SetFlags(log.Lshortfile)
    flag.Parse()
    cfg.Load()
    if cfg.LogFile != "" {
		file, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Panic(err)
		}
		log.SetOutput(file)
	}
	log.Println("Load Config...")
    options := make(service.KeyValue)
    options["Restart"] = "on-success"
    options["SuccessExitStatus"] = "1 2 8 SIGKILL"
    svcConfig := service.Config{
        Name:         "EphorMicroservice",
        DisplayName:  "Ephor microservice",
        Description:  "Microservice of Ephor company",
        Dependencies: []string{},
        Option: options,
    }

    prg := &program{}
    s, err := service.New(prg, &svcConfig)
   
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Start service")
    errs := make(chan error, 5)
    logger, err = s.Logger(errs)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Start Logger")
    if len(*svcFlag) != 0 {
        err := service.Control(s, *svcFlag)
        if err != nil {
            log.Printf("Valid actions: %q\n", service.ControlAction)
            log.Fatal(err)
        }
        return
    }
    err = s.Run()
    if err != nil {
        logger.Error(err)
    }

}




