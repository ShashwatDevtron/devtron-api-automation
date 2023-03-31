package K8sPodsRouter

func (suite *K8sPodsTestSuite) TestClassGetLogs() {
	suite.Run("A=1=GetLogs", func() {
		suite.HitCheckLogsApi("", "", "", "", "", "", "", "")
	})
}
