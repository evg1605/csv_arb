package arb_converter

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLoadArb(t *testing.T) {
	arbData, err := LoadArb(createLogger(), "test_data", "en")

	require.NoError(t, err)
	require.NotNil(t, arbData)
}

func createLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	return logger
}
