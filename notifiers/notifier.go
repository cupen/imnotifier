package notifiers

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gamedev-embers/imnotifier/notifiers/feishu"
)

// Notifier ...
type Notifier struct {
	URL     string `toml:"url"`
	Timeout int    `toml:"timeout"`

	router _Router `toml:"-"`
	url    string  `toml:"-"`
}

// Router ...
type _Router func(url string, msg interface{}, timeout ...time.Duration) error

var routers = map[string]_Router{}

func init() {
	routers["feishu"] = feishu.Send
}

// Send ...
func (n *Notifier) Send(msg interface{}) {
	err := n.send(msg)
	panicIf(err)
}

// SendAsync ...
func (n *Notifier) SendAsync(msg interface{}, ch chan error) {
	var rs error
	defer func() {
		if err := recover(); err != nil {
			if _err, ok := err.(error); ok {
				rs = _err
			} else {
				rs = fmt.Errorf("%v", err)
			}
			ch <- rs
		}
	}()
	go func() {
		rs = n.send(msg)
	}()
}

func (n *Notifier) send(msg interface{}) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if n == nil {
		return fmt.Errorf("nil notifier")
	}
	router, err := n.getRouter()
	if err != nil {
		return err
	}
	return router(n.url, msg, n.getTimeout())
}

func (n *Notifier) getTimeout() time.Duration {
	return time.Duration(n.Timeout) * time.Second
}

func (n *Notifier) getRouter() (_Router, error) {
	if n.router == nil {
		u, err := url.Parse(n.URL)
		if err != nil {
			return nil, err
		}

		tmpArr := strings.Split(u.Scheme, "+")
		name := tmpArr[0]
		if r, ok := routers[name]; ok {
			n.router = r
			n.url = n.URL
			if len(tmpArr) >= 2 {
				u.Scheme = tmpArr[1]
				n.url = u.String()
			}
		} else {
			return nil, fmt.Errorf("invalid scheme:%s", u.Scheme)
		}
	}
	return n.router, nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
