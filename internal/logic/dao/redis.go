package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/go-redis/redis/v8"
	log "github.com/golang/glog"

	"github.com/zhenjl/cityhash"
)

const (
	_prefixMidServer    = "mid_%d"     // mid -> key:server
	_prefixKeyServer    = "key_%s"     // key -> server
	_prefixServerOnline = "ol_%s"      // server -> online
	_prefixmessageCount = "message_%d" // server -> online
)

func messageCount(msgID int64) string {
	return fmt.Sprintf(_prefixmessageCount, msgID)
}

func keyMidServer(mid int64) string {
	return fmt.Sprintf(_prefixMidServer, mid)
}

func keyKeyServer(key string) string {
	return fmt.Sprintf(_prefixKeyServer, key)
}

func keyServerOnline(key string) string {
	return fmt.Sprintf(_prefixServerOnline, key)
}

// pingRedis check redis connection.
func (d *Dao) pingRedis(c context.Context) (err error) {

	return d.redis.Ping(c).Err()
}

// AddMapping add a mapping.
// Mapping:
//	mid -> key_server
//	key -> server
func (d *Dao) AddMapping(ctx context.Context, mid int64, key, server string) (err error) {
	_, err = d.redis.Pipelined(ctx, func(pipe redis.Pipeliner) (e error) {
		if mid > 0 {
			if _, e = pipe.Do(ctx, "HSET", keyMidServer(mid), key, server).Result(); err != nil {
				return
			}
			if _, e = pipe.Do(ctx, "EXPIRE", keyMidServer(mid), d.redisExpire).Result(); err != nil {
				return
			}
		}
		if _, e = pipe.Do(ctx, "SET", keyKeyServer(key), server).Result(); err != nil {
			return
		}
		if _, e = pipe.Do(ctx, "EXPIRE", keyKeyServer(key), d.redisExpire).Result(); err != nil {
			return
		}
		return nil
	})

	return
}

// ExpireMapping expire a mapping.
func (d *Dao) ExpireMapping(ctx context.Context, mid int64, key string) (has bool, err error) {
	cmd, err := d.redis.Pipelined(ctx, func(pipe redis.Pipeliner) (e error) {
		if _, e = pipe.Expire(ctx, keyMidServer(mid), time.Duration(d.redisExpire)*time.Second).Result(); err != nil {
			return
		}
		if _, e = pipe.Expire(ctx, keyKeyServer(key), time.Duration(d.redisExpire)*time.Second).Result(); err != nil {
			return
		}
		return nil
	})
	has = true
	for _, c := range cmd {
		d := c.(*redis.BoolCmd)
		has = has && d.Val()
	}
	return
}

// DelMapping del a mapping.
func (d *Dao) DelMapping(ctx context.Context, mid int64, key, server string) (has bool, err error) {
	cmd, err := d.redis.Pipelined(ctx, func(pipe redis.Pipeliner) (e error) {

		if _, e = pipe.HDel(ctx, keyMidServer(mid), key).Result(); err != nil {
			return
		}
		if _, e = pipe.Del(ctx, keyKeyServer(key)).Result(); err != nil {
			return
		}

		return nil
	})
	has = true
	for _, c := range cmd {
		d := c.(*redis.IntCmd)
		has = has && d.Val() > 0
	}
	return
}

// ServersByKeys get a server by key.
func (d *Dao) ServersByKeys(c context.Context, keys []string) (res []string, err error) {
	cmd, _ := d.redis.Pipelined(c, func(pipe redis.Pipeliner) (e error) {
		for _, key := range keys {
			pipe.Get(c, keyKeyServer(key))
		}
		return
	})

	log.Infof("%v", cmd)
	for _, c := range cmd {

		if c.Err() == redis.Nil {
			continue
		}

		d := c.(*redis.StringCmd)
		if d.Err() != nil {
			return nil, d.Err()
		}

		res = append(res, d.Val())
	}
	return
}

// KeysByMids get a key server by mid.
func (d *Dao) KeysByMids(c context.Context, mids []int64) (ress map[string]string, olMids []int64, err error) {
	ress = make(map[string]string)

	for idx, mid := range mids {
		r := d.redis.HGetAll(c, keyMidServer(mid))
		if err = r.Err(); err != nil {
			log.Errorf("conn.Do(HGETALL %d) error(%v)", mid, err)
			return
		}
		if len(r.Val()) > 0 {
			olMids = append(olMids, mids[idx])
		}

		for k, v := range r.Val() {
			ress[k] = v
		}
	}
	return
}

// AddServerOnline add a server online.
func (d *Dao) AddServerOnline(c context.Context, server string, online *model.Online) (err error) {
	roomsMap := map[uint32]map[string]int32{}
	for room, count := range online.RoomCount {
		rMap := roomsMap[cityhash.CityHash32([]byte(room), uint32(len(room)))%64]
		if rMap == nil {
			rMap = make(map[string]int32)
			roomsMap[cityhash.CityHash32([]byte(room), uint32(len(room)))%64] = rMap
		}
		rMap[room] = count
	}
	key := keyServerOnline(server)
	for hashKey, value := range roomsMap {
		err = d.addServerOnline(c, key, strconv.FormatInt(int64(hashKey), 10), &model.Online{RoomCount: value, Server: online.Server, Updated: online.Updated})
		if err != nil {
			return
		}
	}
	return
}

func (d *Dao) addServerOnline(c context.Context, key string, hashKey string, online *model.Online) (err error) {
	b, _ := json.Marshal(online)
	_, err = d.redis.Pipelined(c, func(pipe redis.Pipeliner) (e error) {
		if _, e = pipe.Do(c, "HSET", key, hashKey, b).Result(); err != nil {
			return
		}
		if _, e = pipe.Do(c, "EXPIRE", key, d.redisExpire).Result(); err != nil {
			return
		}
		return nil
	})
	return
}

// ServerOnline get a server online.
func (d *Dao) ServerOnline(c context.Context, server string) (online *model.Online, err error) {
	online = &model.Online{RoomCount: map[string]int32{}}
	key := keyServerOnline(server)
	for i := 0; i < 64; i++ {
		ol, err := d.serverOnline(c, key, strconv.FormatInt(int64(i), 10))
		if err == nil && ol != nil {
			online.Server = ol.Server
			if ol.Updated > online.Updated {
				online.Updated = ol.Updated
			}
			for room, count := range ol.RoomCount {
				online.RoomCount[room] = count
			}
		}
	}
	return
}

func (d *Dao) serverOnline(c context.Context, key string, hashKey string) (online *model.Online, err error) {
	b, err := d.redis.HGet(c, key, hashKey).Bytes()
	if err != nil {
		if err != redis.Nil {
			log.Errorf("conn.Do(HGET %s %s) error(%v)", key, hashKey, err)
		}
		return
	}
	online = new(model.Online)
	if err = json.Unmarshal(b, online); err != nil {
		log.Errorf("serverOnline json.Unmarshal(%s) error(%v)", b, err)
		return
	}

	_, err = d.redis.Pipelined(c, func(pipe redis.Pipeliner) (e error) {
		if _, e = pipe.Do(c, "HSET", key, hashKey, b).Result(); err != nil {
			return
		}
		if _, e = pipe.Do(c, "EXPIRE", key, d.redisExpire).Result(); err != nil {
			return
		}
		return nil
	})
	return
}

// DelServerOnline del a server online.
func (d *Dao) DelServerOnline(c context.Context, server string) (err error) {
	key := keyServerOnline(server)
	if _, err = d.redis.Del(c, key).Result(); err != nil {
		log.Errorf("conn.Do(DEL %s) error(%v)", key, err)
	}
	return
}

// MessageSeqAdd mesage incr
func (d *Dao) MessageSeqAdd(c context.Context, msgID int64) (err error) {
	key := messageCount(msgID)
	count, err := d.redis.Incr(c, key).Uint64()
	if count <= uint64(2) && err == nil {
		if err = d.redis.Expire(c, key, time.Duration(time.Second*345600)).Err(); err != nil {
			log.Errorf("conn.Send(EXPIRE %s) error(%v)", key, err)
			return
		}
	}
	return
}

// MessageCountStats mesage count
func (d *Dao) MessageCountStats(c context.Context, msgID int64) (count int64, err error) {
	count, err = d.redis.Get(c, messageCount(msgID)).Int64()
	if err == nil {
		return
	}
	m, err := d.GetMessageByID(msgID)
	if err == nil && m != nil {
		count = int64(m.PushCount)
	} else {
		count = int64(0)
		err = nil
	}
	return
}
