package config

import (
	"github.com/Skycoin/git-telegram-bot/pkg/errutil"
	"net/url"
	"os"
	"strconv"
)

const defaultEventCount = "3"

// BotConfig contains all the necessary config for running the bot
type BotConfig struct {
	TgBotToken        string
	TargetOrgUrl      string
	TargetGroupChatId int64
}

// NewBotConfig creates bot config from json
func NewBotConfig() (*BotConfig, error) {
	botToken := os.Getenv("TG_BOT_TOKEN")
	if botToken == "" {
		return nil, errutil.ErrInvalidConfig.Desc("empty bot token")
	}
	targetOrgUrl := os.Getenv("TG_BOT_TARGET_ORG_URL")
	if targetOrgUrl == "" {
		return nil, errutil.ErrInvalidConfig.Desc("empty target org url")
	}
	targetGroupChatId := os.Getenv("TG_BOT_GROUP_CHAT_ID")
	if targetGroupChatId == "" {
		return nil, errutil.ErrInvalidConfig.Desc("empty group chat id")
	}
	groupChat, err := strconv.Atoi(targetGroupChatId)
	if err != nil {
		return nil, errutil.ErrInvalidConfig.Desc("invalid group chat ID, expected int")
	}
	bc := &BotConfig{
		TgBotToken:        botToken,
		TargetOrgUrl:      targetOrgUrl,
		TargetGroupChatId: int64(groupChat),
	}

	u, err := url.Parse(bc.TargetOrgUrl)
	if err != nil {
		return nil, errutil.ErrInvalidUrl.Desc(bc.TargetOrgUrl)
	}
	// set url param for per_page item
	q := u.Query()
	q.Set("per_page", defaultEventCount)
	u.RawQuery = q.Encode()
	bc.TargetOrgUrl = u.String()
	return bc, nil
}
