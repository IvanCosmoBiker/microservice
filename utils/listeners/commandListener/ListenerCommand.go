package commandListener

import(
	"encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
	"time"
    "log"
    configEphor "configEphor"
    connectionPostgresql "connectionDB/connect"
    ConnectionRabbitMQ "lib-rabbitmq"
    listener "listeners/v1"
    requestCommand "data/command"
    commandModel "internal/command"
)

func ProcessingRequest(json_data []byte) {
    response  := make(map[string]interface{})
    json.Unmarshal(json_data, &CommandRequest)
    connectDb.AddLog(fmt.Sprintf("%+v",CommandRequest),"CommandSend" ,fmt.Sprintf("%+v",CommandRequest),"CommandSend")
    entry := make(map[string]interface{})
    entry["id"] = CommandRequest.Id
    entry["sended"] = commandModel.SendSuccess
    log.Println(CommandRequest.Id)
    command := CommandStore.GetOneById(CommandRequest.Id)
    CommandStore.SetByParams(entry)
    cmd,_ := command.Command.Value()
    sum,_ := command.Command_param1.Value()
    response["a"] = cmd.(int64)
    response["m"]  = 2
    response["sum"] = sum.(int64)
    data, _ := json.Marshal(response)
    Rabbit.PublishMessage(data,fmt.Sprintf("ephor.1.dev.%v",CommandRequest.Imei))
}

func handler(w http.ResponseWriter, req *http.Request) {
        switch req.Method {
        case "POST":
            json_data, _ := ioutil.ReadAll(req.Body)
            defer req.Body.Close()
            go ProcessingRequest(json_data)
            go func() {
                select {
                    case <-time.After(time.Duration(conf.Services.EphorCommand.Listener.ExecuteMinutes) * time.Minute):
                        return 
                }
            }()
            return 
        case "GET":
            fmt.Fprintf(w, "%s: Running\n", "Command")
        default:
            fmt.Fprintf(w, "Sorry, only POST and GET method is supported.")
        }
}
var connectDb connectionPostgresql.DatabaseInstance
var conf *configEphor.Config
var Rabbit *ConnectionRabbitMQ.ChannelMQ
var CommandRequest  requestCommand.CommandServerRequest
var CommandStore  commandModel.StoreCommand
var ChannelParentClose chan bool
func StarListen(cfg *configEphor.Config,connectPg connectionPostgresql.DatabaseInstance,rabbitmq *ConnectionRabbitMQ.ChannelMQ,channelParent chan bool ){
    connectDb = connectPg
    conf = cfg
    Rabbit = rabbitmq
    CommandStore.Connection = connectPg
    point := fmt.Sprintf("%s:%s",cfg.Services.Address,cfg.Services.Port)
    listener.StartListener("/command",point,handler)
    log.Println("Start Module Command..")
}
