FROM golang:1.22 AS build

WORKDIR /app

COPY . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -a -o backend .

FROM gcr.io/distroless/static AS release

WORKDIR /app

COPY --from=build /app/backend /app/backend

EXPOSE 8080
    
CMD ["/app/backend"]