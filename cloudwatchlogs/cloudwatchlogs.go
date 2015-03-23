package cloudwatchlogs

import (
	"encoding/xml"
	"errors"
	"strconv"
	"time"

	"github.com/AdRoll/goamz/aws"
	"github.com/feyeleanor/sets"
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

type DescribeLogGroupsRequest struct {
	Limit              int
	LogGroupNamePrefix string
	NextToken          string
}

type DescribeLogGroupsResult struct {
	LogGroups []LogGroup `xml:"LogGroups>member"`
	NextToken string     `xml:"NextToken"`
}

type DescribeLogGroupsResponse struct {
	DescribeLogGroupsResult DescribeLogGroupsResult
	ResponseMetadata        aws.ResponseMetadata
}

var validRetentionTimeInDays = sets.ISet(
	1,
	3,
	5,
	7,
	14,
	30,
	60,
	90,
	120,
	150,
	180,
	365,
	400,
	545,
	731,
	1827,
	3653,
)

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

func (c *CloudWatchLogs) DeleteMetricFilter(logGroup *LogGroup, metricFilter *MetricFilter) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("DeleteMetricFilter")

	switch {
	case logGroup.LogGroupName == "":
		err = errors.New("No LogGroupName supplied")
	case metricFilter.FilterName == "":
		err = errors.New("No FilterName supplied")
	}
	if err != nil {
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName
	params["FilterName"] = metricFilter.FilterName

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) DeleteRetentionPolicy(logGroup *LogGroup) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("DeleteRetentionPolicy")

	if logGroup.LogGroupName == "" {
		err = errors.New("No LogGroupName supplied")
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) DescribeLogGroups(req *DescribeLogGroupsRequest) (result *DescribeLogGroupsResponse, err error) {
	params := aws.MakeParams("DescribeLogGroups")

	params["Limit"] = strconv.Itoa(req.Limit)
	params["LogGroupNamePrefix"] = req.LogGroupNamePrefix
	params["NextToken"] = req.NextToken

	result = new(DescribeLogGroupsResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) DescribeLogStreams(logGroup *LogGroup, retentionInDays int) (result *aws.BaseResponse, err error) {
	return
}

func (c *CloudWatchLogs) DescribeMetricFilters(logGroup *LogGroup, retentionInDays int) (result *aws.BaseResponse, err error) {
	return
}

func (c *CloudWatchLogs) GetLogEvents(logGroup *LogGroup, retentionInDays int) (result *aws.BaseResponse, err error) {
	return
}

func (c *CloudWatchLogs) PutLogEvents(logGroup *LogGroup, retentionInDays int) (result *aws.BaseResponse, err error) {
	return
}

func (c *CloudWatchLogs) PutMetricFilter(logGroup *LogGroup, metricFilter *MetricFilter) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("PutMetricFilter")

	switch {
	case logGroup.LogGroupName == "":
		err = errors.New("No LogGroupName supplied")
	case metricFilter.FilterName == "":
		err = errors.New("No FilterName supplied")
	case metricFilter.FilterPattern == "":
		err = errors.New("No FilterPattern supplied")
	case len(metricFilter.MetricTransformations) < 1:
		err = errors.New("No MetricTransformations supplied")
	}
	if err != nil {
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName
	params["FilterName"] = metricFilter.FilterName
	params["FilterPattern"] = metricFilter.FilterPattern

	for i, d := range metricFilter.MetricTransformations {
		prefix := "MetricTransformations.member." + strconv.Itoa(i+1)
		params[prefix+".MetricName"] = d.MetricName
		params[prefix+".MetricNamespace"] = d.MetricNamespace
		params[prefix+".MetricValue"] = d.MetricValue
	}

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) PutRetentionPolicy(logGroup *LogGroup, retentionInDays int) (result *aws.BaseResponse, err error) {
	params := aws.MakeParams("PutRetentionPolicy")

	switch {
	case logGroup.LogGroupName == "":
		err = errors.New("No LogGroupName supplied")
	case !validRetentionTimeInDays.Member(retentionInDays):
		err = errors.New("RetentionInDays is not a valid value")
	}
	if err != nil {
		return
	}

	params["LogGroupName"] = logGroup.LogGroupName
	params["RetentionInDays"] = strconv.Itoa(retentionInDays)

	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}

func (c *CloudWatchLogs) TestMetricFilter(logGroup *LogGroup, retentionInDays int) (result *aws.BaseResponse, err error) {
	return
}
