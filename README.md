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
   export TG_BOT_GROUP_CHAT_ID="YOUR TELEGRAM GROUP ID"
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
docker run --rm -e TG_BOT_TOKEN=<YOUR BOT CREDENTIAL> -e TG_BOT_TARGET_ORG_URL=https://api.github.com/orgs/skycoin/events -e TG_BOT_GROUP_CHAT_ID=<YOUR GROUP CHAT ID> -it git-telegram-bot:latest
```

### Bot Instruction

To operate on the bot, you have to be an administrator to the group you configured to.

To actually run the bot, type:

```
/startpoll@NAME_OF_YOUR_BOT
```

To stop the bot, type:

```
/stoppoll@NAME_OF_YOUR_BOT
```

To restart the bot, type:

```bash
/resetpoll@NAME_OF_YOUR_BOT
# then
/startpoll@NAME_OF_YOUR_BOT
```

To view its help menu, type:

```
/helppoll@NAME_OF_YOUR_BOT
```