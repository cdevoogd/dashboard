FROM golang:1.21-alpine AS builder
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod go.sum /src/
RUN go mod download -x
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build  \
    go build -o /src/dashboard .

FROM scratch
COPY --from=builder /src/dashboard /
ENTRYPOINT ["/dashboard"]
