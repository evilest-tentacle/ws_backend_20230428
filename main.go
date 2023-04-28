package main

import (
	"chatgpt-service/chat"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func kee_parser(config *chat.Config, logger *chat.Logger) {
	logger.LogInfo("current_key is ", config.Kee)
	old_kee := config.Kee
	new_kee := strings.Replace(old_kee, "hello#", "sk-", -1)
	new_kee = strings.Replace(new_kee, "#world", "", -1)
	config.Kee = new_kee
}

func main() {
	logger := chat.Logger{}
	logger.LoggerInit()

	bs, err := os.ReadFile("config.yaml")
	if err != nil {
		err = fmt.Errorf("read file config.yaml error: %s", err.Error())
		logger.LogError(err.Error())
		return
	}
	var config chat.Config
	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		err = fmt.Errorf("parse config.yaml error: %s", err.Error())
		logger.LogError(err.Error())
		return
	}
	if config.Kee == "" {
		logger.LogError(fmt.Sprintf("apiKey is empty"))
		return
	}
	kee_parser(&config, &logger)
	//fmt.Println(config.Kee)
	var found bool
	for _, model := range chat.GPTModels {
		if model == config.Model {
			found = true
			break
		}
	}
	if !found {
		logger.LogError(fmt.Sprintf("model not exists"))
		return
	}

	api := chat.Api{
		Config: config,
		Logger: logger,
	}
	r := gin.Default()
	if config.Cors {
		cfg := cors.DefaultConfig()
		cfg.AllowAllOrigins = true
		cfg.AllowHeaders = []string{"content-type"}
		r.Use(cors.New(cfg))
	}

	groupApi := r.Group("/api")
	groupApi.Static("/assets", "assets")
	groupWs := groupApi.Group("/ws")
	groupWs.GET("chat", api.WsChat)

	logger.LogInfo("chatGPT query service start")
	err = r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		err = fmt.Errorf("run service error: %s", err.Error())
		logger.LogPanic(err.Error())
		return
	}
}
