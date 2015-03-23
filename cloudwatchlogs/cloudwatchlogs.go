package cloudwatchlogs

import (
	"encoding/xml"
	"errors"
	"time"

	"github.com/AdRoll/goamz/aws"
)

type CloudWatchLogs struct {
	Service aws.AWSService
}

type InputLogEvent struct {
	Message   string
	Timestamp time.Time
}

type LogGroup struct {
	ARN               string
	CreationTime      time.Time
	LogGroupName      string
	MetricFilterCount float64
	RetentionInDays   int
	StoredBytes       float64
}

type LogStream struct {
	ARN                 string
	CreationTime        time.Time
	FirstEventTimestamp time.Time
	LastEventTimestamp  time.Time
	LastIngestionTime   time.Time
	LogStreamName       string
	StoredBytes         float64
	UploadSequenceToken string
}

type MetricTransformation struct {
	MetricName      string
	MetricNamespace string
	MetricValue     string
}

type MetricFilter struct {
	CreationTime          time.Time
	FilterName            string
	FilterPattern         string
	MetricTransformations []MetricTransformation
}

type MetricFilterMatchRecord struct {
	EventMessage    string
	EventNumber     float64
	ExtractedValues map[string]string
}

type OutputLogEvent struct {
	IngestionTime time.Time
	Message       string
	Timestamp     time.Time
}

type RejectedLogEventsInfo struct {
	ExpiredLogEventEndIndex  float64
	TooNewLogEventStartIndex float64
	TooOldLogEventEndIndex   float64
}

func NewCloudWatchLogs(auth aws.Auth, region aws.ServiceInfo) (*CloudWatchLogs, error) {
	service, err := aws.NewService(auth, region)
	if err != nil {
		return nil, err
	}
	return &CloudWatchLogs{
		Service: service,
	}, nil

}

func (c *CloudWatchLogs) query(method, path string, params map[string]string, resp interface{}) error {
	// Add basic Cloudwatch param
	params["Version"] = "2014-03-28"

	r, err := c.Service.Query(method, path, params)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return c.Service.BuildError(r)
	}
	err = xml.NewDecoder(r.Body).Decode(resp)
	return err
}

func (c *CloudWatchLogs) DeleteLogStream(logGroup *LogGroup, logStream *LogStream) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("DeleteLogStream")

	switch {
	case logGroup.LogGroupName == "":
		err = errors.New("No LogGroupName supplied")
	case logStream.LogStreamName == "":
		err = errors.New("No LogStreamName supplied")
	}
	if err != nil {
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName
	params["LogStreamName"] = logStream.LogStreamName

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) DeleteLogGroup(logGroup *LogGroup) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("DeleteLogGroup")

	if logGroup.LogGroupName == "" {
		err = errors.New("No LogGroupName supplied")
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) CreateLogStream(logGroup *LogGroup, logStream *LogStream) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("CreateLogStream")

	switch {
	case logGroup.LogGroupName == "":
		err = errors.New("No LogGroupName supplied")
	case logStream.LogStreamName == "":
		err = errors.New("No LogStreamName supplied")
	}
	if err != nil {
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName
	params["LogStreamName"] = logStream.LogStreamName

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) CreateLogGroup(logGroup *LogGroup) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("CreateLogGroup")

	if logGroup.LogGroupName == "" {
		err = errors.New("No LogGroupName supplied")
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}
