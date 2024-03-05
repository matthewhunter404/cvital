FROM golang:1.16-alpine as builder
WORKDIR /project
COPY server/* .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main main.go

FROM alpine:latest as cvital

WORKDIR /cvital
COPY --from=builder /project/main .
COPY --from=builder /project/config.yml .
RUN chmod +x ./main
CMD ["/cvital/main"]
EXPOSE 3000