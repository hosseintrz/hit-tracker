package test

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestHealthCheck(t *testing.T) {
	var response *http.Response
	var err error
	timer := time.After(15 * time.Second)

	err = backoff.Retry(func() error {
		t.Log(fmt.Sprint("url : ", "localhost:80/ping"))
		response, err = http.Get("http://0.0.0.0:80/ping")
		if err != nil {
			select {
			case <-timer:
				t.FailNow()
			default:
				return err
			}
		}
		return nil
	}, backoff.NewExponentialBackOff())
	require.NoError(t, err, "http error")
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode, "http status code")

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err, "error reading body")

	t.Log("response is : ", string(body))
	require.JSONEqf(t, `{"Status":"ok"}`, string(body), "response body")

}

//func TestHealthCheck(t *testing.T) {
//	dockerPs := func() {
//		out, err := exec.Command("docker", "ps", "-a").Output()
//		if err != nil {
//			t.Fatal(err)
//		}
//		t.Log("ps output : ", string(out))
//	}
//
//	err := exec.Command("docker", "compose", "up", "-d").Run()
//	if err != nil{
//		t.Log("docker compose err -> ",err.Error())
//	}
//	require.NoError(t, err, "docker compose error")
//
//	defer func() {
//		err := exec.Command("docker", "compose", "down").Run()
//		require.NoError(t, err, "error in docker compose down")
//		dockerPs()
//	}()
//
//	dockerPs()
//
//	var response *http.Response
//	err = backoff.Retry(func() error {
//		t.Log(fmt.Sprint("url : ", "localhost:80/ping"))
//		response, err = http.Get("http://0.0.0.0:80/ping")
//		if err != nil {
//			return err
//		}
//		return nil
//	}, backoff.NewExponentialBackOff())
//	require.NoError(t, err, "http error")
//	defer response.Body.Close()
//
//	require.Equal(t, http.StatusOK, response.StatusCode, "http status code")
//
//	body, err := ioutil.ReadAll(response.Body)
//	require.NoError(t, err, "error reading body")
//
//	t.Log("response is : ", string(body))
//	require.JSONEqf(t, `{"Status":"ok"}`, string(body), "response body")
//
//}

//func TestHealthCheck(t *testing.T) {
//	pool, err := dockertest.NewPool("")
//	require.NoError(t, err, "couldn't connect to docker")
//
//	//starting cockroach db
//	roach, err := pool.RunWithOptions(&dockertest.RunOptions{
//		Repository: "cockroachdb/cockroach",
//		Tag:        "latest",
//		Name:       "roach2",
//		Hostname:   "db",
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"26257": {
//				{HostIP: "0.0.0.0", HostPort: "26257"},
//			},
//			"9090": {
//				{HostIP: "0.0.0.0", HostPort: "9090"},
//			},
//		},
//		Cmd: []string{"start-single-node"},
//	}, func(config *docker.HostConfig) {
//		config.AutoRemove = false
//		config.RestartPolicy = docker.RestartPolicy{
//			Name: "no",
//		}
//	})
//	require.NoError(t, err, "error starting db container")
//	t.Log("roach state : ", roach.Container.State.Status)
//
//	hitTracker, err := pool.RunWithOptions(&dockertest.RunOptions{
//		Repository: "hit-tracker",
//		Tag:        "latest",
//		Name:       "hit-tracker",
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"80": {
//				{HostIP: "0.0.0.0", HostPort: "9090"},
//			},
//		},
//		Env: []string{
//			"PGUSER=totoro",
//			"PGHOST=db",
//			"PGPORT=26257",
//			"PGDATABASE:mydb",
//		},
//	})
//
//	require.NoError(t, err, "error starting hit-tracker")
//
//	done := make(chan bool)
//	go func() {
//		for {
//			select {
//			case <-done:
//				return
//			default:
//				t.Log("hit-tracker state : ", hitTracker.Container.State.Status)
//				t.Log("roach state : ", roach.Container.State.Status)
//			}
//			time.Sleep(400 * time.Millisecond)
//		}
//	}()
//
//	t.Cleanup(func() {
//		t.Log("purging ...")
//		require.NoError(t, pool.Purge(hitTracker), "failed to remove hit-tracker")
//		require.NoError(t, pool.Purge(roach), "failed to remove roach db")
//		time.Sleep(1 * time.Second)
//		done <- true
//	})
//
//	var response *http.Response
//	err = pool.Retry(func() error {
//		t.Log(fmt.Sprint("url : ", "http://localhost/ping"))
//		response, err = http.Get(fmt.Sprint("http://localhost/ping"))
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	require.NoError(t, err, "http error")
//	defer response.Body.Close()
//
//	require.Equal(t, http.StatusOK, response.StatusCode, "http status code")
//
//	body, err := ioutil.ReadAll(response.Body)
//	require.NoError(t, err, "error reading body")
//
//	t.Log("response is : ", string(body))
//	require.JSONEqf(t, `{"Status":"ok"}`, string(body), "response body")
//}
