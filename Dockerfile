# Builder 
FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg
COPY internal internal
COPY db db
COPY config config

# RUN go build -o /out/app ./cmd/app

RUN CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w" -o /out/stgorders ./cmd/

EXPOSE 8080

ENV PATH_CONFIG=/src/config/config.yaml

ENTRYPOINT ["/out/app"]


# Runtime
FROM alpine:3.22.0 AS runtime

WORKDIR /src

COPY --from=builder /out/stgorders /src/stgorders

COPY config/config.yaml /src/config/config.yaml
COPY db/migrations /src/db/migrations 

EXPOSE 8080

ENV PATH_CONFIG=/src/config/config.yaml

ENTRYPOINT ["/src/stgorders", "-t"]
