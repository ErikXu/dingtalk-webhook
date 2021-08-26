package router

import (
	"dingtalk-webhook/config"
	"io/ioutil"
	"net/http"

	"log"

	util "dingtalk-webhook/util"

	"github.com/gin-gonic/gin"
)

func print(c *gin.Context) {
	bytes, _ := ioutil.ReadAll(c.Request.Body)
	content := string(bytes[:])
	log.Println(content)

	c.JSON(http.StatusOK, content)
}

func send(c *gin.Context) {
	err := util.SendMarkdownMsg(config.AppConfig.Dingtalk.Webhook, config.AppConfig.Dingtalk.Secret, "Test", "Hello World!")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed")
	}

	c.JSON(http.StatusOK, "success")
}
