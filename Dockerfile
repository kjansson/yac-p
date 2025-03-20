FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM public.ecr.aws/lambda/provided:al2
COPY --from=builder /app/main /main
ENTRYPOINT [ "/main" ]