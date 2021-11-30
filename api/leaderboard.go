package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func (api *Api) Leaderboard() {
	group := api.r.Group("/leaderboard")
	group.GET("/", api.CurrentLeaderboard)
}

func (api *Api) CurrentLeaderboard(c *gin.Context) {
	result, err := api.Cache.ZRangeByScoreWithScores(RedisContext, "articles", &redis.ZRangeBy{
		Min: "0",
		Max: "+inf",
	}).Result()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed to get leaderboard",
		})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
