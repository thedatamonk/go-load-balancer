# Go Load Balancer Go!!

*I genuinely believe that in order to understand a concept, you must try building it from scratch.*

*Here is a simple Golang implementation of a load balancer that was given as a homework exercise in [*The System Design Masterclass by Arpit Bhayani*](https://arpitbhayani.me/masterclass/).*


## Install dependencies

```sh
TODO
```

## Run the load balancer

```sh
# 1. Clone this repo
git clone https://github.com/thedatamonk/go-load-balancer.git

# 2. Go to the main directory
cd go-load-balancer

# 3. Start the HTTP servers
# This will also use config.json as input to get the name of the servers that need to be managed by the load balancer
python manage_servers.py

# 4. Start the load balancer
# The load balancer will be started on the localhost at the port number specified in the config.json
go run .

# 5. Send a request to the load balancer
# The grep command is used for extracting the value of "X-Forwarded-Server" from the server response.
curl -s -i http://localhost:8080 | grep "X-Forwarded-Server"

```

## How to perform load testing of this project?

```
artillery run artillery.yaml --record --key <API_KEY>
```

## Why this might be helpful for you?

- You can use this as a base project for your implementation of a load balancer.
- By playing around with the load balancer tool, you can see what would happen by tweaking different settings of a load balancer.
- By playing around with the load balancer tool, you would understand the load balancer working much better.


## Features
1. API endpoint to add more servers in case of increase in traffic.
2. Remove unhealthy servers - unhealthy is defined by a `threshold`. If a server is not responding continously for `threshold` amount of time, then this server will automatically be removed.
3. Retry mechanism - if a server fails to respond to a request sent by the load balancer, then the load balancer tries to send the same request to another server.
4. Common load balancer strategies are implemented - 
    - ***Round Robin***
    - ***Random***
    - ***Least connections***

