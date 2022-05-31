package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	apphttp "github.com/pgorczyca/url-shortener/internal/app/http"
	"github.com/pgorczyca/url-shortener/internal/app/model"
	"github.com/pgorczyca/url-shortener/internal/app/repository"
	"github.com/pgorczyca/url-shortener/internal/app/shortener"

	"github.com/gin-gonic/gin"
)

func CreateUrl(c *gin.Context, repo repository.UrlRepository, sg *shortener.ShortGenerator) {
	jsonUrl, _ := ioutil.ReadAll(c.Request.Body)
	var req apphttp.CreateUrlRequest
	json.Unmarshal(jsonUrl, &req)
	validationerr, err := validateCreate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": validationerr,
		})

		return
	}
	short, err := sg.GetShort()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err,
		})
		return
	}
	url := model.Url{
		Long:      "https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/",
		Short:     short,
		ExpiredAt: time.Now().Add(time.Hour * 6),
		CreatedAt: time.Now(),
	}
	repo.Add(context.TODO(), url)

	res := apphttp.UrlResponse{Long: req.Long, Short: short, CreatedAt: time.Now(), ExpiredAt: time.Now().Add(time.Hour * 6)}
	c.JSON(http.StatusCreated, res)

}

func GetUrl(c *gin.Context, repo repository.UrlRepository) {
	url, err := repo.GetByShort(c, c.Param("short"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "record not found",
		})
		return
	}
	res := apphttp.UrlResponse{Long: url.Long, Short: url.Short, CreatedAt: url.CreatedAt, ExpiredAt: url.ExpiredAt}
	c.JSON(http.StatusOK, res)
}
