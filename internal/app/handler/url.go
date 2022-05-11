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

	"github.com/gin-gonic/gin"
)

func CreateUrl(repo repository.UrlRepository) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		jsonUrl, _ := ioutil.ReadAll(c.Request.Body)
		var req apphttp.CreateUrlRequest
		json.Unmarshal(jsonUrl, &req)
		validationerr, err := validate(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validationerr,
			})

			return
		}

		url := model.Url{
			Long:      "https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/write-operations/insert/",
			Short:     "ASt",
			ExpiredAt: time.Now().Add(time.Hour * 6),
			CreatedAt: time.Now(),
		}
		repo.Add(context.TODO(), url)

		res := apphttp.CreateUrlResponse{Long: req.Long, Short: "34s", CreatedAt: time.Now()}
		c.JSON(http.StatusOK, res)
	}
	return gin.HandlerFunc(fn)
}
