package api

import (
	"context"
	"crud-with-auth/db"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

type Api struct {
	r     *gin.Engine
	db    db.ProviderDB
	Cache *redis.Client
}

var RedisContext = context.Background()

func (api *Api) Start() {
	api.r.GET("/", api.HomeHandler)

	api.Auth()
	api.Articles()
	api.Leaderboard()
	api.OAuth()

	api.r.Run(":4000")
}

func (api Api) HomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func NewAPI(db db.ProviderDB) *Api {
	r := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	if err := rdb.Ping(RedisContext).Err(); err != nil {
		panic("failed ping redis")
	} else {
		log.Println("redis is connected")
	}

	return &Api{r: r, db: db, Cache: rdb}
}

var ProviderAPI = wire.NewSet(NewAPI)
