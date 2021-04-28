# git-telegram-bot

Telegram bot for notifying org events

### Requirements (for building)

- `Go` version `1.16.x`

### Setup

- If you don't have a telegram bot token yet, create one via [@botfather](https://t.me/botfather)
- Set Environment Variables:

```
   export TG_BOT_TOKEN="YOUR BOT TOKEN CREDENTIAL"
   export TG_BOT_TARGET_ORG_URL="YOUR ORG URL example: https://api.github.com/orgs/skycoin/events"
```

### Build

1. Clone this repo `git clone https://github.com/skycoin/git-telegram-bot`
2. Run `make build`
3. Your binary should be in the root-directory of the cloned repository.

or alternatively via Docker:

1. Clone this repo `git clone https://github.com/skycoin/git-telegram-bot`
2. Run `make docker`

### Run

Make sure you already set the required env variables above.

Run it via

- Baremetal:

```
./git-telegram-bot
```

- Docker:

```
docker run --rm -e TG_BOT_TOKEN=<YOUR BOT CREDENTIAL> -e TG_BOT_TARGET_ORG_URL=https://api.github.com/orgs/skycoin/events -it git-telegram-bot:latest
```