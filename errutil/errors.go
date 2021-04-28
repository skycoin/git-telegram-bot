package errutil

import "fmt"

// BotErr returns formatted error for skygit-bot-tg errors.
type BotErr string

const (
	ErrInvalidUrl        BotErr = "error parsing target org url %s: invalid url"
	ErrCreateRequest     BotErr = "error creating request to %s: %v"
	ErrSendingRequest    BotErr = "error sending request to %s: %v"
	ErrRespBody          BotErr = "error parsing github event: %v"
	ErrNonExistentConfig BotErr = "error non existent config file at specified path %s"
	ErrInvalidConfig     BotErr = "error parsing config file: %v"
	ErrCreatingBot       BotErr = "error creating telegram bot: %v"
	ErrUnhandledEvent    BotErr = "error unhandled event: %v"
	ErrParsePayload      BotErr = "error parsing payload: %v"
)

func (b BotErr) Desc(args ...interface{}) BotErr {
	return BotErr(fmt.Sprintf(b.String(), args...))
}

func (b BotErr) String() string {
	return string(b)
}

func (b BotErr) Error() string {
	return string(b)
}
