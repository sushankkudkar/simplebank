# Build stage
FROM golang:1.22-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM golang:1.22-bullseye
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
# COPY start.sh .
# COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]