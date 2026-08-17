# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /app

# copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/main.go


# ---- runtime stage ----
FROM alpine:3.20
RUN adduser -D appuser
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/migration ./migration
USER appuser
EXPOSE 8080
CMD ["./server"]