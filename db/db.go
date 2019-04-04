package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

const ReconnectTimeoutSec = 5 * time.Second

type DataBase struct {
	DB *gorm.DB
	*Config
}

type Config struct {
	User     string
	Password string
	Host     string
	DbName   string
}

func (c *Config) ToUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.DbName)
}

func (db *DataBase) EstablishConnection() {
	log.Println("Connecting to DB ...")

	dbConn, err := gorm.Open("mysql",  db.Config.ToUrl())

	if err != nil {
		log.Println("Failed to connect to database", err)
		time.Sleep(ReconnectTimeoutSec)

		db.EstablishConnection()
	}

	log.Println("Db connection success")
	db.DB = dbConn
}

func (db *DataBase) MonitorConnection() {
	time.Sleep(ReconnectTimeoutSec)

	if db.HasConnection() && db.IsAlive() {
		db.MonitorConnection()

		return
	}

	log.Println("Database connection was not established or is not alive, reconnecting ...")
	go db.EstablishConnection()
}

func (db *DataBase) HasConnection() bool {
	// check if DB is not an empty value
	return db.DB != (&gorm.DB{})
}

func (db *DataBase) CloseConnection() {
	if !db.HasConnection() {
		log.Println("Connection to database is missing, skipping closing ...")
		return
	}

	err := db.DB.Close()

	if err != nil {
		log.Println("Failed to close db connection gracefully", err)
	}
}

func (db *DataBase) IsAlive() bool {
	sqlDB := db.DB.DB()
	err := sqlDB.Ping()

	if err != nil {
		log.Println("Looks like db is not alive, error after ping", err)
		return false
	}

	return true
}
