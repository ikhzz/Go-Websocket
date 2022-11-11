package config

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"clean_arch_v2/models"
)

var SqlConn *gorm.DB

// InitMysql for connect to mysql database
func InitMysql() *gorm.DB {
	// get database variale
	dbHost := viper.GetString(`database.mysql.host`)
	dbPort := viper.GetString(`database.mysql.port`)
	dbUser := viper.GetString(`database.mysql.user`)
	dbPass := viper.GetString(`database.mysql.pass`)
	dbName := viper.GetString(`database.mysql.db`)

	// set dsn url and try connection
	basedsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", viper.GetString("timezone"))
	dsn := fmt.Sprintf("%s?%s", basedsn, val.Encode())
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	SqlConn = conn

	// advance setting
	// sqldb,err :=db.DB()
	// sqldb.SetMaxOpenConns(10)
	// sqldb.SetMaxIdleConns(10)
	// sqldb.SetConnMaxLifetime(10)
	conn.AutoMigrate(
		&models.AuthModel{}, 
		&models.GroupModel{},
		&models.GroupMemberModel{},
		&models.GeneralMessageListModel{},
		&models.GroupMessageListModel{},
	)
	return conn
}
	