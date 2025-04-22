package tests

// Package tests provides test types and functions for yac-p

// GetTestConfigLoader returns a function that loads a static test configuration
func GetTestConfigLoader() func() ([]byte, error) {
	return func() ([]byte, error) {
		return []byte(`
apiVersion: v1alpha1
discovery:
  exportedTagsOnMetrics:
    AWS/EC2:
      - Environment
  jobs:
    - type: AWS/EC2
      roles:
        - roleArn: "arn:aws:iam::12345678912:role/test"
      regions: 
        - eu-north-1
      metrics:
        - name: CPUUtilization
        - name: NetworkIn
        - name: NetworkOut
        - name: EBSReadOps
        - name: EBSWriteOps
      dimensionNameRequirements:
        - InstanceId
      includeContextOnInfoMetrics: true
      statistics:
        - Average
      period: 300
      length: 300
`), nil
	}
}
