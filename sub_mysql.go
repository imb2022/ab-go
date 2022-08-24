package ab

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/xwi88/kit4go/mysql"
	"github.com/xwi88/log4go"

	"github.com/imb2022/ab-go/scheme"
)

var (
	mysqlClient *mysql.Client
)

func initMySql(cfg MySql) (err error) {
	charset := mysql.Charset(cfg.Charset)
	collation := mysql.Collation(cfg.Collation)
	maxOpenConnections := mysql.MaxOpenConnections(cfg.MaxOpenConnections)
	maxIdleConnections := mysql.MaxIdleConnections(cfg.MaxIdleConnections)
	debug := mysql.Debug(cfg.Debug)
	parseTime := mysql.ParseTime(cfg.ParseTime)
	loc := mysql.Loc(cfg.LOC)
	tls := mysql.TLS(cfg.TLS)
	connMaxLifetime := mysql.ConnMaxLifetime(cfg.ConnMaxLifetime)
	timeout := mysql.Timeout(cfg.Timeout)
	readTimeout := mysql.ReadTimeout(cfg.ReadTimeout)
	writeTimeout := mysql.WriteTimeout(cfg.WriteTimeout)

	if mysqlClient, err = mysql.Dial(cfg.Addr, cfg.User, cfg.Password, cfg.DBName,
		charset, collation, maxOpenConnections, maxIdleConnections,
		debug, parseTime, loc, tls, connMaxLifetime,
		timeout, readTimeout, writeTimeout); err != nil {
		return err
	}

	return nil
}

func closeMySql() error {
	if mysqlClient != nil {
		return mysqlClient.Close()
	}
	return nil
}

func loadABScheme(appNameOrFlag string) {
	_sql := fmt.Sprintf("SELECT app_flag, config from experiment_config WHERE app_flag = '%v' order by created_at desc limit 1", appNameOrFlag)
	row := mysqlClient.MDB.QueryRow(_sql)
	if row != nil {
		var appFlag sql.NullString
		var config sql.NullString
		err := row.Scan(&appFlag, &config)
		if err != nil {
			log4go.Error("[loadABScheme] app_flag:%v, err:%v", appNameOrFlag, err.Error())
			return
		}
		log4go.Info("[loadABScheme] app_flag:%v, config:%v", appNameOrFlag, config)
		if len(config.String) > 0 {
			configDataBytes, err := base64.URLEncoding.DecodeString(config.String)
			if err != nil {
				log4go.Error("[loadABScheme] app_flag:%v, err:%v", appNameOrFlag, err.Error())
				return
			}
			var abScheme *scheme.ABScheme
			err = json.Unmarshal(configDataBytes, &abScheme)
			if err != nil {
				log4go.Error("[loadABScheme] app_flag:%v, err:%v", appNameOrFlag, err.Error())
				return
			}
			log4go.Info("[loadABScheme] app_flag:%v, scheme:%v", appNameOrFlag, abScheme)
			if abScheme != nil {
				appABScheme = abScheme
			}
		}
	}
}
