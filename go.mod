module github.com/ohkinozomu/redash-visualizer

go 1.16

require (
	github.com/goccy/go-graphviz v0.0.9
	github.com/snowplow-devops/redash-client-go v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.2.1
)

replace github.com/snowplow-devops/redash-client-go => github.com/ohkinozomu/redash-client-go v0.4.2
