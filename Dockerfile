FROM golang:1.17 AS base

WORKDIR /go/src/background-checks

RUN apt-get update && apt-get install -y socat && rm -rf /var/lib/apt/lists/*

COPY --from=temporalio/admin-tools /usr/local/bin/tctl /usr/local/bin/tctl

COPY go.mod go.sum ./

RUN go mod download

COPY activities ./activities/
COPY api ./api/
COPY cli ./cli/
COPY config ./config/
COPY mappings ./mappings/
COPY mocks ./mocks/
COPY queries ./queries/
COPY signals ./signals/
COPY types ./types/
COPY temporal ./temporal
COPY workflows ./workflows/

RUN go install -v ./cli/bgc-dataconverter-plugin
ENV TEMPORAL_CLI_PLUGIN_DATA_CONVERTER=bgc-dataconverter-plugin

RUN go install -v ./cli/bgc-backend
RUN go install -v ./cli/bgc-company
RUN go install -v ./cli/bgc-candidate
RUN go install -v ./cli/bgc-researcher
RUN go install -v ./cli/thirdparty-simulator
