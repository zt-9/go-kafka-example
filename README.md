# go-kafka-example
Data processing with Apache Kafka, REST API and Redis.

We process blockchain on-chain data with Apache Kafka. 
There are three topics: "address", "label" and "transaction". We have one producer and one consumer listen to and process these topics.
We use Redis as in memory database to store data.

The producer runs on port 8080 and the consumer runs on port 8081. 
Users send REST API requests to producer. Producer then send messages to consumer.

## Prerequistes
- docker ([install](https://docs.docker.com/engine/install/))

## Run the application
1. start docker img: 
    ```bash
    make run/docker
    ```
2. start the producer in a new terminal
    ```bash
    make run/producer
    ```
3. start the consumer in a new terminal
    ```bash
    make run/consumer
    ```
4. send api requests to localhost:8080
    examples:
    create an address:

    send POST request to `localhost:8080/addresses/0x87631B45877794f9cdd50a70c827403e3C36d072`
    with body
    ```json
    {
    "address": "0x87631B45877794f9cdd50a70c827403e3C36d072",
    "labels": ["eoa"]
    
    }
    ```

    get an address:
    
    send GET request to `localhost:8081/addresses/0x87631B45877794f9cdd50a70c827403e3C36d072`   

## Project structure
- `/cmd`
main applications for this project
- `/cmd/producer`
kafka producer. it sends messages to kafka.
- `/cmd/consumer`
kafka consumer. It process the received 
===
- `config`
configuration
===
- `/pkg/controllers`
controllers that handles api endpoints request
- `/pkg/models`
data structures and models used in REST API
- `/utils`
util/helper functions