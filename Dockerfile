FROM golang:1.17 AS base

WORKDIR /go/src/tctl

COPY deployment/tctl/go.mod deployment/tctl/go.sum ./

RUN go mod download

WORKDIR /go/src/background-checks

COPY go.mod go.sum ./

RUN go mod download

FROM base AS build

WORKDIR /go/src/tctl

COPY deployment/tctl ./

RUN make

WORKDIR /go/src/background-checks

COPY activities ./activities
COPY api ./api
COPY cli ./cli
COPY temporal ./temporal
COPY utils ./utils
COPY ui ./ui
COPY workflows ./workflows

RUN go install -v ./cli/bgc-backend
RUN go install -v ./cli/bgc-company
RUN go install -v ./cli/bgc-candidate
RUN go install -v ./cli/bgc-researcher
RUN go install -v ./temporal/dataconverter-plugin
RUN go install -v ./temporal/dataconverter-server

FROM golang:1.17 AS app

ENV TEMPORAL_CLI_PLUGIN_DATA_CONVERTER=dataconverter-plugin

COPY --from=build /go/src/tctl/tctl /usr/local/bin/tctl
COPY --from=build /go/bin/dataconverter-plugin /usr/local/bin/dataconverter-plugin
COPY --from=build /go/bin/dataconverter-server /usr/local/bin/dataconverter-server

COPY --from=build /go/bin/bgc-backend /usr/local/bin/bgc-backend
COPY --from=build /go/bin/bgc-company /usr/local/bin/bgc-company
COPY --from=build /go/bin/bgc-candidate /usr/local/bin/bgc-candidate
COPY --from=build /go/bin/bgc-researcher /usr/local/bin/bgc-researcher