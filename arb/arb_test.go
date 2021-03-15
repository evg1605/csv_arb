package arb

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLoadArb(t *testing.T) {
	arbData, err := LoadArb(createLogger(), "test_data", "en")

	require.NoError(t, err)
	require.NotNil(t, arbData)

	require.Len(t, arbData.Cultures, 2)
	require.Contains(t, arbData.Cultures, "ru")
	require.Contains(t, arbData.Cultures, "en")

	require.Contains(t, arbData.Items, "aa1")
	require.Contains(t, arbData.Items, "aa2")
	require.Contains(t, arbData.Items, "aa3")
	require.Contains(t, arbData.Items, "aa4")
	require.Contains(t, arbData.Items, "myName")

	require.Len(t, arbData.Items["myName"].Cultures, 2)
	require.NotEmpty(t, arbData.Items["myName"].Description)
	require.Len(t, arbData.Items["myName"].Parameters, 2)
	require.Contains(t, arbData.Items["myName"].Parameters, "p1")
	require.Contains(t, arbData.Items["myName"].Parameters, "par2")

	require.Len(t, arbData.Items["aa4"].Cultures, 1)
	require.Empty(t, arbData.Items["aa3"].Parameters)

}

func createLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	return logger
}
