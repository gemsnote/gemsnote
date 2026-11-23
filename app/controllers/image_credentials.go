package controllers

import (
	"strings"
	"time"

	"github.com/gemsnote/gemsnote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
)

const defaultImageTokenTTL = 10 * time.Minute

func imageSecret() string {
	if revel.Config == nil {
		return ""
	}
	if configured := strings.TrimSpace(revel.Config.StringDefault("image.secret", "")); configured != "" {
		return configured
	}
	appSecret, _ := revel.Config.String("app.secret")
	return lea.DeriveImageSecret(appSecret)
}

// ImageTokenTTL returns the configured lifetime of an image capability.
func ImageTokenTTL() time.Duration {
	if revel.Config == nil {
		return defaultImageTokenTTL
	}
	raw := strings.TrimSpace(revel.Config.StringDefault("image.token.expires", defaultImageTokenTTL.String()))
	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		return defaultImageTokenTTL
	}
	return value
}

// SignedImagePath builds a short-lived authorised image URL for the viewer.
func SignedImagePath(fileId, userId string) string {
	if !bson.IsObjectIdHex(fileId) || !bson.IsObjectIdHex(userId) {
		return ""
	}
	token := lea.ImageToken(imageSecret(), userId, fileId, time.Now(), ImageTokenTTL())
	if token == "" {
		return ""
	}
	return "/api2/file/getImage?fileId=" + fileId + "&uid=" + userId + "&token=" + token
}

// ImageCapabilityUser returns the viewer identity attested by a signed image
// URL. Permission checks still run against the file afterwards.
func ImageCapabilityUser(c *revel.Controller, fileId string) string {
	if c == nil || c.Params == nil || c.Params.Values == nil || !bson.IsObjectIdHex(fileId) {
		return ""
	}
	userId := strings.TrimSpace(c.Params.Values.Get("uid"))
	token := strings.TrimSpace(c.Params.Values.Get("token"))
	if !bson.IsObjectIdHex(userId) || token == "" {
		return ""
	}
	if !lea.ValidImageToken(imageSecret(), userId, fileId, token, time.Now(), ImageTokenTTL()) {
		return ""
	}
	return userId
}
