package robot

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/db"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/util"
	"gopkg.in/mgo.v2"
	"strconv"
)

//RecoveryRobot robô que faz a resiliência de boletos
func RecoveryRobot(ex string) {

	if ex == "true" {
		go func() {
			e, _ := strconv.ParseUint(config.Get().RecoveryRobotExecutionInMinutes, 10, 64)
			gocron.Every(e).Minutes().Do(executionTask)
			gocron.RunAll()
			<-gocron.Start()
		}()
	}
}

func executionTask() {

	lg := log.CreateLog()
	lg.Operation = "RecoveryRobot"

	redis := db.CreateRedis()
	keys, err := redis.GetAllJSON(lg)
	lg.InitRobot(len(keys))

	if util.CheckErrorRobot(err) == false && len(keys) != 0 {
		mongo, errMongo := db.CreateMongo(lg)
		if util.CheckErrorRobot(errMongo) == false {
			for _, key := range keys {
				bol, errRedis := redis.GetBoletoJSONByKey(key, lg)
				lg.RequestKey = bol.Boleto.RequestKey

				if util.CheckErrorRobot(errRedis) == false {
					idBoleto, _ := bol.ID.MarshalText()

					b, _ := mongo.GetBoletoByID(string(idBoleto), bol.PublicKey)
					
					if b.ID == "" {
						err = mongo.SaveBoleto(bol)

						/*If the error when saving to the mongo is because the key is duplicated,
						we will ignore it, as this is a way to control the flow, if there are two instances running.*/
						if lastErr, ok := err.(*mgo.LastError); (ok && lastErr.Code != 11000) || (!ok && err != nil) {
							lg.Warn(err.Error(), fmt.Sprintf("Error saving to mongo - %s", key))
						}
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
	}

	lg.EndRobot()
}