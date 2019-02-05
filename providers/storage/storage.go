package storage

import (
	"database/sql"
	"fmt"

	"translator/config"
	"translator/providers"
	"translator/providers/yandex"

	_ "github.com/go-sql-driver/mysql"
)

// IStorage storage interface
type IStorage interface {
	MySQL() *sql.DB
	YandexTranslator() providers.TranslateProvider
}

// Storage contain connects for db
type Storage struct {
	mysql  *sql.DB
	yandex *yandex.Client
}

// Init constructor
func Init(c *config.Config) *Storage {
	return &Storage{
		mysql:  mustMySQLConnect(c.MySQL),
		yandex: mustYandexTranslator(c.Translators.Yandex),
	}
}

// MySQL return MySQL connection
func (s *Storage) MySQL() *sql.DB {
	return s.mysql
}

// YandexTranslator return yandex translator
func (s *Storage) YandexTranslator() providers.TranslateProvider {
	return s.yandex
}

func mustMySQLConnect(c config.MySQL) *sql.DB {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%dms",
			c.User,
			c.Password,
			c.Address,
			c.DBName,
			3000,
		))

	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	println("db ping success") // todo

	return db
}

func mustYandexTranslator(c config.Yandex) *yandex.Client {
	return yandex.New(c.ApiKey, c.ApiVer)
}
