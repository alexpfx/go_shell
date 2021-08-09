FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/go_shell

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /bin/go_shell /go_shell
USER nonroot:nonroot

ENTRYPOINT ["/go_shell"]
