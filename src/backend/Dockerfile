# use official Golang image
# FROM golang:1.24.2-alpine3.21
FROM golang:1.24.2-alpine3.21 AS buildstage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8000
RUN go build -o /flab ./cmd/server/main.go

FROM alpine:latest
WORKDIR /
COPY --from=buildstage /flab /flab
EXPOSE 8000
ENTRYPOINT [ "/flab" ]


# # set working directory
# WORKDIR /app

# # Copy the source code
# COPY . .

# # Download and install the dependencies
# RUN go get -d -v ./...

# # Build the Go app
# RUN go build -o api .

# #EXPOSE the port
# EXPOSE 8000

# # Run the executable
# CMD ["./api"]