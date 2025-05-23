// Package prom provides a client for persisting metrics to a Prometheus remote write endpoint. It implements the types.MetricPersister interface.
package prom

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/golang/snappy"
	"github.com/kjansson/yac-p/v3/pkg/types"
	"github.com/prometheus/prometheus/prompb"
)

type PromClient struct {
	RemoteWriteURL   string // URL of the Prometheus remote write endpoint
	AuthType         string // Type of authentication to use (AWS, BASIC, TOKEN)
	AuthToken        string // Token to use for authentication (if using TOKEN auth)
	Username         string // Username to use for authentication (if using BASIC auth)
	Password         string // Password to use for authentication (if using BASIC auth)
	Region           string // AWS region to use for authentication (if using AWS auth)
	PrometheusRegion string // AWS region of the Prometheus remote write endpoint (if using Amazon Managed Prometheus)
	AWSRoleARN       string // ARN of the AWS role to assume for remote write (if using Amazon Managed Prometheus cross-account)
}

func NewPromClient(
	remoteWriteURL string,
	authType string,
	authToken string,
	username string,
	password string,
	region string,
	prometheusRegion string,
	awsRoleARN string,
) (*PromClient, error) {

	if remoteWriteURL == "" {
		return nil, fmt.Errorf("prometheus remote write URL must be set")
	}

	if authType == "BASIC" { // Basic auth requires username and password
		if username == "" || password == "" {
			return nil, fmt.Errorf("username and password must be set for BASIC auth")
		}
	}
	if authType == "TOKEN" { // Token auth requires token
		if authToken == "" {
			return nil, fmt.Errorf("auth token must be set for TOKEN auth")
		}
	}

	return &PromClient{
		RemoteWriteURL:   remoteWriteURL,
		AuthType:         authType,
		AuthToken:        authToken,
		Username:         username,
		Password:         password,
		Region:           region,
		PrometheusRegion: prometheusRegion,
		AWSRoleARN:       awsRoleARN,
	}, nil
}

// PeristMetrics creates a Prometheus remote write request and sends it to the remote write URL
func (p *PromClient) PersistMetrics(timeSeries []prompb.TimeSeries, logger types.Logger) error {

	logger.Log("debug", "Sending timeseries", slog.Int("timeseries_count", len(timeSeries)))
	logger.Log("debug", "Auth type", slog.String("auth_type", p.AuthType))

	r := &prompb.WriteRequest{
		Timeseries: timeSeries,
	}
	tsProto, err := r.Marshal()
	if err != nil {
		return err
	}

	encoded := snappy.Encode(nil, tsProto)
	body := bytes.NewReader(encoded)

	req, err := http.NewRequest("POST", p.RemoteWriteURL, body)
	if err != nil {
		return err
	}

	switch p.AuthType {
	case "AWS":
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(p.Region),
		})
		if err != nil {
			return err
		}

		var awsCredentials *credentials.Credentials
		if p.AWSRoleARN != "" {
			logger.Log("debug", "Using AWS role", slog.String("role_arn", p.AWSRoleARN))
			awsCredentials = stscreds.NewCredentials(sess, p.AWSRoleARN, func(p *stscreds.AssumeRoleProvider) {
				host, err := os.Hostname()
				if err != nil {
					host = "unknown"
				}
				p.RoleSessionName = "aws-sigv4-proxy-" + host
				logger.Log("debug", "Using AWS role session name", slog.String("role_session_name", p.RoleSessionName))
			})
		} else {
			awsCredentials = sess.Config.Credentials
		}

		_, err = v4.NewSigner(awsCredentials).Sign(req, body, "aps", p.PrometheusRegion, time.Now())
		if err != nil {
			return err
		}
	case "BASIC":
		logger.Log("debug", "Using basic auth")
		req.SetBasicAuth(p.Username, p.Password)
	case "TOKEN":
		logger.Log("debug", "Using token auth")
		req.Header.Set("Authorization", "Bearer "+p.AuthToken)
	}

	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

	logger.Log("debug", "Sending request", slog.String("url", p.RemoteWriteURL), slog.Int("body_size", len(encoded)))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send metrics: %s", response.Status)
	}
	logger.Log("debug", "Response", slog.String("status", response.Status), slog.Int("status_code", response.StatusCode))

	return nil
}
