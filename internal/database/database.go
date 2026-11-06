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
		"host=%v user=%v password=%v dbname=%v TimeZone=%v",
		*cfg.DBHost, *cfg.DBUser, *cfg.DBPassword, *cfg.DBName, *cfg.DBTimeZone,
	)
	dialec := postgres.Open(dsn)
	gormCfg := new(gorm.Config{
		TranslateError: true,
	})
	if db, err = gorm.Open(dialec, gormCfg); err != nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	sqlDB.SetMaxOpenConns(int(*cfg.DBMaxOpenCons))
	sqlDB.SetConnMaxIdleTime(*cfg.DBConnMaxIdleTime)

	return
}
