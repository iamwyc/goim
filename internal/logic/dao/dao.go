package dao

import (
	"context"
	"time"

	"github.com/Terry-Mao/goim/internal/logic/conf"
	"github.com/golang/glog"
	"github.com/robfig/cron"

	"github.com/go-redis/redis/v8"
	kafka "gopkg.in/Shopify/sarama.v1"
	"gopkg.in/mgo.v2"
)

// Dao dao.
type Dao struct {
	c            *conf.Config
	kafkaPub     kafka.SyncProducer
	redis        *redis.ClusterClient
	redisExpire  int32
	MongoSession *mgo.Session
	idWorker     *IDWorker
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
		c:            c,
		kafkaPub:     newKafkaPub(c.Kafka),
		redis:        newRedis(c.Redis),
		redisExpire:  int32(time.Duration(c.Redis.Expire) / time.Second),
		MongoSession: mSession,
		idWorker:     NewIDWorker(c.MessagePush.WorkerID),
	}
	if c.MessagePush.EnableCron {
		glog.Infof("cron start...")
		cr := cron.New()
		cr.AddFunc("0 0 2 * * ?", func() {
			d.MessageStats()
		})
		cr.Start()
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

func newRedis(c *conf.Redis) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{":6304", ":6305", ":6300", ":6301", "6302", ":6303"},
		ReadOnly: false,
	})
}

// Close close the resource.
func (d *Dao) Close() error {
	return d.redis.Close()
}

// Ping dao ping.
func (d *Dao) Ping(c context.Context) error {
	return d.pingRedis(c)
}
