FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod download
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o qnap-lcd-display-manager .

FROM scratch AS binaries
COPY --from=builder /app/qnap-lcd-display-manager /