FROM golang:1.17

WORKDIR /go/src/background-checks

COPY go.mod go.sum ./

RUN go mod download

COPY activities ./activities/
COPY api ./api/
COPY cli ./cli/
COPY cmd ./cmd/
COPY config ./config/
COPY mappings ./mappings/
COPY mocks ./mocks/
COPY queries ./queries/
COPY signals ./signals/
COPY thirdparty ./thirdparty/
COPY types ./types/
COPY workflows ./workflows/
COPY main.go ./

RUN go install -v .
RUN go install -v ./cli/bgc-candidate
RUN go install -v ./cli/bgc-company
RUN go install -v ./cli/bgc-researcher