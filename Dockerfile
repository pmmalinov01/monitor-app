FROM golang:alpine as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o app .

FROM alpine:3.14 as worker
WORKDIR /
COPY --from=builder /build/app ./

# Export necessary port
EXPOSE 8001

# Command to run when starting the container
CMD ["/app"]
