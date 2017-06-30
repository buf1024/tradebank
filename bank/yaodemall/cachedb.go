package main

import (
	"database/sql"
	"tradebank/util"

	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type YaodeMallDB struct {
	mall *YaodeMall
	db   *sql.DB
	lock sync.Locker
}

type InoutLog struct {
	extflow     string
	iotype      int
	amount      float64
	status      int
	operatetime int64
	checkdate   string
	payway      int
}

func (y *YaodeMallDB) Init(file string) error {
	sqlStr := `
	begin;
	create table inout_log(
		extflow varchar(64) primary key not null,
		type number(8, 0) not null,
		amount number(16, 4) not null,
		status number(8, 0) not null,
		operatetime number(20, 0) not null,
		checkdate varchar(32) not null,
		payway number(8, 0)
	);
	commit;
	`
	var err error
	y.db, err = sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	_, err = y.db.Exec("select count(1) from inout_log")
	if err != nil {
		y.mall.Log.Info("db not exists, create one\n")
		_, err = y.db.Exec(sqlStr)
		if err != nil {
			y.mall.Log.Error("create database failed\n")
			return err
		}
		y.mall.Log.Info("database create\n")
	}
	return nil

}
func (y *YaodeMallDB) Close() {
	y.lock.Lock()
	defer y.lock.Unlock()
	y.db.Close()
}

func (y *YaodeMallDB) QueryCheckLog(operatetime int64) ([]InoutLog, error) {
	y.lock.Lock()
	defer y.lock.Unlock()
	y.mall.Log.Debug("QueryLog, t=%d\n", operatetime)
	rows, err := y.db.Query("select extflow, type, amount, status, operatetime, checkdate, payway from inout_log where status = 0 and operatetime <= ?",
		operatetime)

	if err != nil {
		y.mall.Log.Error("QueryCheckLog error = %s\n", err.Error())
		return nil, err
	}

	var result []InoutLog

	for rows.Next() {
		var tmpRes InoutLog
		if err := rows.Scan(&tmpRes.extflow, &tmpRes.iotype, &tmpRes.amount, &tmpRes.status, &tmpRes.operatetime, &tmpRes.checkdate, &tmpRes.payway); err != nil {
			y.mall.Log.Error("Scan error =  %s\n", err)
			return nil, nil
		}
		result = append(result, tmpRes)
	}
	//y.mall.Log.Debug("result: %v\n", result)
	return result, err
}
func (y *YaodeMallDB) InsertLog(lg InoutLog) error {
	y.lock.Lock()
	defer y.lock.Unlock()
	now := util.CurrentUTCMicroSec()
	y.mall.Log.Debug("now->%d\n", now)

	_, err := y.db.Exec("insert into inout_log(extflow, type, amount, status, operatetime, checkdate, payway) values(?, ?, ?, ?, ?, ?, ?)",
		lg.extflow, lg.iotype, lg.amount, 0, now, "", lg.payway)
	if err != nil {
		y.mall.Log.Error("InsertLog failed, error = %s\n", err)
		return err
	}
	return err
}
func (y *YaodeMallDB) UpdateLog(extflow string, status int, checkdate string) error {
	y.lock.Lock()
	defer y.lock.Unlock()
	_, err := y.db.Exec("update inout_log set status = ?, checkdate = ? where extflow = ?",
		status, checkdate, extflow)
	if err != nil {
		y.mall.Log.Error("UpdateLog failed, error = %s\n", err)
		return err
	}
	return err
}
