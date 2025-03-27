import json
import subprocess
import time
import os


def start_server(server):
    port = server.split(":")[-1]

    command = f"python -m http.server {port}"
    log_path = f"./logs/server_{port}.log"

    try:
        os.makedirs(os.path.dirname(log_path), exist_ok=True)
        # check if the server is already running
        # start the server in a new subprocess
        with open(log_path, "w") as log:
            process = subprocess.Popen(command, shell=True, stdout=log, stderr=log)
            print (f"Server started on port {port} with PID {process.pid}, logging to {log_path}")
            return process
    except Exception as e:
        print (f"Failed to start server on port {port}: {e}")
    
def shutdown_servers(processes):
    for process in processes:
        print (f"Shutting down server with PID {process.pid}")
        process.terminate()
        process.wait()
    
    print (f"All servers have been shut down.")

if __name__ == "__main__":
    # load the list of servers from the load balancer config json
    with open("config.json") as f:
        config = json.load(f)

    servers = config['servers']
    processes = []

    for server in servers:
        print (f"Starting server@: {server}")
        process = start_server(server)
        if process:
            processes.append(process)
    
    # run this script continuously
    try:
        while True:
            time.sleep(2)
    except KeyboardInterrupt:
        print (f"Stopping the servers...")
        shutdown_servers(processes)

