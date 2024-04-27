package setting

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

var MainDatabase DB

type DB struct {
	Orm orm.Ormer
}

func InitDB() {

	dbType, _ := beego.AppConfig.String("db_type")
	dbMode, _ := beego.AppConfig.String("db_mode")

	dbAlias, _ := beego.AppConfig.String(dbMode + "::db_alias")
	dbUser, _ := beego.AppConfig.String(dbMode + "::db_user")
	dbPwd, _ := beego.AppConfig.String(dbMode + "::db_pwd")
	dbHost, _ := beego.AppConfig.String(dbMode + "::db_host")
	dbPort, _ := beego.AppConfig.String(dbMode + "::db_port")
	dbName, _ := beego.AppConfig.String(dbMode + "::db_name")

	logs.Info("config file", dbType)

	switch dbType {
	case "sqlite3":
		orm.RegisterDataBase(dbAlias, dbType, dbName)
	case "mysql":
		dbCharset, _ := beego.AppConfig.String(dbType + "::db_charset")
		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+
			dbPort+")/"+dbName+"?loc=Asia%2FJakarta&charset="+dbCharset)
	}
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
}
