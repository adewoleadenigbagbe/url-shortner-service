# Use an official Golang runtime as a parent image
FROM golang:1:20.5

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Installs Go dependencies
RUN go mod download

# Build the Go application inside the container
RUN CGO_ENABLED=0 go build -o urlshortner

# Tells Docker which network port your container listens on
EXPOSE 8321

# Define the command to run your application
ENTRYPOINT ["./urlshortner"]