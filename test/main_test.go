package test

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "couldn't connect to docker")

	container, err := pool.Run("hit-tracker", "latest", []string{})
	require.NoError(t, err, "couldn't start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(container), "failed to remove container")
	})
	var response *http.Response
	err = pool.Retry(func() error {
		t.Log(fmt.Sprint("url : ", "http://localhost:", container.GetPort("9090/tcp"), "/ping"))
		response, err = http.Get(fmt.Sprint("http://localhost:", container.GetPort("9090/tcp"), "/ping"))
		if err != nil {
			return err
		}
		return nil
	})
	require.NoError(t, err, "http error")
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode, "http status code")

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err, "error reading body")

	require.JSONEqf(t, `{"Status":"ok"}`, string(body), "response body")
}
