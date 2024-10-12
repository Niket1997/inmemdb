# inmemdb
A simple redis-compatible asynchronous in-memory KV store.

### How to run this application in linux environment on mac
1. Install Docker for Desktop
2. Add following function to your `.zshrc` file, ensure to update the project directory
```bash
startinmemdblinux() {
    # Navigate to the Docker Compose project directory
    PROJECT_DIR="/Users/aniket.mahangare/myProjects/inmemdb"
    
    cd "$PROJECT_DIR/ubuntu-linux-docker" || {
        echo "❌ Failed to navigate to project directory."
        return 1
    }

    # Start Docker Compose services in detached mode
    docker-compose up -d || {
        echo "❌ Docker Compose failed to start."
        return 1
    }

    # Define the cleanup function to stop Docker containers
    cleanup() {
        echo "🛑 Stopping Docker containers..."
        docker-compose down
    }

    # Trap the SIGINT signal (Control + C) to execute the cleanup function
    trap cleanup SIGINT

    # Execute the Go application inside the Docker container
    docker exec -it ubuntu_container bash -c "cd inmemdb && go run main.go" || {
        echo "❌ Failed to execute 'go run main.go' inside the container."
        cleanup
    }

    # If the Go application exits normally, perform cleanup
    cleanup
}
```
3. restart your bash shell
4. execute the function 
```bash
startinmemdblinux
```
