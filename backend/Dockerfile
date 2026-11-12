FROM golang:1.27-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/pizzahub-api .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/pizzahub-api ./pizzahub-api
USER app
EXPOSE 8080
ENTRYPOINT ["./pizzahub-api"]
