package cloudwatchlogs_test

import (
	"testing"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/cloudwatchlogs"
	"github.com/goamz/goamz/testutil"

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
