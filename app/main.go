package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	_roleRepo "github.com/sicozz/papyrus/role/repository/postgres"
	_userHttpDelivery "github.com/sicozz/papyrus/user/delivery/http"
	_userRepo "github.com/sicozz/papyrus/user/repository/postgres"
	_userUsecase "github.com/sicozz/papyrus/user/usecase"
	_userStateRepo "github.com/sicozz/papyrus/user_state/repository/postgres"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("sslmode", "disable")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`postgres`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	e.Use(middleware.CORS())

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	rr := _roleRepo.NewPostgresRoleRepository(dbConn)
	ur := _userRepo.NewPostgresUserRepository(dbConn)
	usr := _userStateRepo.NewPostgresUserStateRepository(dbConn)
	uu := _userUsecase.NewUserUsecase(ur, rr, usr, timeoutContext)
	_userHttpDelivery.NewUserHandler(e, uu)
	e.Logger.Fatal(e.Start(":9090"))
}
