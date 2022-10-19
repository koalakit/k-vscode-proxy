package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var mysqldb *gorm.DB

func MySQLInit(url string, debug bool) (err error) {
	mysqldb, err = gorm.Open("mysql", url)
	if err != nil {
		return
	}

	mysqldb.LogMode(debug)
	mysqldb.DB().SetMaxIdleConns(5)
	mysqldb.DB().SetMaxOpenConns(10)

	mysqldb.AutoMigrate(new(UserData))

	go func() {
		for {
			<-time.After(2 * time.Minute)

			if err := mysqldb.DB().Ping(); err != nil {
				log.Println("[MySQL] ping:", err)
			}
		}
	}()

	if err = mysqldb.DB().Ping(); err != nil {
		log.Fatal("[MySQL] ping:", err)
	} else {
		log.Println("[MySQL] success.")
	}

	return
}

func MySQLPing() (err error) {
	err = mysqldb.DB().Ping()
	return
}
