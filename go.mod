module news-topic-api

// +heroku goVersion go1.16
go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/google/uuid v1.1.2 //
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.7.3
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.10.0
	github.com/magiconair/properties v1.8.5
	github.com/stretchr/testify v1.7.0
	google.golang.org/api v0.40.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.3
)
