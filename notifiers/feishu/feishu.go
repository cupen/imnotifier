package feishu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gamedev-embers/imnotifier/models"
)

func Send(url string, msg interface{}, timeout ...time.Duration) error {
	if url == "" {
		return fmt.Errorf("empty url")
	}
	if msg == nil {
		return fmt.Errorf("nil message")
	}

	realTimeout := 6 * time.Second
	_ = realTimeout
	if len(timeout) > 0 {
		realTimeout = timeout[0]
		_ = realTimeout
		// TODO: timeout
	}

	switch _msg := msg.(type) {
	case *models.Text:
		return sendText(url, _msg.Content)
	default:
		return fmt.Errorf("invalid message: %+v", msg)
	}
}

func sendText(_url string, text string) error {
	parseSignKey := func(_url string) (string, string) {
		tmpArr := strings.Split(_url, "?")
		if len(tmpArr) == 1 {
			return _url, ""
		} else if len(tmpArr) == 2 {
			_url = tmpArr[0]
			q, err := url.ParseQuery(tmpArr[1])
			if err != nil {
				panic(err)
			}
			return _url, q.Get("signKey")
		} else {
			panic(fmt.Errorf("invalid url: %s", _url))
		}
	}
	sign := func(secret string, timestamp int64) string {
		stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)
		var data []byte
		h := hmac.New(sha256.New, []byte(stringToSign))
		_, err := h.Write(data)
		if err != nil {
			panic(err)
		}
		return base64.StdEncoding.EncodeToString(h.Sum(nil))
	}
	if _url == "" {
		return nil
	}
	_url, signKey := parseSignKey(_url)
	body := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": text,
		},
	}

	if signKey != "" {
		ts := time.Now().Unix()
		body["timestamp"] = ts
		body["sign"] = sign(signKey, ts)
	}
	return postJson(_url, body)
}

// postJson ...
func postJson(url string, body map[string]interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return post(url, "application/json", b)
}

// post ...
func post(url string, contentType string, body []byte) error {
	c := http.Client{}
	buf := bytes.NewBuffer(body)
	resp, err := c.Post(url, contentType, buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid request. url:%s resp:%d", url, resp.StatusCode)
	}
	return nil
}
