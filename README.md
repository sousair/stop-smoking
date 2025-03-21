Just a Telegram bot to help me stop smoking by a strict schedule.

## Requirements
- Docker (with compose plugin)

## Setup

```bash
cp .env.example .env
```

You need to get a Telegram bot token from the [BotFather](https://t.me/botfather) and set it in the `.env` file.

```bash
mkdir db && touch db/db.sqlite
```
You can change the database path in the `.env` file.

## Run

```bash
make run
```

Or, to run in detached mode:

```bash
make rund
```
