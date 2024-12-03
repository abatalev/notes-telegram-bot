# Telegram Bot for Create Notes from Messages 

## Build and run

```sh
./mk.sh
export TELEGRAM_APITOKEN=....
export TELEGRAM_FROMID=....
./notes-telegram
```

## Use docker compose

```sh
docker build -t notes-telegram-bot .
cp docker-compose.tmpl docker-compose.yaml
echo "TELEGRAM_APITOKEN=..." > .env
docker compose up -d
```

## Environment variables

TELEGRAM_APITOKEN
: Token from BotFather

TELEGRAM_FROMID
: telegram user ID 

TELEGRAM_DIR
: directory for save notes

TELEGRAM_PREFIX
: prefix saved file ({{prefix}}{{NNN}}.md) 

## Usage pattern

TBD

## see also

- https://github.com/soberhacker/obsidian-telegram-sync
- https://t.me/smartspeech_sber_bot
