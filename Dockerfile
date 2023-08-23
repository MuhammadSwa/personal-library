# ############################
# # STEP 1 build executable binary
# ############################
FROM golang:alpine AS build-stage
# # Install git.
# # Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git 

WORKDIR /app
COPY . .
COPY .env .

# # Fetch dependencies.
# # Using go get.
# RUN go get -d -v
RUN go mod download


# RUN go get github.com/githubnemo/CompileDaemon
# RUN go get -v golang.org/x/tools/gopls
RUN go install github.com/mitranim/gow@latest

EXPOSE 8080

ENTRYPOINT gow -c -e=go,mod,tmpl run ./cmd/main.go
# # Build the binary.
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/personal-lib ./cmd/web
# CMD ["/bin/personal-lib"]
#
# # ############################
# # # STEP 2 build a small image
# # ############################
# FROM scratch
# WORKDIR /
# # # Copy our static executable.
# COPY --from=build-stage /app/bin/personal-lib /bin/personal-lib
# COPY --from=build-stage /app/.env ./
#
# EXPOSE 8080
# # USER nonroot:nonroot
#
# CMD ["/bin/personal-lib"]
#
# # FROM golang:alpine AS builder
# # # Install git.
# # # Git is required for fetching the dependencies.
# # RUN apk update && apk add --no-cache git
# # WORKDIR $GOPATH/src/mypackage/myapp/
# # COPY . .
# # # Fetch dependencies.
# # # Using go get.
# # RUN go get -d -v
# # # Build the binary.
# # # RUN go build -o /go/bin/hello
# # RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o go/bin/personal-lib ./src/main.go
# # ############################
# # # STEP 2 build a small image
# # ############################
# # FROM scratch
# # # Copy our static executable.
# # COPY --from=builder /go/bin/hello /go/bin/hello
# # # Run the hello binary.
# # ENTRYPOINT ["/go/bin/hello"]
