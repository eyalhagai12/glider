package main

import (
	"fmt"
	"glider/api"
	"glider/workerpool"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DatabaseURL             string
	DatabaseMaxIdleConns    int64
	DatabaseMaxOpenConns    int64
	DatabaseConnMaxLifetime int64
}

func loadConfig() Config {
	godotenv.Load()

	var conf Config

	conf.DatabaseURL, _ = os.LookupEnv("DB_CONNECTION_STRING")
	conf.DatabaseMaxIdleConns, _ = strconv.ParseInt(os.Getenv("DB_MAX_IDLE_CONNS"), 10, 64)
	conf.DatabaseMaxOpenConns, _ = strconv.ParseInt(os.Getenv("DB_MAX_OPEN_CONNS"), 10, 64)
	conf.DatabaseConnMaxLifetime, _ = strconv.ParseInt(os.Getenv("DB_CONN_MAX_LIFETIME"), 10, 64)

	return conf
}

func main() {
	cfg := loadConfig()

	r := gin.Default()
	os.Setenv("DOCKER_BUILDKIT", "1")

	logger := log.Default()
	logger.Printf("connecting to database %s", cfg.DatabaseURL)
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}
	db.SetMaxIdleConns(int(cfg.DatabaseMaxIdleConns))
	db.SetMaxOpenConns(int(cfg.DatabaseMaxOpenConns))
	db.SetConnMaxLifetime(time.Duration(cfg.DatabaseConnMaxLifetime) * time.Minute)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	workerPool := workerpool.NewWorkerPool(10, 10)
	defer workerPool.Close()

	initViews(r, db, cli, workerPool)

	r.Run()
}

func initViews(r *gin.Engine, db *sqlx.DB, cli *client.Client, workerPool *workerpool.WorkerPool) {
	deployHandlers := api.NewDeployHandlers(db, workerPool, cli)
	nodesHandlers := api.NewNodeHandlers(db)

	r.POST("/deploy", api.HandlerFromFunc(deployHandlers.Deploy, http.StatusAccepted))
	r.POST("/nodes/register", api.HandlerFromFunc(nodesHandlers.RegisterNewNode, http.StatusCreated))
}
