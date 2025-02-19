# Go Load Balancer

## Features
1. API endpoint to add more servers in case of increase in traffic.
2. Remove unhealthy servers - unhealthy is defined by a `threshold`. If a server is not responding continously for `threshold` amount of time, then this server will automatically be removed.
3. Retry mechanism - if a server fails to respond to a request sent by the load balancer, then the load balancer tries to send the same request to another server.
4. Round robin strategy implementation
5. Rate limiting the load balancer *(TODO)*


## Upcoming...

1. Functional and load testing
2. Nice and simple UI in React and TailwindCSS
