package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type App struct {
	ctx    context.Context
	cancel func()
	opts   options
}

type Option func(o *options)

type options struct {
	ctx  context.Context
	sigs []os.Signal

	servers []http.Server
}

func Server(srv ...http.Server) Option {
	return func(o *options) { o.servers = srv }
}

// 创建生命周期管理器
func New(opts ...Option) *App {
	options := options{
		ctx:  context.Background(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
	}
}

// 启动
func (a *App) Run() error {
	eg, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done()
			return srv.Shutdown(ctx)
		})
		eg.Go(func() error {
			return srv.ListenAndServe()
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				a.Stop()
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// context取消
func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func NewServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Handler: handler,
		Addr:    port,
	}
}

func NewWelcomeServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Welcome, Gopher!")
	})
	return mux
}

func NewHelloServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello, Gopher!")
	})
	return mux
}

func main() {
	welcomeServer := NewServer(":8080", NewWelcomeServer())
	helloServer := NewServer(":8081", NewHelloServer())

	app := New(
		Server(*welcomeServer, *helloServer),
	)
	time.AfterFunc(5*time.Second, func() {
		app.Stop()
	})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
