FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Move to working directory /build
WORKDIR /build
RUN go get github.com/Kagami/go-face
RUN go get github.com/codingsince1985/geo-golang
RUN go get github.com/disintegration/imaging

RUN go get github.com/gin-gonic/contrib

RUN go get github.com/gin-gonic/gin

RUN go get github.com/h2non/filetype

RUN go get github.com/rwcarlsen/goexif

RUN go get golang.org/x/image

RUN go get gorm.io/driver/sqlite

RUN go get gorm.io/gorm

# Copy the code into the container
COPY . .

# Export necessary port
EXPOSE 3000