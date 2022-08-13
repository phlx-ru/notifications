# Notifications Service

Send few types of messages (email, push, telegram) to recipients.

Based on [Kratos](https://github.com/go-kratos/kratos) Golang framework.

## TODO

    // TODO Make this README better

## Prerequisites

- [required] PostgreSQL Server
- [optional] Graphite backend (for metrics)

All prerequisites resolve with running Docker Swarm environment, created 
by repo [infra](https://github.com/phlx-ru/infra).

## How to start

Make file `./configs/.env` with correct variables:

```
LOG_LEVEL=warn

METRICS_ADDRESS=graphite:8125
METRICS_MUTE=false

DATA_DATABASE_MIGRATE=hard
DATA_DATABASE_DEBUG=false

AUTH_JWT_SECRET=doedoesdeezdodge

POSTGRES_USER=postgres
POSTGRES_PASS=postgres
POSTGRES_DB=notifications
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

SENDERS_PLAIN_FILE=./plain_messages.log
SENDERS_EMAIL_FROM="John Doe <johndoe@mail.example>"
SENDERS_EMAIL_ADDRESS=smtp.mail.example:587
SENDERS_EMAIL_USERNAME=johndoe@mail.example
SENDERS_EMAIL_PASSWORD=ilovejanedoe
```

Run service
```
make run-server
```

Test with plain notification:
```
curl --location --request POST 'https://localhost:8000/send' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data-raw '{
  "type": 0,
  "payload": {
    "message": "Hello from the outside!"
  }
}'
```

## Any help?

Try
```
make help
```

and [Kratos Docs](https://go-kratos.dev/en/docs/)
