package cloudwatchlogs_test

import (
	"testing"

	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/cloudwatchlogs"
	"github.com/AdRoll/goamz/testutil"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	cwl *cloudwatchlogs.CloudWatchLogs
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.cwl, _ = cloudwatchlogs.NewCloudWatchLogs(auth, aws.ServiceInfo{testServer.URL, aws.V2Signature})
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

// TODO: add invalid entities to test proper validation?
func getSomeTestLogGroup() *cloudwatchlogs.LogGroup {
	logGroup := new(cloudwatchlogs.LogGroup)
	logGroup.LogGroupName = "someLogGroupName"
	return logGroup
}

func getSomeTestLogStream() *cloudwatchlogs.LogStream {
	logStream := new(cloudwatchlogs.LogStream)
	logStream.LogStreamName = "someLogStreamName"
	return logStream
}

func getSomeTestMetricTransformation() *cloudwatchlogs.MetricTransformation {
	metricTransformation := new(cloudwatchlogs.MetricTransformation)
	metricTransformation.MetricName = "someMetricName"
	metricTransformation.MetricNamespace = "someMetricNamespace"
	metricTransformation.MetricValue = "someMetricValue"
	return metricTransformation
}

func getSomeTestMetricFilter() *cloudwatchlogs.MetricFilter {
	metricFilter := new(cloudwatchlogs.MetricFilter)
	metricFilter.FilterName = "someMetricFilterName"
	metricFilter.FilterPattern = "someMetricFilterPattern"

	metricTransformation := getSomeTestMetricTransformation()
	metricFilter.MetricTransformations = []cloudwatchlogs.MetricTransformation{*metricTransformation}

	return metricFilter
}

func (s *S) TestPutMetricFilter(c *C) {
	testServer.Response(200, nil, "<RequestId>123</RequestId>")

	logGroup := getSomeTestLogGroup()
	metricFilter := getSomeTestMetricFilter()

	_, err := s.cwl.PutMetricFilter(logGroup, metricFilter)
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Form["Action"], DeepEquals, []string{"PutMetricFilter"})
	c.Assert(req.Form["FilterName"], DeepEquals, []string{"someMetricFilterName"})
	c.Assert(req.Form["FilterPattern"], DeepEquals, []string{"someMetricFilterPattern"})
	c.Assert(req.Form["LogGroupName"], DeepEquals, []string{"someLogGroupName"})

	metricsPrefix := "MetricTransformations.member.1"

	c.Assert(req.Form[metricsPrefix+".MetricName"], DeepEquals, []string{metricFilter.MetricTransformations[0].MetricName})
	c.Assert(req.Form[metricsPrefix+".MetricNamespace"], DeepEquals, []string{metricFilter.MetricTransformations[0].MetricNamespace})
	c.Assert(req.Form[metricsPrefix+".MetricValue"], DeepEquals, []string{metricFilter.MetricTransformations[0].MetricValue})
}

func (s *S) TestDeleteLogStream(c *C) {
	testServer.Response(200, nil, "<RequestId>123</RequestId>")

	logGroup := getSomeTestLogGroup()
	logStream := getSomeTestLogStream()

	_, err := s.cwl.DeleteLogStream(logGroup, logStream)
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteLogStream"})
	c.Assert(req.Form["LogGroupName"], DeepEquals, []string{"someLogGroupName"})
	c.Assert(req.Form["LogStreamName"], DeepEquals, []string{"someLogStreamName"})
}

func (s *S) TestDeleteLogGroup(c *C) {
	testServer.Response(200, nil, "<RequestId>123</RequestId>")

	logGroup := getSomeTestLogGroup()

	_, err := s.cwl.DeleteLogGroup(logGroup)
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Form["Action"], DeepEquals, []string{"DeleteLogGroup"})
	c.Assert(req.Form["LogGroupName"], DeepEquals, []string{"someLogGroupName"})
}

func (s *S) TestCreateLogGroup(c *C) {
	testServer.Response(200, nil, "<RequestId>123</RequestId>")

	logGroup := getSomeTestLogGroup()

	_, err := s.cwl.CreateLogGroup(logGroup)
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateLogGroup"})
	c.Assert(req.Form["LogGroupName"], DeepEquals, []string{"someLogGroupName"})
}

func (s *S) TestCreateLogStream(c *C) {
	testServer.Response(200, nil, "<RequestId>123</RequestId>")

	logGroup := getSomeTestLogGroup()
	logStream := getSomeTestLogStream()

	_, err := s.cwl.CreateLogStream(logGroup, logStream)
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateLogStream"})
	c.Assert(req.Form["LogGroupName"], DeepEquals, []string{"someLogGroupName"})
	c.Assert(req.Form["LogStreamName"], DeepEquals, []string{"someLogStreamName"})
}
