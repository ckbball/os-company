package cmd

import (
  "context"
  "database/sql"
  "flag"
  "fmt"
  "os"
  "strconv"

  "github.com/ckbball/os-company/pkg/logger"
  companyGrpc "github.com/ckbball/os-company/pkg/protocol/grpc"
  v1 "github.com/ckbball/os-company/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
  // gRPC server start parameters section
  // gRPC is TCP port to listen by gRPC server
  GRPCPort string

  // the port to listen for http calls
  HTTPPort string

  // DB Datastore parameters section
  // DatastoreDBHost is host of database
  DatastoreDBHost string
  // DatastoreDBUser is username to connect to database
  DatastoreDBUser string
  // DatastoreDBPassword password to connect to database
  DatastoreDBPassword string
  // DatastoreDBSchema is schema of database
  DatastoreDBSchema string
  // address for single redis node
  RedisAddress string

  // LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
  LogLevel int
  // LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
  LogTimeFormat string

  // user service address
  JobSvcAddress string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
  ctx := context.Background()

  // get configuration
  var cfg Config
  flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
  flag.StringVar(&cfg.HTTPPort, "http-port", "", "http port to bind")
  flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
  flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
  flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
  flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
  flag.StringVar(&cfg.RedisAddress, "redis-address", "", "Redis address")
  flag.Parse()

  if len(cfg.GRPCPort) == 0 {
    cfg.GRPCPort = os.Getenv("GRPC_PORT")
    cfg.HTTPPort = os.Getenv("HTTP_PORT")
    cfg.DatastoreDBHost = os.Getenv("DB_HOST")
    cfg.DatastoreDBUser = os.Getenv("DB_USER")
    cfg.DatastoreDBPassword = os.Getenv("DB_PASSWORD")
    cfg.DatastoreDBSchema = os.Getenv("DB_SCHEMA")
    cfg.RedisAddress = os.Getenv("REDIS_ADDRESS")
    cfg.JobSvcAddress = os.Getenv("JOB_ADDRESS")
    cfg.LogLevel, _ = strconv.Atoi(os.Getenv("LOG_LEVEL"))
    cfg.LogTimeFormat = os.Getenv("LOG_TIME")
  }

  if len(cfg.GRPCPort) == 0 {
    return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
  }

  if len(cfg.HTTPPort) == 0 {
    return fmt.Errorf("invalid TCP port for http server: '%s'", cfg.HTTPPort)
  }

  // SET up mongo client
  // retry := false
  clientOptions := options.Client().ApplyURI(cfg.MongoAddress)
  client, err := mongo.Connect(context.TODO(), clientOptions)
  if err != nil {
    return err
  }
  collection := client.Database(cfg.MongoName).Collection(cfg.MongoCollection)

  // create repository
  repository := v1.NewCompanyRepository(collection)

  // create auth service
  tokenService := v1.NewTokenService()

  // pass in fields of handler directly to method
  v1API := v1.NewCompanyServiceServer(repository, tokenService) // may need to add Job Service address

  // initialize logger
  if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
    return fmt.Errorf("failed to initialize logger: %v", err)
  }

  return companyGrpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
