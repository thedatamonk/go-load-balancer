**Functionality to test**

- Add new servers to the pool.
- Automatic removal of unhealthy servers from the pool.
- Ability to switch load balancing strategy without server restarts.
- visualise load balancing metrics
    - for this take help from lovable.dev and build a sweet little UI.

## Important commands
1. *Start a HTTP server at port 8000*
```
python -m http.server 8000
```

2. *Send a request to a server @ `localhost:8000`*
```
curl -s -i http://localhost:8000
```
