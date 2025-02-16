
2. monitor metrics of the load balancer - 
    - note down the metrics
    - implement them

5. I need to package the project in such a way so that it's easy to build, run and test it.
6. Implement test setup.
7. Complete the load balancer and push to github. Forget about the UI part of it right now.
8. Then complete the image server part. I would say do not focus on the UI part as of now.

- Core requirements 
    1. has an interface to
        - ~~add and remove backend servers~~
        - ~~see which of the configured backend servers are healthy~~
        - visualize load balancer metrics
        - change load balancing strategy on the fly
        - changes should not require a reboot to take effect

## Important commands
1. *Start a HTTP server at port 8000*
```
python -m http.server 8000
```

2. *Send a request to a server @ `localhost:8000`*
```
curl -s -i http://localhost:8000
```
