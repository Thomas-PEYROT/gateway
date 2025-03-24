# Gateway

This project is the API gateway of [Go microservices project](https://github.com/Thomas-PEYROT/go-microservices-architecture).
Note that if a microservice is launched multiple times, the method used for load balancing is **random load balancing**.

## Setup

First of all, clone this repository :

```
git clone https://github.com/Thomas-PEYROT/gateway
```

After that, copy the example `.env` file :

```
cp .env.example .env
```

Finally, launch the server :

```
go run .
```
