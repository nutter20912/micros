package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"micros/config"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"

	"github.com/golang-migrate/migrate/v4"
	dm "github.com/golang-migrate/migrate/v4/database/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() {
	newDatabase()

	dbInstance, _ := db.DB()
	migration(dbInstance)
}

func Get() *gorm.DB {
	return db
}

func newDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.database"),
	)

	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Millisecond * 50,
			LogLevel:                  gormLogger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)

	var err error
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("init db")

	return db
}

func migration(db *sql.DB) {
	driver, _ := dm.WithInstance(db, &dm.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:%s/database/migrations/%s", config.BasePath, viper.GetString("db.database")),
		viper.GetString("db.database"),
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("migration done")
}
