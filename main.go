package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/gamedev-embers/imnotifier/conf"
	"github.com/gamedev-embers/imnotifier/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	// cli args
	cfgpath = flag.String("config", "./conf/conf.toml", "config file")
	listen  = flag.String("listen", "", "listen address")

	// web server
	app = echo.New()
	log = app.Logger
)

func setup() {
	flag.Parse()
	if err := conf.Init(*cfgpath); err != nil {
		log.Fatalf("loading config failed:%v", err)
		return
	}
	if *listen == "" {
		*listen = conf.Config().Server.Listen
		if *listen == "" {
			log.Fatalf("missing --listen")
		}
	}

	for _, notifier := range conf.Config().Notifiers {
		log.Infof("notifier: %s", notifier)
	}
}

func main() {
	setup()

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} [${method}] ${uri} ${status} in/out(${bytes_in}/${bytes_out}) ${latency_human}\n",
	}))
	app.Use(middleware.Recover())

	app.GET("/status", status)
	app.POST("/notify", notify)

	log.Infof("listening: %s", *listen)
	if err := app.Start(*listen); err != nil {
		log.Fatalf("starting app failed: %v", err)
		return
	}
}

func status(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func notify(c echo.Context) error {
	args, err := c.FormParams()
	if err != nil {
		return err
	}
	channels := args.Get("channels")
	text := args.Get("text")
	if channels == "" || text == "" {
		return c.String(600, "invalid params")
	}

	channelsArr := strings.Split(channels, ",")
	// TODO: async task queue
	if len(channelsArr) > 1 {
		return c.String(600, fmt.Sprintf("too many channels: %v", channelsArr))
	}
	notifiers := conf.Config().Notifiers
	msg := models.Text{Content: text}
	invalidChannels := []string{}
	for _, channel := range channelsArr {
		if n, ok := notifiers[channel]; ok {
			log.Infof("send msg to [%s]", channel)
			n.Send(&msg)
		} else {
			invalidChannels = append(invalidChannels, channel)
		}
	}

	if len(invalidChannels) > 0 {
		return c.String(601, fmt.Sprintf("invalid channels: %v", invalidChannels))
	}
	return c.NoContent(http.StatusOK)
}
