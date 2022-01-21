FROM golang:1.17 AS build

WORKDIR /go/src/thirdparty-simulator

COPY go.mod go.sum ./

RUN go mod download

COPY api ./api
COPY cmd ./cmd
COPY main.go ./main.go

RUN go install -v ./

FROM golang:1.17

COPY --from=build /go/bin/thirdparty-simulator /go/bin/thirdparty-simulator
