package database

import (
	"errors"
	"fmt"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrCreateDB = errors.New("CREATE DATABASE ERROR")

func NewDB(cfg *config.Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"user=%v password=%v dbname=%v",
		cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	dialec := postgres.Open(dsn)
	gormCfg := new(gorm.Config{
		TranslateError: true,
	})

	db, err = gorm.Open(dialec, gormCfg)
	return
}
