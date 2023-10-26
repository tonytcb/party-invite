########################################## Build Stage ##########################################

FROM golang:1.21-bullseye as build

RUN mkdir /app
ADD . /app
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o party-invite ./cmd/main.go

########################################## Deploy Stage ##########################################

FROM alpine:latest

COPY --from=build /app .
COPY app.env app.env

CMD ["./party-invite"]
