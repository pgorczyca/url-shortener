package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	apphttp "github.com/pgorczyca/url-shortener/internal/app/http"
	"github.com/pgorczyca/url-shortener/internal/app/model"
	"github.com/pgorczyca/url-shortener/internal/app/repository"
	"github.com/pgorczyca/url-shortener/internal/app/shortener"
	"github.com/pgorczyca/url-shortener/internal/app/utils"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const expirationDuration = time.Hour * 8765 // links are valid for 1 year
var prefixUrl = utils.GetConfig().PrefixUrl

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
		internalServerErrorResponse(c)
		utils.Logger.Error("Not able to generate new short", zap.Error(err))
		return
	}
	url := model.Url{
		Long:      req.Long,
		Short:     short,
		ExpiredAt: time.Now().Add(expirationDuration),
		CreatedAt: time.Now(),
	}
	err = repo.Add(context.TODO(), url)
	if err != nil {
		internalServerErrorResponse(c)
		utils.Logger.Error("Not able to insert to repository.", zap.Error(err))
		return
	}

	res := apphttp.UrlResponse{Long: req.Long,
		Short:     strings.Join([]string{prefixUrl, url.Short}, "/"),
		CreatedAt: url.CreatedAt,
		ExpiredAt: url.ExpiredAt,
	}
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

	c.Redirect(http.StatusMovedPermanently, url.Long)

}
