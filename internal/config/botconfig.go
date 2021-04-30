package config

import (
	errutil2 "github.com/Skycoin/git-telegram-bot/pkg/errutil"
	"net/url"
	"os"
)

const defaultEventCount = "3"

// BotConfig contains all the necessary config for running the bot
type BotConfig struct {
	TgBotToken   string
	TargetOrgUrl string
}

// NewBotConfig creates bot config from json
func NewBotConfig() (*BotConfig, error) {
	botToken := os.Getenv("TG_BOT_TOKEN")
	if botToken == "" {
		return nil, errutil2.ErrInvalidConfig.Desc("empty bot token")
	}
	targetOrgUrl := os.Getenv("TG_BOT_TARGET_ORG_URL")
	if targetOrgUrl == "" {
		return nil, errutil2.ErrInvalidConfig.Desc("empty target org url")
	}
	bc := &BotConfig{
		TgBotToken:   botToken,
		TargetOrgUrl: targetOrgUrl,
	}

	u, err := url.Parse(bc.TargetOrgUrl)
	if err != nil {
		return nil, errutil2.ErrInvalidUrl.Desc(bc.TargetOrgUrl)
	}
	// set url param for per_page item
	q := u.Query()
	q.Set("per_page", defaultEventCount)
	u.RawQuery = q.Encode()
	bc.TargetOrgUrl = u.String()
	return bc, nil
}
