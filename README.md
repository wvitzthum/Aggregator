## Aggregator

This repo is a test for event sourcing using Goka and Kafka.

The service is consuming a websocket stream of bitcoin transactions, writing them to a Kafka topic and then aggregates them into a window based on the source address.

Its a merged and refined version [Mike's](https://github.com/mikedewar) repos:
- [BTC Dispersion](https://github.com/mikedewar/btcDispersion)
- [Aggregator](https://github.com/mikedewar/aggregator)



## Architecture
![Flow Diagram](docs/flowDiagram.jpg)



## Observability
Messages and statistics from the topics can be observed using the web interface (docker-compose) on http://localhost:8080.

Data can be retrieve thorough the rest interfaces on:
- http://localhost:9095/window/{key}
- http://localhost:9095/window/

## General Observations
It appears that the default log compaction settings on Goka are not sufficient for the rate of data and amount of partitions.


## Currently under development