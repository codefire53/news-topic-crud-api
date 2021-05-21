package infrastructures

import (
	"fmt"
	"news-topic-api/helpers"
	"news-topic-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

//InitDB ...
func InitDB() {
	conn, err := dbSetup()
	if err != nil {
		panic(err)
	}
	db = conn
	doMigration()
}

//GetDB ...
func GetDB() *gorm.DB {
	return db
}

//SetDB ...
func SetDB(newDB *gorm.DB) {
	db = newDB
}

func doMigration() {
	db.AutoMigrate(&models.News{})
	db.AutoMigrate(&models.Tag{})
}

func dbSetup() (*gorm.DB, error) {
	username := helpers.GetEnv("postgres_username", "postgres")
	password := helpers.GetEnv("postgres_password", "postgres")
	dbName := helpers.GetEnv("postgres_dbname", "postgres")
	host := helpers.GetEnv("postgres_host", "localhost")
	dbParams := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", username, password, dbName, host)
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  helpers.GetEnv("DATABASE_URL", dbParams), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                                     // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})
	return conn, err
}
