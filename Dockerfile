FROM golang:1.17 AS build

WORKDIR /go/src/background-checks

COPY go.mod go.sum ./

RUN go mod download

COPY activities ./activities
COPY api ./api
COPY cli ./cli
COPY temporal ./temporal
COPY types ./types
COPY utils ./utils
COPY ui ./ui
COPY workflows ./workflows

RUN go install -v ./cli/bgc-dataconverter-plugin
RUN go install -v ./cli/bgc-backend
RUN go install -v ./cli/bgc-company
RUN go install -v ./cli/bgc-candidate
RUN go install -v ./cli/bgc-researcher
RUN go install -v ./cli/thirdparty-simulator

FROM golang:1.17

ENV TEMPORAL_CLI_PLUGIN_DATA_CONVERTER=bgc-dataconverter-plugin

COPY --from=temporalio/admin-tools:1.14.0 /usr/local/bin/tctl /usr/local/bin/tctl

COPY --from=build /go/bin/bgc-dataconverter-plugin /usr/local/bin/bgc-dataconverter-plugin
COPY --from=build /go/bin/bgc-backend /usr/local/bin/bgc-backend
COPY --from=build /go/bin/bgc-company /usr/local/bin/bgc-company
COPY --from=build /go/bin/bgc-candidate /usr/local/bin/bgc-candidate
COPY --from=build /go/bin/bgc-researcher /usr/local/bin/bgc-researcher
COPY --from=build /go/bin/thirdparty-simulator /usr/local/bin/thirdparty-simulator
