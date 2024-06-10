FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o bitbucket-archiver .

# Copy bin into execution env
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/bitbucket-archiver .

CMD ["/app/bitbucket-archiver"]
