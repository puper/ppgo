package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/puper/p2pwatch/app"
	"github.com/puper/p2pwatch/models"
)

func Index(c *routing.Context) error {
	var posts []models.Post
	app.ChooseDB(models.PostModel, false).Select("id, title, content").Find(&posts)
	c.Write(posts)
	r, _ := app.MustGetRedis("redis").Do("get", "a")
	c.Write(r)
	log.Debug("request: %v", r)
	return nil
}
