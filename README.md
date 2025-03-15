# Go Load Balancer

I like building things from scratch. Lately I was reading about load balancers and thought of building a little prototype to understand the inner workings better.


## How to build & run this project?
```
Commands
```

## How to perform load testing of this project?

```
Commands
```

## Why this might be helpful for you?

- You can use this as a base project for your implementation of a load balancer.
- By playing around with the load balancer tool, you can see what would happen by tweaking different settings of a load balancer.
- By playing around with the load balancer tool, you would understand the load balancer working much better.


## Features
1. API endpoint to add more servers in case of increase in traffic.
2. Remove unhealthy servers - unhealthy is defined by a `threshold`. If a server is not responding continously for `threshold` amount of time, then this server will automatically be removed.
3. Retry mechanism - if a server fails to respond to a request sent by the load balancer, then the load balancer tries to send the same request to another server.
4. Round robin strategy implementation
5. Rate limiting the load balancer *(TODO)*
6. Implement popular load balancing strategies.
    - Random strategy
    - Round robin strategy
    - *(TODO)* Add more strategies


## Upcoming...

1. Functional and load testing
2. Nice and simple UI in React and TailwindCSS
3. How should we demo this app? *(TODO)*
4. I wanna load test this app.
    -- First we need to think about what we need to test, when I say load testing.
    -- Then we need to think about how to implement it in Artillery.
    -- I think we also need to read a little bit about how artillery works and in what possible ways we can manipulate the artillery yaml file.



We need to test the following - 

1. Basic load testing
    - max number of requests that the load balancer can handle without errors or significant performance degradation
    - Monitor response times and latency and error rates as the load increases

2. Burst traffic and sustained high load
    - when we finally have max load we must sustain it for some time  and see how the system performs

3. Failover and redundancy
    - fail one or more servers and see how the system performs
    - in particular how rerouting of the requests are done.
    - in this case, do you see any difference in response times and latency.

4. Compare various traffic distribution algorithms - 
    - In our case, we have implement round robin, random and least connections strategy

5. How are the resources being utilised?
    - By resources, I mean CPU, Memory, Disk and Network usage.


