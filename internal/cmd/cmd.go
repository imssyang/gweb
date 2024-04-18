package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imssyang/gweb/internal/api"
	"github.com/imssyang/gweb/internal/api/formatify"
	"github.com/imssyang/gweb/internal/conf"
	"github.com/imssyang/gweb/internal/log"
	"github.com/imssyang/gweb/public"
	"github.com/imssyang/gweb/templates"
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Category:    "Miscellaneous:",
		Name:        "debug",
		Aliases:     []string{"d"},
		Usage:       "Enable debug mode",
		Value:       conf.App.Debug,
		Destination: &conf.App.Debug,
	},
	&cli.BoolFlag{
		Category:    "Miscellaneous:",
		Name:        "silent",
		Aliases:     []string{"s"},
		Usage:       "Suppress log messages from output",
		Value:       conf.App.Silent,
		Destination: &conf.App.Silent,
	},
	&cli.StringFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Usage:       "Load configuration from file",
		Destination: &conf.App.Config,
	},
	&cli.StringFlag{
		Name:        "bind",
		Aliases:     []string{"b"},
		Usage:       "Address to use",
		Value:       conf.App.Service.Host,
		Destination: &conf.App.Service.Host,
	},
	&cli.IntFlag{
		Name:        "port",
		Aliases:     []string{"p"},
		Usage:       "Port to use",
		Value:       conf.App.Service.Port,
		Destination: &conf.App.Service.Port,
		Action: func(ctx *cli.Context, v int) error {
			if v >= 65536 {
				return fmt.Errorf("flag: port value %v out of range[0-65535]", v)
			}
			return nil
		},
	},
}

func Action(ctx *cli.Context) error {
	if len(conf.App.Config) > 0 {
		conf.App.Load(conf.App.Config)
	} else {
		conf.App.Encap()
	}

	log.Zap.Debugf("Config: %+v", conf.App)

	gin.SetMode(conf.GinMode())
	engine := gin.New()
	engine.Use(log.GinLogger())
	engine.Use(gin.Recovery())
	templates.Init(engine)
	public.Init(engine)

	api.Register(engine)
	formatify.Register(engine)

	server := &http.Server{
		Addr:           conf.App.Service.Address,
		Handler:        engine,
		ReadTimeout:    conf.App.Service.Timeout.Read,
		WriteTimeout:   conf.App.Service.Timeout.Write,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Zap.Fatalf("Server run error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Zap.Infof("Server shutdown by signal: %v", sig)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	if err := server.Shutdown(ctxTimeout); err != nil {
		log.Zap.Fatalln("Server exit error:", err)
	} else {
		log.Zap.Infoln("Server exit")
	}

	return nil
}
