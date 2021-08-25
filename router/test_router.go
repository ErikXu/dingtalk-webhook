package router

import (
	"io/ioutil"
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	bytes, _ := ioutil.ReadAll(c.Request.Body)
	content := string(bytes[:])
	log.Println(content)

	c.JSON(http.StatusOK, content)
}
