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
	_dirHttpDelivery "github.com/sicozz/papyrus/dir/delivery/http"
	_dirRepo "github.com/sicozz/papyrus/dir/repository/postgres"
	_dirUsecase "github.com/sicozz/papyrus/dir/usecase"
	_pFileHttpDelivery "github.com/sicozz/papyrus/pfile/delivery/http"
	_pFileRepo "github.com/sicozz/papyrus/pfile/repository/postgres"
	_pFileUsecase "github.com/sicozz/papyrus/pfile/usecase"
	_roleRepo "github.com/sicozz/papyrus/role/repository/postgres"
	_userHttpDelivery "github.com/sicozz/papyrus/user/delivery/http"
	_userRepo "github.com/sicozz/papyrus/user/repository/postgres"
	_userUsecase "github.com/sicozz/papyrus/user/usecase"
	_userStateRepo "github.com/sicozz/papyrus/user_state/repository/postgres"
	"github.com/sicozz/papyrus/utils"
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
	err := utils.InitFsDir()
	if err != nil {
		log.Fatal("Failed to initialize fs directory -> ", err)
	}

	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
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

	// This enables debuggin. Change it later to an adequte time
	// timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second * 300
	rr := _roleRepo.NewPostgresRoleRepository(dbConn)
	ur := _userRepo.NewPostgresUserRepository(dbConn)
	usr := _userStateRepo.NewPostgresUserStateRepository(dbConn)
	uu := _userUsecase.NewUserUsecase(ur, rr, usr, timeoutContext)
	_userHttpDelivery.NewUserHandler(e, uu)

	dr := _dirRepo.NewPostgresDirRepository(dbConn)
	dPfr := _pFileRepo.NewPostgresPFileRepository(dbConn)
	du := _dirUsecase.NewDirUsecase(dr, dPfr, timeoutContext)
	_dirHttpDelivery.NewDirHandler(e, du)

	pfr := _pFileRepo.NewPostgresPFileRepository(dbConn)
	pfDr := _dirRepo.NewPostgresDirRepository(dbConn)
	pfUr := _userRepo.NewPostgresUserRepository(dbConn)
	pfu := _pFileUsecase.NewPFileUsecase(pfr, pfDr, pfUr, timeoutContext)
	_pFileHttpDelivery.NewPFileHandler(e, pfu)

	e.Logger.Fatal(e.Start(":9090"))
	/**
	* TODO: Improve comments
	* TODO: Add unit testing for everything created
	* TODO: Make a deep check and desing of fs constraints, enforce them in DB and in app-code
	**/
}
