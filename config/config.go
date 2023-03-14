package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

/* Configuration structure, load from config.json, global */
type Config struct {
	RabbitMq struct {
		Login, Password, Address, Port string
		MaxAttempts                    int
		ExecuteTimeSeconds             time.Duration
		PoolChannel                    int8
		BackOffPolicySendMassage       []time.Duration
		BackOffPolicyConnection        []time.Duration
	}
	Transport struct {
		Mqtt struct {
			Login, Password, Address, Port, ClientID string
			BackOffPolicySendMassage                 []time.Duration
			BackOffPolicyConnection                  []time.Duration
			ExecuteTimeSeconds                       int
			Subscribers                              []string
			KeepAlive                                uint32
			ProtocolVersion                          uint
			Disconnect                               uint
			Consumers                                []string
			PoolPublisher                            uint8
		}
		Http struct {
			Address, Port string
		}
	}
	Db struct {
		ReconnectSecond                                          int
		Login, Password, Address, DatabaseName                   string
		PreferSimpleProtocol                                     bool
		Port, PgConnectionPool, PgConnectionMin, PgConnectionMax uint16
	}
	Services struct {
		EnableControl bool
		Http          struct {
			Address, Port string
		}
		EphorPayment struct {
			Config struct {
				ExecuteMinutes time.Duration // this parametr for time run work with bank
				IntervalTime   time.Duration
			}
			Transport struct {
				Grpc struct {
					Address, Port string
				}
				Http struct {
					Address, Port string
				}
				Mqtt struct {
					NameQueue string
				}
			}
		}
		EphorPay struct {
			NameQueue string
			Bank      struct {
				ExecuteMinutes time.Duration // this parametr for time run work with bank
				PollingTime    time.Duration
			}
			Controller struct {
				Transport struct {
					Http struct {
						Address, Port string
					}
				}
			}
		}
		EphorCommand struct {
			NameQueue      string
			ExecuteMinutes time.Duration
			Listener       struct {
				ExecuteMinutes time.Duration
			}
			Controller struct {
				Transport struct {
					Http struct {
						Address, Port string
					}
				}
			}
		}
		EphorFiscal struct {
			NameQueue                   string
			ResponseUrl                 string
			PathCert                    string
			ExecuteMinutes, SleepSecond int
			Listener                    struct {
				ExecuteMinutes time.Duration
			}
			Controller struct {
				Transport struct {
					Http struct {
						Address, Port string
					}
				}
			}
		}
	}
	ErrorCount     int
	Log            string
	PrefixLog      string
	LogAvalable    []int
	LogEnable      bool
	ExecuteMinutes time.Duration // this parametr work execute time for one transaction
	Debug          bool
}

func (c *Config) Load() {
	file, _ := os.Open("config.json")
	byteValue, _ := ioutil.ReadAll(file)
	defer file.Close()
	json.Unmarshal(byteValue, &c)
	fmt.Printf("%+v", c)
}
