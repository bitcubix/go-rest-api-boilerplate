package server

import (
	"github.com/gabrielix29/go-rest-api/pkg/logger"
	"gorm.io/driver/postgres"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
	Router   *mux.Router
	Database *gorm.DB
}

func New() *Server {
	var server Server
	server.Router = mux.NewRouter()
	return &server
}

func (s *Server) Run() {
	s.InitDatabase()
	s.initServices()
	addr := viper.GetString("server.host") + ":" + viper.GetString("server.port")
	httpserver := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}
	http2.ConfigureServer(httpserver, &http2.Server{})
	logger.Info("HTTP Server started listening on ", addr)
	logger.Fatal(httpserver.ListenAndServe())
}

func (s *Server) InitDatabase() {
	var err error
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	dbname := viper.GetString("database.name")

	switch viper.GetString("database.driver") {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
		s.Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		break
	case "postgres":
		dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", username, password, dbname, port, host)
		s.Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		break
	default:
		logger.Fatal("invalid database driver please use 'mysql' or 'postgresql'")
	}

	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Debug("Connected to database")
	}
}
