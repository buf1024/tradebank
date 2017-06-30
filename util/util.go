package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type CallMeFunc func(int64, interface{})

const time19701900 = 2208988800000000

func SID() string {
	return fmt.Sprintf("SIM-%d-%d",
		time.Now().Unix(), rand.Int63())
}
func CurrentDate() string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
func CurrentUTCMicroSec() int64 {

	t := time.Now()
	t = t.UTC()
	return t.UnixNano()/int64(1000) + time19701900
}
func DateStrToUTCMicroSec(date string) (int64, error) {
	err := error(nil)
	year, err := strconv.Atoi(date[:4])
	if err != nil {
		return -1, err
	}
	month := 0
	if date[4] == '0' {
		month, err = strconv.Atoi(date[5:6])
		if err != nil {
			return -1, err
		}
	} else {
		month, err = strconv.Atoi(date[4:6])
		if err != nil {
			return -1, err
		}
	}
	day, err := strconv.Atoi(date[6:])
	if err != nil {
		return -1, err
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, nil)

	t = t.UTC()
	return t.UnixNano()/int64(1000) + time19701900, nil
}

func CallMeLater(timeout int64, fun CallMeFunc, data interface{}) {
	go func() {
		t := time.NewTimer((time.Duration)((int64)(time.Second) * timeout))
		<-t.C

		fun(timeout, data)
	}()
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

func FtpPut(host string, user string, pass string, localpath string, ftppath string, ftpfile string) error {
	svr, err := ftp.Dial(host)
	if err != nil {
		return err
	}
	defer svr.Quit()

	if err = svr.Login(user, pass); err != nil {

		return err
	}
	if err = svr.ChangeDir(ftppath); err != nil {
		return err
	}
	data, err := ioutil.ReadFile(localpath)
	if err != nil {
		return err
	}
	reader := bytes.NewBuffer(data)
	if err = svr.Stor(ftpfile, reader); err != nil {
		return err
	}

	return nil
}
func FtpGet(host string, user string, pass string, ftppath string, ftpfile string, localpath string) error {
	svr, err := ftp.Dial(host)
	if err != nil {
		return err
	}
	defer svr.Quit()

	if err = svr.Login(user, pass); err != nil {

		return err
	}
	if err = svr.ChangeDir(ftppath); err != nil {
		return err
	}
	rsp, err := svr.Retr(ftpfile)
	if err != nil {
		return err
	}
	defer rsp.Close()

	data, err := ioutil.ReadAll(rsp)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(localpath, data, 0644); err != nil {
		return err
	}

	return nil
}

func PostData(url string, msg []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(msg))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-agent", "Revenco Dream Die HTTP Client")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
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
