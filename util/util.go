package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type CallMeFunc func(int64, interface{})

func SID() string {
	return fmt.Sprintf("SIM-%d-%d",
		time.Now().Unix(), rand.Int63())
}

func CallMeLater(timeout int64, fun CallMeFunc, data interface{}) {
	t := time.NewTimer((time.Duration)((int64)(time.Second) * timeout))
	<-t.C

	fun(timeout, data)
}
func CurrentDate() string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
func GetSplitData(rev string, key string) string {
	if rev == "" {
		return ""
	}

	revVec := strings.Split(rev, "|")
	for _, v := range revVec {
		index := strings.Index(v, key)
		if index >= 0 {
			val := v[len(key):]

			return strings.TrimSpace(val)
		}
	}

	return ""
}
func PostData(url string, msg []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(msg))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-agent", "Revenco Dream Die HTTP Client")
	client := http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil

}

func init() {
	rand.NewSource(time.Now().UnixNano())
}
