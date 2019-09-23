package orochi_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/pankona/orochi"
)

func setup() ([]*orochi.Orochi, func() []error) {
	var (
		portlist   = []int{3000, 3001, 3002}
		serverlist = []*orochi.Orochi{}
	)
	for _, v := range portlist {
		go func(p int) {
			o := &orochi.Orochi{PortList: portlist}
			serverlist = append(serverlist, o)
			err := o.Serve(p)
			if err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("failed to launch server: %v", err))
			}
		}(v)
	}

	time.Sleep(100 * time.Millisecond)

	return serverlist, func() []error {
		errorlist := []error{}
		for i := range serverlist {
			err := serverlist[i].Shutdown()
			if err != nil {
				errorlist = append(errorlist, err)
			}
		}
		return errorlist
	}
}

func TestTypicalUsecase(t *testing.T) {
	serverlist, teardown := setup()
	defer func() {
		errlist := teardown()
		if len(errlist) != 0 {
			t.Log(errlist)
		}
	}()

	ret, err := get(serverlist[0], "hoge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ret != "" {
		t.Fatalf("unexpected result. got [%s], want [%s]", ret, "")
	}

	err = post(serverlist[0], "hoge", "fuga")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, server := range serverlist {
		ret, err := get(server, "hoge")
		if err != nil {
			t.Fatal(err)
		}

		if ret != "fuga" {
			t.Errorf("unexpected result. got [%s], want [%s]", ret, "fuga")
		}
	}

	err = post(serverlist[1], "hoge", "piyo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, server := range serverlist {
		ret, err := get(server, "hoge")
		if err != nil {
			t.Fatal(err)
		}

		if ret != "piyo" {
			t.Errorf("unexpected result. got [%s], want [%s]", ret, "piyo")
		}
	}

	err = post(serverlist[2], "foo", "bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, server := range serverlist {
		ret, err := get(server, "foo")
		if err != nil {
			t.Fatal(err)
		}

		if ret != "bar" {
			t.Errorf("unexpected result. got [%s], want [%s]", ret, "bar")
		}
	}

	err = restart(serverlist[1])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, server := range serverlist {
		ret, err := get(server, "foo")
		if err != nil {
			t.Fatal(err)
		}

		if ret != "bar" {
			t.Errorf("unexpected result. got [%s], want [%s]", ret, "bar")
		}
	}
}

func TestUnsupportedPath(t *testing.T) {
	serverlist, teardown := setup()
	defer func() {
		errlist := teardown()
		if len(errlist) != 0 {
			t.Log(errlist)
		}
	}()

	c := http.Client{}
	p := strconv.Itoa(serverlist[0].Port())
	key := "hoge"

	resp, err := c.Get(fmt.Sprintf("http://127.0.0.1:%s/unknown/path/%s", p, key))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected result: got [%d], want [%d]", resp.StatusCode, http.StatusNotFound)
	}
}

func restart(o *orochi.Orochi) error {
	err := o.Shutdown()
	if err != nil {
		return err
	}

	go func() {
		_ = o.Serve(o.Port())
	}()

	time.Sleep(500 * time.Millisecond)
	return nil
}

func post(o *orochi.Orochi, key, value string) error {
	c := http.Client{}
	p := strconv.Itoa(o.Port())

	resp, err := c.Post(fmt.Sprintf("http://127.0.0.1:%s/%s", p, key), "", bytes.NewBuffer([]byte(value)))
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

func get(o *orochi.Orochi, key string) (string, error) {
	c := http.Client{}
	p := strconv.Itoa(o.Port())

	resp, err := c.Get(fmt.Sprintf("http://127.0.0.1:%s/%s", p, key))
	if err != nil {
		return "", err
	}

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(ret), nil
}
