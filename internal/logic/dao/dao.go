package dao

import (
	"context"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/Terry-Mao/goim/internal/logic/conf"
	"github.com/gomodule/redigo/redis"
	kafka "gopkg.in/Shopify/sarama.v1"
)

// Dao dao.
type Dao struct {
	c           *conf.Config
	kafkaPub    kafka.SyncProducer
	redis       *redis.Pool
	redisExpire int32
	MongoSession    *mgo.Session
}

// New new a dao and return.
func New(c *conf.Config) *Dao {

	mSession, err := mgo.Dial(c.Mongodb.MongoUrl)
	if err != nil {
		panic(err)
	}
	mSession.SetMode(mgo.Monotonic, true)
	mSession.SetPoolLimit(c.Mongodb.PoolLimit)
	mSession.SetSocketTimeout(time.Second * 10)
	d := &Dao{
		c:           c,
		kafkaPub:    newKafkaPub(c.Kafka),
		redis:       newRedis(c.Redis),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
		MongoSession:    mSession,
	}
	return d
}

func newKafkaPub(c *conf.Kafka) kafka.SyncProducer {
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(c.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

func newRedis(c *conf.Redis) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     c.Idle,
		MaxActive:   c.Active,
		IdleTimeout: time.Duration(c.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(c.Network, c.Addr,
				redis.DialConnectTimeout(time.Duration(c.DialTimeout)),
				redis.DialReadTimeout(time.Duration(c.ReadTimeout)),
				redis.DialWriteTimeout(time.Duration(c.WriteTimeout)),
				redis.DialPassword(c.Auth),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

// Close close the resource.
func (d *Dao) Close() error {
	return d.redis.Close()
}

// Ping dao ping.
func (d *Dao) Ping(c context.Context) error {
	return d.pingRedis(c)
}
