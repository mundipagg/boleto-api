package robot

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/util"
	"strconv"
)

//RecoveryRobot robô que faz a resiliência de boletos
func RecoveryRobot(ex string) {

	if ex == "true" {
		go func() {
			e, _ := strconv.ParseUint(config.Get().RecoveryRobotExecutionInMinutes, 10, 64)
			gocron.Every(e).Minutes().Do(executionTask)
			<-gocron.Start()
		}()
	}
}

func executionTask() {

	lg := log.CreateLog()
	lg.Operation = "RecoveryRobot"

	lg.InitRobot()

	redis := db.CreateRedis()
	keys, _ := redis.GetAllJSON(lg)

	mongo, errMongo := db.CreateMongo(lg)
	if util.CheckErrorRobot(errMongo) == false {
		for _, key := range keys {
			bol, errRedis := redis.GetBoletoJSONByKey(key, lg)

			if util.CheckErrorRobot(errRedis) == false {
				lg.RequestKey = bol.Boleto.RequestKey

				err := mongo.SaveBoleto(bol)
				if err != nil{
					lg.Warn(err.Error(), fmt.Sprintf("Error saving to mongo - %s", err.Error()))
				}

				if util.CheckErrorRobot(err) == false {
					errRedis = redis.DeleteBoletoJSONByKey(key, lg)

					if util.CheckErrorRobot(errRedis) == false{
						lg.ResumeRobot(key)
					}
				}
			}
		}
	}

	lg.EndRobot()
}
