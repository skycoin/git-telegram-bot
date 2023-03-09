package githandler

import (
	"encoding/json"
	config2 "github.com/Skycoin/git-telegram-bot/internal/config"
	errutil2 "github.com/Skycoin/git-telegram-bot/pkg/errutil"
	gh "github.com/google/go-github/v32/github"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigAndFetch(t *testing.T) {
	r := require.New(t)
	cases := map[string]*config2.BotConfig{
		"notValidUrl": {
			TgBotToken:   "test",
			TargetOrgUrl: "5678",
		},
		"validUrl": {
			TgBotToken:   "test",
			TargetOrgUrl: "https://api.github.com/orgs/skycoin/events",
		},
	}
	expected := map[string]errutil2.BotErr{
		"notValidUrl": errutil2.ErrSendingRequest,
	}
	for k, v := range cases {
		res, err := fetchGhEvent(v.TargetOrgUrl)
		if strings.Contains(k, "valid") {
			r.NoError(err)
			r.NotNil(res)
			continue
		}
		r.IsType(expected[k], err)
	}
}

func TestGitHandler(t *testing.T) {
	r := require.New(t)
	var res []gh.Event
	data, err := ioutil.ReadFile(filepath.Join("testdata", "example.json"))
	r.NoError(err)
	err = json.Unmarshal(data, &res)
	r.NoError(err)

	cases := map[string]gh.Event{
		"push":              res[5],
		"pullRequest":       res[1],
		"pullRequestReview": res[46],
		"issue":             res[2],
		"issueComment":      res[0],
	}


	successes := map[string]string{
		"push":              "SkycoinSynth pushed new commit(s) to skycoin/cx-game\n Link: https://github.com//skycoin/cx-game/commits/f5aefd94d8c744807ab1b9e49773048e2c994ff6\nhttps://github.com//skycoin/cx-game/commits/299a2378642e0f8d91d878cfcb00dad5c5eb2db7\nhttps://github.com//skycoin/cx-game/commits/1e123440b087d1a4a0d0c458d6cb8ec227501649\n\n",
		"pullRequest":       "ted537 opened a PR to the repo skycoin/cx-game\n Link: https://github.com/skycoin/cx-game/pull/73\n",
		"pullRequestReview": "jdknives reviewed PR in skycoin/skywire\n Link: https://github.com/skycoin/skywire/pull/744#pullrequestreview-645902433\n",
		"issue":             "jdknives opened an issue to the repo skycoin/skycoin-services\n Link: https://github.com/skycoin/skycoin-services/issues/17\n",
		"issueComment":      "CLAassistant commented on an issue to the repo skycoin/cx-game\n Link: https://github.com/skycoin/cx-game/pull/73#issuecomment-828603924\n",
	}

	for k, v := range cases {
		msg, err := handleGithubEvent(v)
		r.NoError(err)
		r.Equal(successes[k], msg)
	}
}
