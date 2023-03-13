package AuditLog

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAuditLogsSuite(t *testing.T) {
	suite.Run(t, new(AuditLogsTestSuite))
}
