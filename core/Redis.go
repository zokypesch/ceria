package core

import (
	"fmt"

	"github.com/mediocregopher/radix.v2/redis"
)

// RedisCore struct for redis
type RedisCore struct {
	Conn *redis.Client
}

// RedisCoreInter for interfacing RedisCore
type RedisCoreInter interface {
	CreateOrUpdateDocument(name string, id string, params ...interface{}) error
	DeleteDocument(name string, id string) error
	GetDocument(name string, id string, field string) (string, error)
	GetAllDocument(name string, id string) (map[string]string, error)
}

var rediscore *RedisCore

// NewServiceRedisCore for New RedisCore instance
func NewServiceRedisCore(host string, port string) (*RedisCore, error) {
	if rediscore == nil {
		con, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			return nil, err
		}

		rediscore = &RedisCore{
			Conn: con,
		}
	}
	return rediscore, nil
}

// CreateOrUpdateDocument func for insert or update redis
func (red *RedisCore) CreateOrUpdateDocument(name string, id string, params ...interface{}) error {

	if name == "" || id == "" || len(params) == 0 {
		return fmt.Errorf("params cannot be empty")
	}

	resp := red.Conn.Cmd("HMSET", fmt.Sprintf("%s:%s", name, id), params)
	// Check the Err field of the *Resp object for any errors.
	if resp.Err != nil {
		return resp.Err
	}
	return nil
}

// DeleteDocument func for delete redis
func (red *RedisCore) DeleteDocument(name string, id string) error {

	if name == "" || id == "" {
		return fmt.Errorf("params cannot be empty")
	}

	resp := red.Conn.Cmd("DEL", fmt.Sprintf("%s:%s", name, id))
	// Check the Err field of the *Resp object for any errors.
	if resp.Err != nil {
		return resp.Err
	}
	return nil
}

// GetDocument for get document return as string
func (red *RedisCore) GetDocument(name string, id string, field string) (string, error) {

	if name == "" || id == "" || field == "" {
		return "", fmt.Errorf("params cannot be empty")
	}

	resp, err := red.Conn.Cmd("HGET", fmt.Sprintf("%s:%s", name, id), field).Str()

	// Check the Err field of the *Resp object for any errors.
	return resp, err
}

// GetAllDocument for get all document return a new map
func (red *RedisCore) GetAllDocument(name string, id string) (map[string]string, error) {
	// var res map[string]string

	reply, err := red.Conn.Cmd("HGETALL", fmt.Sprintf("%s:%s", name, id)).Map()

	return reply, err
}
