ARG BASE_BUILDER_IMAGE=golang:1.27-alpine3.24
ARG TEMPORAL_TAG=1.10.0

FROM ${BASE_BUILDER_IMAGE} AS base

WORKDIR /go/src/background-checks

RUN apk upgrade

COPY go.mod go.sum ./

RUN go mod download

FROM base AS build

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
RUN go install -v ./temporal/dataconverter-server

FROM temporalio/temporal:${TEMPORAL_TAG} AS app

COPY --from=build /go/bin/dataconverter-server /usr/local/bin/dataconverter-server

COPY --from=build /go/bin/bgc-backend /usr/local/bin/bgc-backend
COPY --from=build /go/bin/bgc-company /usr/local/bin/bgc-company
COPY --from=build /go/bin/bgc-candidate /usr/local/bin/bgc-candidate
COPY --from=build /go/bin/bgc-researcher /usr/local/bin/bgc-researcher

ENTRYPOINT ["/usr/bin/env"]
CMD ["/bin/sh"]