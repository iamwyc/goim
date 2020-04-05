package main

// Start Commond eg: ./client 1 1000 localhost:3101
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	log "github.com/golang/glog"
)

const (
	opHeartbeat           = int32(2)
	opHeartbeatReply      = int32(3)
	opAuth                = int32(7)
	opAuthReply           = int32(8)
	OpGetOfflineMessage   = int32(18)
	OpBusinessMessagePush = int32(900)
	OpBusinessMessageAck  = int32(901)
)

const (
	rawHeaderLen = uint16(16)
	heart        = 10 * time.Second
)

const snPrefix = "00:11:22:33:44:"

// Proto proto.
type Proto struct {
	PackLen   int32  // package length
	HeaderLen int16  // header length
	Ver       int16  // protocol version
	Operation int32  // operation for request
	Seq       int32  // sequence number chosen by client
	Body      []byte // body
}

// AuthToken auth token.
type AuthToken struct {
	Mid      int64   `json:"mid"`
	Key      string  `json:"key"`
	RoomID   string  `json:"room_id"`
	Platform string  `json:"platform"`
	Accepts  []int32 `json:"accepts"`
}

var (
	countDown  int64
	aliveCount int64
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	begin := 12
	num := 1
	ip := "127.0.0.1:3101"
	go result()
	for i := begin; i < begin+num; i++ {
		go client(int64(i), ip)
	}
	// signal
	var exit chan bool
	<-exit
}

func result() {
	var (
		lastTimes int64
		interval  = int64(5)
	)
	for {
		nowCount := atomic.LoadInt64(&countDown)
		nowAlive := atomic.LoadInt64(&aliveCount)
		diff := nowCount - lastTimes
		lastTimes = nowCount
		fmt.Println(fmt.Sprintf("%s alive:%d down:%d down/s:%d", time.Now().Format("2006-01-02 15:04:05"), nowAlive, nowCount, diff/interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func client(mid int64, ip string) {
	for {
		startClient(mid, ip)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func startClient(key int64, ip string) {
	time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
	atomic.AddInt64(&aliveCount, 1)
	quit := make(chan bool, 1)
	defer func() {
		close(quit)
		atomic.AddInt64(&aliveCount, -1)
	}()
	// connnect to server
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Errorf("net.Dial(%s) error(%v)", ip, err)
		return
	}
	seq := int32(0)
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)
	authToken := &AuthToken{
		key,
		snPrefix + strconv.FormatInt(key, 10),
		"test://1",
		"ios",
		[]int32{1000, 1001, 1002},
	}
	proto := new(Proto)
	proto.Ver = 1
	proto.Operation = opAuth
	proto.Seq = seq
	proto.Body, _ = json.Marshal(authToken)
	if err = tcpWriteProto(wr, proto); err != nil {
		log.Errorf("tcpWriteProto() error(%v)", err)
		return
	}
	if err = tcpReadProto(rd, proto); err != nil {
		log.Errorf("tcpReadProto() error(%v)", err)
		return
	}
	log.Infof("key:%d auth ok, proto: %v", key, proto)
	seq++
	// writer
	go func() {
		hbProto := new(Proto)
		for {
			// heartbeat
			hbProto.Operation = opHeartbeat
			if seq%2 == 0 {
				hbProto.Operation = OpGetOfflineMessage
			}
			hbProto.Seq = seq
			hbProto.Body = nil
			if err = tcpWriteProto(wr, hbProto); err != nil {
				log.Errorf("key:%d tcpWriteProto() error(%v)", key, err)
				return
			}
			log.Infof("key:%d Write heartbeat", key)
			time.Sleep(heart)
			seq++
			select {
			case <-quit:
				return
			default:
			}
		}
	}()
	// reader
	for {
		if err = tcpReadProto(rd, proto); err != nil {
			log.Errorf("key:%d tcpReadProto() error(%v)", key, err)
			quit <- true
			return
		}
		if proto.Operation == opAuthReply {
			log.Infof("key:%d auth success", key)
		} else if proto.Operation == opHeartbeatReply {
			log.Infof("key:%d receive heartbeat", key)
			if err = conn.SetReadDeadline(time.Now().Add(heart + 60*time.Second)); err != nil {
				log.Errorf("conn.SetReadDeadline() error(%v)", err)
				quit <- true
				return
			}
			log.Infof("key:%d seq:%d op:%d msg: %v", key, proto.Seq, proto.Operation, proto.Body)
		} else if proto.Operation == OpBusinessMessagePush {
			proto.Body = nil
			proto.Operation = OpBusinessMessageAck
			tcpWriteProto(wr, proto)
			log.Infof("packlen:%d key:%d seq:%d op:%d msglen:%d msg: %s", proto.PackLen, key, proto.Seq, proto.Operation, len(proto.Body), string(proto.Body))

		} else {
			log.Infof("packlen:%d key:%d seq:%d op:%d msglen:%d msg: %s", proto.PackLen, key, proto.Seq, proto.Operation, len(proto.Body), string(proto.Body))
			atomic.AddInt64(&countDown, 1)
		}
	}
}

func tcpWriteProto(wr *bufio.Writer, proto *Proto) (err error) {
	// write
	if err = binary.Write(wr, binary.BigEndian, uint32(rawHeaderLen)+uint32(len(proto.Body))); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, rawHeaderLen); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Ver); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Operation); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Seq); err != nil {
		return
	}
	if proto.Body != nil {
		if err = binary.Write(wr, binary.BigEndian, proto.Body); err != nil {
			return
		}
	}
	err = wr.Flush()
	return
}

func tcpReadProto(rd *bufio.Reader, proto *Proto) (err error) {
	var (
		packLen   int32
		headerLen int16
	)
	// read
	if err = binary.Read(rd, binary.BigEndian, &packLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &headerLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Ver); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Operation); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Seq); err != nil {
		return
	}
	var (
		n, t    int
		bodyLen = int(packLen - int32(headerLen))
	)
	proto.PackLen = packLen
	proto.HeaderLen = headerLen
	if bodyLen > 0 {
		proto.Body = make([]byte, bodyLen)
		for {
			if t, err = rd.Read(proto.Body[n:]); err != nil {
				return
			}
			if n += t; n == bodyLen {
				break
			}
		}
	} else {
		proto.Body = nil
	}
	return
}
