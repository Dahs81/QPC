# QPC
A simple golang example for a producer and consumer using rabbitMQ

### Usage

```
go get ./...
```

In one terminal, cd to the producer directory and run the producer with optional QPC_RABBITMQ_URL environment variable

```
QPC_RABBITMQ_URL=amqp://guest:guest@localhost:5672/ go run main.go
```

In another terminal, cd to the consumer directory and run the consumer with optional QPC_RABBITMQ_URL environment variable

```
QPC_RABBITMQ_URL=amqp://guest:guest@localhost:5672/ go run main.go
```

**Note:** If `QPC_RABBITMQ_URL` is not set, it will default to `amqp://guest:guest@localhost:5672/`

curl the API endpoint in another terminal (or Postman):

```
curl -X POST -H 'Content-Type:application/json' http://localhost:8080/send -d '{"name": "Testing this out"}';
```

OR

```
for (( c=1; c<=50; c++ )); do
    curl -X POST -H 'Content-Type:application/json' http://localhost:8080/send -d '{"name": "Testing this out"}';
done
```

### API
#### Endpoint
POST /send

##### Request Body

```
{
    "name": "<some string value>"
}
```