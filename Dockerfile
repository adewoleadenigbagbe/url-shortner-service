# Use an official Golang runtime as a parent image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY ./backend .


# Installs Go dependencies
RUN go mod tidy && go mod download

# Build the Go application inside the container
RUN CGO_ENABLED=1 go build -o urlshortner -a -ldflags '-linkmode external -extldflags "-static"'

# Tells Docker which network port your container listens on
EXPOSE 8080

# Define the command to run your application
# ENTRYPOINT ["./urlshortner","serveapi"]