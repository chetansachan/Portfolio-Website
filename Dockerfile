# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /app

# copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 go build -o out ./cmd/main.go


# ---- runtime stage ----
FROM alpine:3.20
RUN adduser -D appuser
WORKDIR /app
COPY --from=build /app/out .
COPY --from=build /app/migration ./migration
RUN mkdir -p data/uploads && chown -R appuser:appuser /app
ENV LOCAL_SYSTEM_DB=false
USER appuser
EXPOSE 8080
CMD ["./out"]