## Almost useless Discord bot

__Text commands__

* !roll [max] [times] (!ролл)
* !duel [@user] [bet] (!дуэль)
* !bet  [bet] (!ставка)
* !city [city] (!город) требуется заполнить бд
* !score (!очки)
* !top (!лидеры, !leaderboard)
* !help (!помощь)

__Slash commands__

>_it is necessary to register commands with the -regcommands flag_

>_use the -delcommands flag to remove slash commands after the bot is turned off_

* /roll
* /duel
* /bet
* /city
* /score
* /top

### Installation

Build the bot binary using the make target `build`:
```shell script
make build
```

Run the db migrations ([golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) is needed for this):
```shell script
make migrate-up
```

Once it is built and the migrations applied, then the bot can be run at minimum with the `-t` and the `-db` flags:
```shell script
./bin/bot -t <TOKEN> -db <PATH_TO_DB>
```

To acquire the bot token, follow the [Discord API Docs](https://discord.com/developers/docs/intro).

### Docker
A Docker image containing the bot and migration utility can be built using
```shell script
make docker-image
```

To run the image, pass the token in via the `BOT_TOKEN` environment variable:
```shell script
docker run -e BOT_TOKEN=<TOKEN> bot
```

To make the store persistent, just mount `/bot/data` to an external volume:
```shell script
docker run -e BOT_TOKEN=<TOKEN> -v data:/bot/data bot
```
