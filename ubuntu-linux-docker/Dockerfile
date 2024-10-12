# Use the official Ubuntu base image
FROM ubuntu:latest

# Set the environment variable for DEBIAN_FRONTEND to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages and dependencies
RUN apt-get update && \
    apt-get install -y wget tar git build-essential && \
    rm -rf /var/lib/apt/lists/*

# Define Go version
ENV GO_VERSION=1.23.2

# Download and install Go
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-arm64.tar.gz && \
    rm go${GO_VERSION}.linux-arm64.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin
ENV MYPROJECTSPATH=/home/docker_user/myProjects

# Create Go workspace directory
RUN mkdir -p $GOPATH

# Set working directory
WORKDIR $MYPROJECTSPATH

# Copy your application code (if any)
# COPY . $GOPATH

# Expose necessary ports (if any)
# EXPOSE 8080

# Define the command to run when starting the container
CMD ["/bin/bash"]
