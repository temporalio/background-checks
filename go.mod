module github.com/temporalio/background-checks

go 1.16

replace go.temporal.io/sdk => github.com/temporalio/sdk-go v1.13.1-0.20220207161017-b2a682807bad

require (
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-plugin v1.4.3
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/uber-go/tally/v4 v4.1.1
	github.com/xhit/go-simple-mail/v2 v2.10.0
	go.temporal.io/api v1.6.1-0.20211110205628-60c98e9cbfe2
	go.temporal.io/sdk v1.13.1-0.20220206141821-eb0f2d2a6719
	go.temporal.io/sdk/contrib/tally v0.1.0
	go.temporal.io/server v1.13.1
)
