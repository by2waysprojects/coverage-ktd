FROM --platform=$BUILDPLATFORM golang:1.22.4 AS builder

ARG GH_PAT

WORKDIR /app

COPY go.mod go.sum ./

RUN git config --global url."https://${GH_PAT}@github.com/".insteadOf "https://github.com/"

RUN go mod download

COPY . .

ARG TARGETARCH TARGETOS

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o coverage-ktd ./cmd/security-tester/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/coverage-ktd /usr/local/bin/coverage-ktd

RUN mkdir -p /attacks
COPY --from=builder /app/attacks /attacks

RUN mkdir -p /reports

CMD ["/usr/local/bin/coverage-ktd"]