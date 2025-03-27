# Go Load Balancer Go!!

*I genuinely believe that in order to understand a concept, you must try building it from scratch.*

*Here is a simple Golang implementation of a load balancer that was given as a homework exercise in [*The System Design Masterclass by Arpit Bhayani*](https://arpitbhayani.me/masterclass/).*


## Install dependencies

```sh
TODO
```

## *[Optional]* Update config.json
You can customize the load balancer's behavior by modifying config.json. The following options are available:

`port`: Specifies the port on which the load balancer runs locally. You can change it as needed.

`healthCheckInterval`: Defines the interval (in milliseconds) at which the load balancer continuously polls all registered servers to check their health.

`lbAlgo`: Determines the load-balancing strategy. Available options:<br>
-  `round-robin`
- `random`
- `least-connections`

`servers`: Lists the backend servers that the load balancer can route traffic to.

`maxRetries`: Specifies the number of retry attempts when a server fails to process a request. The request is forwarded to the next available server based on the selected load-balancing strategy.

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

## How to Perform Load Testing

You can use [Artillery](https://www.artillery.io/) to simulate high traffic and measure the performance of your load balancer.

### Prerequisites  
Ensure you have Artillery installed. If not, install it globally using:  
```sh
npm install -g artillery@latest
```

### Running the Load Test
Execute the following command to run a load test using `artillery.yaml`:

```sh
artillery run artillery.yaml --record --key ARTILLERY_API_KEY
```

## Why This Might Be Helpful for You

- Use this as a base project to implement your own load balancer.
- Experiment with different load balancer settings to observe their impact.
- Gain a deeper understanding of how a load balancer works by interacting with its configurations and mechanisms.

## Features

1. **Dynamic Server Scaling**  
   - Add new servers dynamically via an API endpoint to handle increased traffic.

2. **Automatic Removal of Unhealthy Servers**  
   - Servers are monitored continuously.  
   - If a server fails to respond for a specified `threshold` duration, it is automatically removed from the pool.

3. **Intelligent Retry Mechanism**  
   - If a server fails to process a request, the load balancer retries the request with another available server.

4. **Multiple Load Balancing Strategies**  
   - Supports the following strategies:  
     - **Round Robin**  
     - **Random**  
     - **Least Connections**


