FROM golang:1.13-alpine3.10 as builder
RUN addgroup -S ridesharing && adduser -S -G ridesharing ridesharing
RUN apk add --update npm

# Build frontend
WORKDIR /ridesharing
ADD src ./src
ADD public ./public
ADD package.json .
ADD babel.config.js .
RUN npm install && npm run build

# Build backend
ADD server ./server
WORKDIR server
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
 -a -installsuffix cgo -o main main.go

# ------------------- Cut Here ------------------ #

FROM scratch
WORKDIR /server
COPY --from=builder /ridesharing/dist/ /dist
COPY --from=builder /ridesharing/server/main .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER ridesharing
ENTRYPOINT ["/server/main"]