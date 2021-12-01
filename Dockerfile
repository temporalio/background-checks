FROM golang:1.17-alpine

WORKDIR /go/src/background-checks

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
COPY workflows ./workflows/

RUN go install -v ./cli/bgc-backend
RUN go install -v ./cli/bgc-company
RUN go install -v ./cli/bgc-candidate
RUN go install -v ./cli/bgc-researcher
RUN go install -v ./cli/thirdparty-simulator