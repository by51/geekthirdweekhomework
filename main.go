package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//初始化根节点context
	var ctx = context.Background()

	//创建一个用于控制子gorouting的context
	var canCtx,cancel = context.WithCancel(ctx)

	//子协程里面会用到
	var g, _ = errgroup.WithContext(canCtx)

	var addrs = []string{"httpaddres1","httpaddres2","httpaddres3"}

	for _, addr := range addrs {
		srv := &myServer{
			server: &http.Server{Addr: addr},
		}
		g.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.stopHttpSrver(canCtx)
		})
		g.Go(func() error {
			return srv.startHttpSrver()
		})
	}

	//接收linux signal信号
	// 创建一个os.Signal channel
	sigs := make(chan os.Signal, 1)

	//注册要接收的信号，syscall.SIGINT:接收ctrl+c ,syscall.SIGTERM:程序退出
	//信号没有信号参数表示接收所有的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-sigs:
				cancel()
			}
		}
	})
}

type myServer struct {
	server *http.Server
}

func(myServer *myServer) startHttpSrver () error {
	return myServer.server.ListenAndServe()
}

func(myServer *myServer) stopHttpSrver (ctx context.Context) error {
	return myServer.server.Shutdown(ctx)
}