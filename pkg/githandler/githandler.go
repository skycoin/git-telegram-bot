package githandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Skycoin/git-telegram-bot/pkg/errutil"
	gh "github.com/google/go-github/v32/github"
	"net/http"
	"time"
)

const (
	defaultRequestTimeout = 5 * time.Second

	tplPushEvent     = "%s pushed new commit(s) to %s\n Link: %s\n"
	tplPREvent       = "%s %s a PR to the repo %s\n Link: %s\n"
	tplPRReviewEvent = "%s reviewed PR in %s\n Link: %s\n"
	tplIssuesEvent   = "%s %s an issue to the repo %s\n Link: %s\n"
	tplCommentEvent  = "%s commented on an issue to the repo %s\n Link: %s\n"
	tplReleaseEvent  = "%s %s a release to the repo %s\n Link: %s\n"
)

// HandleStartCommand will run whenever '/start' command is issued to the bot.
func HandleStartCommand(
	previousEventId string,
	currentEventId string,
	ghUrl string,
	sendFunc func(string) error,
) (string, error) {
	var curEvt gh.Event
	res, err := fetchGhEvent(ghUrl)
	if err != nil {
		return "", err
	}
	// skip if current event is the same as previous event
	if previousEventId == "" && res != nil {
		previousEventId = res[0].GetID()
	}
	if res != nil {
		currentEventId = res[0].GetID()
		curEvt = res[0]
	}
	if currentEventId == previousEventId {
		return previousEventId, nil
	}
	previousEventId = currentEventId

	msgText, err := handleGithubEvent(curEvt)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s sent %s", time.Now().UTC().String(), msgText)
	if err = sendFunc(msgText); err != nil {
		return previousEventId, fmt.Errorf("error sending message: %v", err)
	}
	return previousEventId, nil
}

// handleGithubEvent takes current event and return message in a string format
func handleGithubEvent(curEvt gh.Event) (string, error) {
	payload, err := curEvt.ParsePayload()
	if err != nil {
		return "", errutil.ErrParsePayload.Desc(err)
	}

	var msgText string

	switch curEvt.GetType() {
	case "PushEvent":
		link := bytes.Buffer{}
		evt := payload.(*gh.PushEvent)

		for _, commit := range evt.Commits {
			link.WriteString(fmt.Sprintf(
				"https://github.com/%s/%s/commits/%s",
				curEvt.GetOrg().GetName(),
				curEvt.GetRepo().GetName(),
				commit.GetSHA(),
			))
			link.WriteRune('\n')
		}
		msgText = fmt.Sprintf(
			tplPushEvent,
			curEvt.GetActor().GetLogin(),
			curEvt.GetRepo().GetName(),
			link.String(),
		)
	case "PullRequestEvent":
		evt := payload.(*gh.PullRequestEvent)

		msgText = fmt.Sprintf(
			tplPREvent,
			curEvt.GetActor().GetLogin(),
			evt.GetAction(),
			curEvt.GetRepo().GetName(),
			evt.GetPullRequest().GetHTMLURL(),
		)
	case "PullRequestReviewEvent":
		evt := payload.(*gh.PullRequestReviewEvent)

		msgText = fmt.Sprintf(
			tplPRReviewEvent,
			curEvt.GetActor().GetLogin(),
			curEvt.GetRepo().GetName(),
			evt.GetReview().GetHTMLURL(),
		)
	case "IssuesEvent":
		evt := payload.(*gh.IssuesEvent)

		msgText = fmt.Sprintf(
			tplIssuesEvent,
			curEvt.GetActor().GetLogin(),
			evt.GetAction(),
			curEvt.GetRepo().GetName(),
			evt.GetIssue().GetHTMLURL(),
		)
	case "IssueCommentEvent":
		evt := payload.(*gh.IssueCommentEvent)

		msgText = fmt.Sprintf(
			tplCommentEvent,
			curEvt.GetActor().GetLogin(),
			curEvt.GetRepo().GetName(),
			evt.GetComment().GetHTMLURL(),
		)
	case "ReleaseEvent":
		evt := payload.(*gh.ReleaseEvent)

		msgText = fmt.Sprintf(
			tplReleaseEvent,
			curEvt.GetActor().GetLogin(),
			evt.GetAction(),
			curEvt.GetRepo().GetName(),
			evt.GetRelease().GetHTMLURL(),
		)
	default:
		return "", errutil.ErrUnhandledEvent.Desc(curEvt.GetType())
	}
	return msgText, nil
}

// fetchGhEvent fetches event from an org in github
func fetchGhEvent(ghUrl string) ([]gh.Event, error) {
	httpClient := http.Client{
		Timeout: defaultRequestTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, ghUrl, nil)
	if err != nil {
		return nil, errutil.ErrCreateRequest.Desc(ghUrl, err)
	}

	req.Header.Set("User-Agent", "curl/7.64.1")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errutil.ErrSendingRequest.Desc(ghUrl, err)
	}

	var ghEvt []gh.Event
	if err = json.NewDecoder(res.Body).Decode(&ghEvt); err != nil {
		return nil, errutil.ErrRespBody.Desc(err)
	}
	_ = res.Body.Close()
	return ghEvt, nil
}
