FROM golang:1.13-alpine3.10 as builder
RUN addgroup -S ridesharing && adduser -S -G ridesharing ridesharing
RUN apk add build-base

ADD server /ridesharing/server
WORKDIR /ridesharing/server
ENV GO111MODULE=on
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
 -ldflags "-linkmode external -extldflags -static" \
 -a -installsuffix cgo -o main main.go

# ------------------- Cut Here ------------------ #

FROM scratch
ADD dist /dist
WORKDIR /server
ADD server/sqlite.db $WORKDIR
COPY --from=builder /ridesharing/server/main $WORKDIR
COPY --from=builder /etc/passwd /etc/passwd

USER ridesharing
ENTRYPOINT ["/server/main"]