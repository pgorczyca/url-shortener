package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	apphttp "github.com/pgorczyca/url-shortener/internal/app/http"

	"github.com/gin-gonic/gin"
)

func CreateUrl(c *gin.Context) {
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
	res := apphttp.CreateUrlResponse{Long: req.Long, Short: "34s", CreatedAt: time.Now()}
	c.JSON(http.StatusOK, res)
}
