package api

import (
	"bytes"
	"crud-with-auth/db"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func (api *Api) Articles() {
	group := api.r.Group("/articles")
	group.Use(api.AuthMiddleware)
	group.POST("/", api.CreateArticleHandler)
	group.GET("/", api.ListArticlesHandler)
	group.GET("/export", api.ExportArticlesHandler)
	group.GET("/:id", api.ShowArticleHandler)
	group.PUT("/:id", api.UpdateArticleHandler)
	group.DELETE("/:id", api.DeleteArticleHandler)
	group.POST("/:id/vote/up", api.VoteUpArticleHandler)
}

func (api Api) AuthMiddleware(c *gin.Context) {
	tokenString := c.Request.Header.Get("authorization")

	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized",
		})
		return
	}

	splittedToken := strings.Split(tokenString, " ")
	token := splittedToken[1]

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized",
		})
		return
	}

	var user db.User
	if _, err := user.ValidateToken(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Token invalid or expired",
		})
		return
	}

	c.Next()
}

func (api Api) CreateArticleHandler(c *gin.Context) {
	var article db.Article
	err := c.BindJSON(&article)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Fatalln(err)
		return
	}

	api.db.Storage.Create(&article)

	c.JSON(http.StatusOK, article)
}

func (api Api) ListArticlesHandler(c *gin.Context) {
	var articles []db.Article
	if err := api.db.Storage.Find(&articles).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, articles)
}

func (api Api) ShowArticleHandler(c *gin.Context) {
	var article db.Article
	id := c.Params.ByName("id")

	if err := api.db.Storage.Where("id = ?", id).First(&article).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, article)
}

func (api Api) UpdateArticleHandler(c *gin.Context) {
	var article db.Article
	id := c.Params.ByName("id")

	if err := api.db.Storage.Where("id = ?", id).First(&article).Error; err != nil {
		c.AbortWithStatus(200)
		return
	}

	c.BindJSON(&article)
	api.db.Storage.Save(&article)
	c.JSON(200, article)
}

func (api Api) DeleteArticleHandler(c *gin.Context) {
	var article db.Article
	id := c.Params.ByName("id")

	if err := api.db.Storage.Where("id = ?", id).Delete(&article).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, gin.H{"id": article.ID, "message": "Successfully deleted"})
}

func (api Api) ExportArticlesHandler(c *gin.Context) {
	var articles []db.Article

	if err := api.db.Storage.Find(&articles).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Can't find articles",
		})
	}

	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	heading := []string{
		"No. ", "ID", "Title", "Author",
	}

	if err := w.Write(heading); err != nil {
		log.Fatalln("error writing heading to csv:", err)
	}

	for index, article := range articles {
		row := []string{
			strconv.Itoa(index + 1), fmt.Sprintf("%d", article.ID), article.Title, article.Author,
		}

		if err := w.Write(row); err != nil {
			log.Fatalln("error writing row:", err)
		}
	}

	w.Flush()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=articles.csv")
	c.Writer.Write(b.Bytes())
}

type VoteResponse struct {
	ID    uint    `json:"id"`
	Title string  `json:"title"`
	Votes float64 `json:"votes"`
}

func (api *Api) VoteUpArticleHandler(c *gin.Context) {
	var article db.Article
	id := c.Params.ByName("id")

	if err := api.db.Storage.Where("id = ?", id).First(&article).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "article not found",
		})
		return
	}
	member := fmt.Sprintf("id:%s", id)

	score := api.Cache.ZScore(RedisContext, "articles", member).Val() + 1
	api.Cache.ZAdd(RedisContext, "articles", &redis.Z{
		Score:  score,
		Member: member,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": VoteResponse{
			ID:    article.ID,
			Title: article.Title,
			Votes: score,
		},
	})
}
