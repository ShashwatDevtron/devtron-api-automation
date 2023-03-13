package AuditLog

import (
	Base "automation-suite/testUtils"
	"github.com/stretchr/testify/suite"
)

type AuditLogsTestSuite struct {
	suite.Suite
	authToken string
}

func (suite *AuditLogsTestSuite) SetupSuite() {
	suite.authToken = Base.GetAuthToken()
}
