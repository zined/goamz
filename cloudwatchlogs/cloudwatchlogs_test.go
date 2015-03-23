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

func getTestLogGroup() *cloudwatchlogs.LogGroup {
	logGroup := new(cloudwatchlogs.LogGroup)
	logGroup.LogGroupName = "someLogGroupName"
	return logGroup
}

func getTestLogStream() *cloudwatchlogs.LogStream {
	logStream := new(cloudwatchlogs.LogStream)
	logStream.LogStreamName = "someLogStreamName"
	return logStream
}

func (s *S) TestDeleteLogStream(c *C) {
	testServer.Response(200, nil, "<RequestId>123</RequestId>")

	logGroup := getTestLogGroup()
	logStream := getTestLogStream()

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

	logGroup := getTestLogGroup()

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

	logGroup := getTestLogGroup()

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

	logGroup := getTestLogGroup()
	logStream := getTestLogStream()

	_, err := s.cwl.CreateLogStream(logGroup, logStream)
	c.Assert(err, IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, Equals, "POST")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateLogStream"})
	c.Assert(req.Form["LogGroupName"], DeepEquals, []string{"someLogGroupName"})
	c.Assert(req.Form["LogStreamName"], DeepEquals, []string{"someLogStreamName"})
}
