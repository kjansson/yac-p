module github.com/kjansson/yac-p

replace github.com/kjansson/yac-p/pkg/config => /Users/kjansson/dev/swb-observability-poc/yac-p/pkg/config

replace github.com/kjansson/yac-p/pkg/types => /Users/kjansson/dev/swb-observability-poc/yac-p/pkg/types

replace github.com/kjansson/yac-p/pkg/prom => /Users/kjansson/dev/swb-observability-poc/yac-p/pkg/prom

replace github.com/kjansson/yac-p/pkg/yace => /Users/kjansson/dev/swb-observability-poc/yac-p/pkg/yace

replace github.com/kjansson/yac-p/pkg/logger => /Users/kjansson/dev/swb-observability-poc/yac-p/pkg/logger

replace github.com/kjansson/yac-p/pkg/tests => /Users/kjansson/dev/swb-observability-poc/yac-p/pkg/tests

replace github.com/kjansson/yac-p/pkg/loaders => /Users/kjansson/dev/swb-observability-poc/yac-p/pkg/loaders

go 1.23.0

toolchain go1.23.8

require (
	github.com/aws/aws-lambda-go v1.48.0
	github.com/aws/aws-sdk-go v1.55.6
	github.com/golang/snappy v1.0.0
	github.com/prometheus-community/yet-another-cloudwatch-exporter v0.62.1
	github.com/prometheus/client_golang v1.22.0
	github.com/prometheus/client_model v0.6.2
	github.com/prometheus/prometheus v0.303.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.32.7 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.28.7 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.48 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/amp v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.28.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.24.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.51.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.43.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.45.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.198.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.38.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.25.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/shield v1.29.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/storagegateway v1.34.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.3 // indirect
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/grafana/regexp v0.0.0-20240607082908-2cb410fa05da // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
