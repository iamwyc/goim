package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Terry-Mao/goim/internal/logic"
	"github.com/Terry-Mao/goim/internal/logic/conf"
	"github.com/Terry-Mao/goim/internal/logic/grpc"

	"github.com/Terry-Mao/goim/internal/logic/http"
	"github.com/Terry-Mao/goim/internal/logic/model"
	"github.com/Terry-Mao/goim/pkg/ip"
	"github.com/bilibili/discovery/naming"
	resolver "github.com/bilibili/discovery/naming/grpc"
	log "github.com/golang/glog"
	ggrpc "google.golang.org/grpc"
)

const (
	ver = "2.0.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	var (
		httpSrv *http.Server
		rpcSrv  *ggrpc.Server
		appID   = conf.Conf.Env.AppID
	)
	log.Infof("%s [version: %s env: %+v push: %+v] start", appID, ver, conf.Conf.Env, conf.Conf.MessagePush)
	// grpc register naming
	dis := naming.New(conf.Conf.Discovery)
	resolver.Register(dis)
	// logic
	srv := logic.New(conf.Conf)

	if conf.Conf.HTTPServer.Enable {
		httpSrv = http.New(conf.Conf.HTTPServer, srv)
	}
	if conf.Conf.RPCServer.Enable {
		rpcSrv = grpc.New(conf.Conf.RPCServer, srv)
	}
	cancel := register(dis, srv, conf.Conf.Env.AppID)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("%s get a signal %s", appID, s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			srv.Close()
			if httpSrv != nil {
				httpSrv.Close()
			}
			if rpcSrv != nil {
				rpcSrv.GracefulStop()
			}
			log.Infof("%s [version: %s] exit", appID, ver)
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func register(dis *naming.Discovery, srv *logic.Logic, appid string) context.CancelFunc {
	env := conf.Conf.Env
	addr := ip.InternalIP()
	_, port, _ := net.SplitHostPort(conf.Conf.RPCServer.Addr)
	ins := &naming.Instance{
		Region:   env.Region,
		Zone:     env.Zone,
		Env:      env.DeployEnv,
		Hostname: env.Host,
		AppID:    appid,
		Addrs: []string{
			"grpc://" + addr + ":" + port,
		},
		Metadata: map[string]string{
			model.MetaWeight: strconv.FormatInt(env.Weight, 10),
		},
	}
	cancel, err := dis.Register(ins)
	if err != nil {
		panic(err)
	}
	return cancel
}
