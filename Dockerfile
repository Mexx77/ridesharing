FROM golang:1-alpine as builder
RUN addgroup -S ridesharing && adduser -S -G ridesharing ridesharing

WORKDIR /ridesharing
ADD . $WORKDIR

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main server/main.go

# ------------------- Cut Here ------------------ #

FROM scratch
COPY --from=builder /ridesharing/main /
COPY --from=builder /etc/passwd /etc/passwd

USER ridesharing
ENTRYPOINT ["/main"]