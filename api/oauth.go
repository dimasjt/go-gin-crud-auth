package api

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OAuthAPI struct {
	Api
	Config *oauth2.Config
}

func (api *Api) OAuth() {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:4000/oauth/github/callback",
	}

	oauthApi := OAuthAPI{
		Api:    *api,
		Config: conf,
	}

	oauthApi.BuildRoutes()
}

func (oauth *OAuthAPI) BuildRoutes() {
	group := oauth.Api.r.Group("/oauth/github")
	group.GET("/", oauth.GithubAuthorizeHandler)
	group.GET("/callback", oauth.GithubCallbackHandler)
}

func (oauth *OAuthAPI) GithubAuthorizeHandler(c *gin.Context) {
	url := oauth.Config.AuthCodeURL("abcd")

	c.Redirect(http.StatusTemporaryRedirect, url)
}

type GithubUserResponse struct {
	ID int `json:"id"`
}

func (oauth *OAuthAPI) GithubCallbackHandler(c *gin.Context) {
	code, ok := c.GetQuery("code")

	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "code is empty",
		})
		return
	}

	ctx := context.Background()
	token, err := oauth.Config.Exchange(ctx, code)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer res.Body.Close()

	c.JSON(http.StatusOK, &token)
}
