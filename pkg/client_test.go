package pkg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProcessTestSuite struct {
	suite.Suite
	PygeoapiContainer PygeoapiContainer
}

func (suite *ProcessTestSuite) SetupSuite() {
	container, err := NewPygeoapiContainer()
	suite.Require().NoError(err)
	suite.PygeoapiContainer = container
}

func (s *ProcessTestSuite) TearDownSuite() {
	err := s.PygeoapiContainer.testcontainer.Terminate(context.Background())
	s.Require().NoError(err)
}

func (s *ProcessTestSuite) TestListProcesses() {
	client, err := NewProcessesClient(s.PygeoapiContainer.connectionUrl)
	s.Require().NoError(err)

	processesList, err := client.ListProcesses()
	s.Require().NoError(err)
	s.Require().NotEmpty(processesList.ProcessInfo)
}

func (s *ProcessTestSuite) TestGetProcessInfo() {
	client, err := NewProcessesClient(s.PygeoapiContainer.connectionUrl)
	s.Require().NoError(err)

	processesList, err := client.ListProcesses()
	s.Require().NoError(err)

	processInfo, err := client.GetProcessInfo(processesList.ProcessInfo[0].Id)
	s.Require().NoError(err)
	s.Require().NotEmpty(processInfo.Id)
}

func (s *ProcessTestSuite) TestExecuteSync() {
	client, err := NewProcessesClient(s.PygeoapiContainer.connectionUrl)
	s.Require().NoError(err)

	info, err := client.GetProcessInfo("hello-world")
	s.Require().NoError(err)
	s.Require().NotEmpty(info.Id)
	s.Require().Contains(info.JobControlOptions, SyncSupport)

	response, err := client.ExecuteSync("hello-world", map[string]any{"name": "test"})
	s.Require().NoError(err)
	s.Require().Equal("Hello test!", response.Value)
}

// Run the entire test suite
func TestNabuIntegrationClientSuite(t *testing.T) {
	suite.Run(t, new(ProcessTestSuite))
}
