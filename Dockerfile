FROM golang:1.17 AS base

WORKDIR /go/src/background-checks

COPY go.mod go.sum ./

RUN go mod download

FROM base AS build-app

COPY activities ./activities
COPY api ./api
COPY cli ./cli
COPY temporal ./temporal
COPY types ./types
COPY utils ./utils
COPY ui ./ui
COPY workflows ./workflows

RUN go install -v ./cli/bgc-backend
RUN go install -v ./cli/bgc-company
RUN go install -v ./cli/bgc-candidate
RUN go install -v ./cli/bgc-researcher

FROM golang:1.17 AS app

COPY --from=build-app /go/bin/bgc-backend /usr/local/bin/bgc-backend
COPY --from=build-app /go/bin/bgc-company /usr/local/bin/bgc-company
COPY --from=build-app /go/bin/bgc-candidate /usr/local/bin/bgc-candidate
COPY --from=build-app /go/bin/bgc-researcher /usr/local/bin/bgc-researcher

FROM base AS build-plugin

COPY temporal/dataconverter ./temporal/dataconverter
COPY temporal/dataconverter-plugin ./temporal/dataconverter-plugin

RUN go install -v ./temporal/dataconverter-plugin

FROM golang:1.17 AS tools

ENV TEMPORAL_CLI_PLUGIN_DATA_CONVERTER=dataconverter-plugin

COPY --from=temporalio/admin-tools:1.14.0 /usr/local/bin/tctl /usr/local/bin/tctl
COPY --from=build-plugin /go/bin/dataconverter-plugin /go/bin/dataconverter-plugin