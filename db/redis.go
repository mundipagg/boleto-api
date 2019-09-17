package db

import (
	"encoding/json"
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
		lg.Warn(err, fmt.Sprintf("OpenConnection [SetBoletoHTML] - Could not connection to Redis Database "))
	} else {

		key := fmt.Sprintf("%s:%s:%s", "HTML", mID, pk)
		ret, err := r.conn.Do("SETEX", key, config.Get().RedisExpirationTime, b)

		res := fmt.Sprintf("%s", ret)

		if res != "OK" {
			lg.Warn(err, fmt.Sprintf("SetBoletoHTML [SetBoletoHTML] - Could not record HTML in Redis Database: %s", err.Error()))
		}

		r.closeConnection()
	}
}

//GetBoletoHTMLByID busca um boleto pelo ID que vem na URL
func (r *Redis) GetBoletoHTMLByID(id string, pk string, lg *log.Log) string {

	err := r.openConnection()

	if err != nil {
		lg.Warn(err, fmt.Sprintf("OpenConnection [GetBoletoHTMLByID] - Could not connection to Redis Database"))
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
		lg.Warn(err, fmt.Sprintf("OpenConnection [SetBoletoJSON] - Could not connection to Redis Database "))
		return err
	}

	key := fmt.Sprintf("%s:%s:%s", "JSON", mID, pk)
	ret, err := r.conn.Do("SET", key, b)
	r.closeConnection()

	res := fmt.Sprintf("%s", ret)

	if res != "OK" {
		lg.Warn(err, fmt.Sprintf("SetBoletoHTML [SetBoletoJSON] - Could not record HTML in Redis Database: %s", err.Error()))
		return err
	}

	return nil
}

// GetBoletoJSONByKey Recupera um boleto do tipo JSON do Redis
func (r *Redis) GetBoletoJSONByKey(key string, lg *log.Log) (models.BoletoView, error) {
	err := r.openConnection()

	if err != nil {
		lg.Warn(err, fmt.Sprintf("OpenConnection [GetBoletoJSONByKey] - Could not connection to Redis Database "))
		return models.BoletoView{}, err
	}

	ret, err := r.conn.Do("GET", key)
	r.closeConnection()

	if ret != nil {
		result := models.BoletoView{}
		r := fmt.Sprintf("%s", ret)
		_ = json.Unmarshal([]byte(r), &result)
		return result, nil
	}

	return models.BoletoView{}, err
}

// DeleteBoletoJSONByKey Deleta um boleto do tipo JSON do Redis
func (r *Redis) DeleteBoletoJSONByKey(key string, lg *log.Log) {
	err := r.openConnection()

	if err != nil {
		lg.Warn(err, fmt.Sprintf("OpenConnection [DeleteBoletoJSONByKey] - Could not connection to Redis Database "))
	} else {

		_, err = r.conn.Do("DEL", key)
		r.closeConnection()

	}

}

// GetAllJSON Recupera todas as keys JSON do Redis
func (r *Redis) GetAllJSON() ([]string, error) {

	err := r.openConnection()
	if err != nil {
		return nil, err
	}

	var keys []string

	arr, err := redis.Values(r.conn.Do("SCAN", 0, "MATCH", "JSON:*", "COUNT", 500))
	if err != nil {
		return nil, err
	}

	keys, _ = redis.Strings(arr[1], nil)

	r.closeConnection()

	return keys, nil

}
