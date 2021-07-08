package conf

import (
	"sync/atomic"

	"github.com/BurntSushi/toml"
	"github.com/gamedev-embers/imnotifier/notifiers"
)

var cfg atomic.Value

func Init(fpath string) error {
	c := config{}
	_, err := toml.DecodeFile(fpath, &c)
	cfg.Store(&c)
	return err
}

type config struct {
	Server struct {
		Listen string `toml:"listen"`
	} `toml:"server"`

	Notifiers map[string]*notifiers.Notifier `toml:"notifier"`
}

func Config() *config {
	return cfg.Load().(*config)
}
