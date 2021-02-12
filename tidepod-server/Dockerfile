FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=darwin \
  GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .

# Build the application
RUN go mod download
# RUN go build -o main main.go

# # Export necessary port
# EXPOSE 3000

# CMD ["/build/main"]