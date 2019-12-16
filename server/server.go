package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/afairon/smeter/config"
	"github.com/afairon/smeter/internal/impl"
	"github.com/afairon/smeter/service"
	"google.golang.org/grpc"
)

// Server structure to access config and database.
type Server struct {
	DB     *sql.DB
	Server *grpc.Server
	Config *config.Config
}

// NewServer create new connection to database and
// create a new instance of grpc server.
func NewServer(c *config.Config) *Server {
	uri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.UserName,
		c.Database.Password,
		c.Database.DBName,
	)

	if c.Database.SSL {
		uri = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			c.Database.Host,
			c.Database.Port,
			c.Database.UserName,
			c.Database.Password,
			c.Database.DBName,
		)
	}

	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer(createServerUnaryInterceptor(), createServerStreamInterceptor())

	service.RegisterDeviceServiceServer(srv, impl.NewDeviceServiceImpl(db))
	service.RegisterSensorServiceServer(srv, impl.NewSensorServiceImpl(db))
	service.RegisterPowerServiceServer(srv, impl.NewPowerServiceImpl(db))
	service.RegisterTemperatureServiceServer(srv, impl.NewTemperatureServiceImpl(db))
	service.RegisterHumidityServiceServer(srv, impl.NewHumidityServiceImpl(db))

	server := Server{
		DB:     db,
		Server: srv,
		Config: c,
	}

	return &server
}

// Run listens to socket.
func (s *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening %s:%s", s.Config.Server.Host, s.Config.Server.Port)

	if err := s.Server.Serve(listener); err != nil {
		panic(err)
	}
}
