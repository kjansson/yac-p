#!/bin/bash



#s3api create-bucket --bucket sample-bucket --region eu-north-1 --create-bucket-configuration LocationConstraint=eu-north-1
#awslocal s3api put-object --bucket sample-bucket --key config.yaml --body ~/dev/swb-observability-poc/yac-p/config.yaml


# awslocal cloudwatch put-metric-data --namespace "AWS/EC2" --metric-name "TestBytesIn" --value 0.1

# echo "BUILD"
# GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -C ./yac-p -tags lambda.norpc -mod=readonly -ldflags='-s -w' -o ../bootstrap 
# echo "BUILD DONE"

# zip -r ./yac-p.zip ./bootstrap

sam build -t sam-template.yaml

cp .aws-sam/build/yacp/bootstrap yac-p
sam local invoke yacp -t sam-template.yaml






#rm ./bootstrap ./yac-p.zip