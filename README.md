# AWS Cloudwatch Alarm Manager
This service is meant to set the Cloudwatch Alarm to a desired state (which can be invoked by an eventbrigde or etc) to allow the cloudwatch alarm to re-invoke a lambda function. It's a workaround to work as the scale-in & scale-out cooldown period for the [docdb-autoscaler](https://github.com/cheelim1/docdb-autoscaler) service.
This service was built to work together with the `docdb-autoscaler` service [Link](https://github.com/cheelim1/docdb-autoscaler). However, it can be used for other services as well.

### Testing
```
   go fmt ./...
   go clean -cache -testcache
   go test ./... -v
```

### Debugging
1. Go to the AWS Lambda function -> Monitor & check if the Lambda function was invoked
2. Further debug using cloudwatch logs.