package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/mundipagg/boleto-api/config"
	"github.com/mundipagg/boleto-api/log"
	"github.com/mundipagg/boleto-api/models"
)

//Redis Classe de Conexão com o Banco REDIS
type Redis struct {
	conn redis.Conn
}

//CreateRedis Cria instancia do Struct Redis
func CreateRedis() *Redis {
	return new(Redis)
}

func (r *Redis) openConnection() error {
	dbID, _ := strconv.Atoi(config.Get().RedisDatabase)
	o := redis.DialDatabase(dbID)
	ps := redis.DialPassword(config.Get().RedisPassword)
	tOut := redis.DialConnectTimeout(15 * time.Second)
	userTLS := redis.DialUseTLS(config.Get().RedisSSL)

	c, err := redis.Dial("tcp", config.Get().RedisURL, o, ps, tOut, userTLS)
	if err != nil {
		return err
	}

	r.conn = c
	return nil
}

func (r *Redis) closeConnection() {
	r.conn.Close()
}

//SetBoletoHTML Grava um boleto em formato Html no Redis
func (r *Redis) SetBoletoHTML(b, mID, pk string, lg *log.Log) {
	err := r.openConnection()
	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("OpenConnection [SetBoletoHTML] - Could not connection to Redis Database "))
	} else {

		key := fmt.Sprintf("%s:%s:%s", "HTML", mID, pk)
		ret, err := r.conn.Do("SETEX", key, config.Get().RedisExpirationTime, b)

		res := fmt.Sprintf("%s", ret)

		if res != "OK" {
			lg.Warn(res, fmt.Sprintf("SetBoletoHTML [SetBoletoHTML] - Could not record HTML in Redis Database: %s", key))
		} else if err != nil{
			lg.Warn(err.Error(), fmt.Sprintf("Error Redis [SetBoletoHTML] - Could not record HTML in Redis Database: %s", key))
		}

		r.closeConnection()
	}
}

//GetBoletoHTMLByID busca um boleto pelo ID que vem na URL
func (r *Redis) GetBoletoHTMLByID(id string, pk string, lg *log.Log) string {

	err := r.openConnection()

	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("OpenConnection [GetBoletoHTMLByID] - Could not connection to Redis Database"))
		return ""
	}

	key := fmt.Sprintf("%s:%s:%s", "HTML", id, pk)
	ret, _ := r.conn.Do("GET", key)
	r.closeConnection()

	if ret == nil {
		return ""
	}

	return fmt.Sprintf("%s", ret)
}

//SetBoletoJSON Grava um boleto em formato JSON no Redis
func (r *Redis) SetBoletoJSON(b, mID, pk string, lg *log.Log) error {
	err := r.openConnection()

	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("OpenConnection [SetBoletoJSON] - Could not connection to Redis Database "))
		return err
	}

	key := fmt.Sprintf("%s:%s:%s", "JSON", mID, pk)
	ret, err := r.conn.Do("SET", key, b)
	res := fmt.Sprintf("%s", ret)

	r.closeConnection()

	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("SetBoletoJSON [SetBoletoJSON] - Could not record JSON in Redis Database: %s", key))
		return err
	}

	if res != "OK" {
		lg.Warn("could not record JSON", fmt.Sprintf("SetBoletoJSON [SetBoletoJSON] - Error could not record JSON in Redis Database: %s", key))
		return errors.New("could not record JSON")
	}

	return nil
}

// GetBoletoJSONByKey Recupera um boleto do tipo JSON do Redis
func (r *Redis) GetBoletoJSONByKey(key string, lg *log.Log) (models.BoletoView, error) {
	err := r.openConnection()

	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("OpenConnection [GetBoletoJSONByKey] - Could not connection to Redis Database "))
		return models.BoletoView{}, err
	}

	ret, err := r.conn.Do("GET", key)
	r.closeConnection()

	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("GetData [GetBoletoJSONByKey] - Error could not to get data - " + key))
		return models.BoletoView{}, err
	}

	if ret != nil {
		result := models.BoletoView{}
		r := fmt.Sprintf("%s", ret)
		err = json.Unmarshal([]byte(r), &result)

		if err != nil {
			lg.Warn(err.Error(), fmt.Sprintf("Deserialize [GetBoletoJSONByKey] - Could not deserialize json - " + key))
		}

		return result, err
	} else {
		lg.Warn("not found data", fmt.Sprintf("GetData [GetBoletoJSONByKey] - not found data - " + key))
		return models.BoletoView{}, errors.New("not found data")
	}
}

// DeleteBoletoJSONByKey Deleta um boleto do tipo JSON do Redis
func (r *Redis) DeleteBoletoJSONByKey(key string, lg *log.Log) error {
	err := r.openConnection()

	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("OpenConnection [DeleteBoletoJSONByKey] - Could not connection to Redis Database "))
		return err
	} else {
		result, err := r.conn.Do("DEL", key)
		r.closeConnection()

		if err != nil {
			lg.Warn(err.Error(), fmt.Sprintf("Delete data [DeleteBoletoJSONByKey] - Error on delete key: " + key))
			return err
		}

		if result == int64(0) {
			lg.Warn(result, fmt.Sprintf("Delete data [DeleteBoletoJSONByKey] - Data not deleted: " + key))
			return errors.New("data not deleted")
		}
	}

	return nil
}

// GetAllJSON Recupera todas as keys JSON do Redis
func (r *Redis) GetAllJSON(lg *log.Log) ([]string, error) {

	err := r.openConnection()
	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("OpenConnection [GetAllJson] - Could not connection to Redis Database "))
		return nil, err
	}

	var keys []string

	arr, err := redis.Values(r.conn.Do("SCAN", 0, "MATCH", "JSON:*", "COUNT", 500))
	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("ReadDb [GetAllJson] - Could not get all json from database"))
		return nil, err
	}

	keys, err = redis.Strings(arr[1], nil)

	r.closeConnection()

	if err != nil {
		lg.Warn(err.Error(), fmt.Sprintf("Stringify Keys [GetAllJson] - Error when setting keys in the array"))
		return nil, err
	}

	return keys, nil

}
