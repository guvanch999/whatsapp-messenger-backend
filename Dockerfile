FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
RUN apt update

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /medium-messenger

ARG BRANCH_NAME=dev

# Run
CMD ["/whatsapp-messenger-backend", "-ldflags='-X main.release=$BRANCH_NAME'"]