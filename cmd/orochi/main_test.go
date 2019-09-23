package main

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

var portlist = []int{3000, 3001, 3002}

func setup() func() []error {
	serverlist := []*orochi.Server{}
	for _, v := range portlist {
		go func(p int) {
			o := &orochi.Server{PortList: portlist}
			serverlist = append(serverlist, o)
			o.Serve(p)
		}(v)
	}

	time.Sleep(100 * time.Millisecond)

	return func() []error {
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
	teardown := setup()

	defer func() {
		errlist := teardown()
		if len(errlist) != 0 {
			t.Log(errlist)
		}
	}()

	err := post(portlist[0], "hoge", "fuga")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, port := range portlist {
		ret, err := get(port, "hoge")
		if err != nil {
			t.Fatal(err)
		}

		if string(ret) != "fuga" {
			t.Errorf("unexpected result. got [%s], want [%s]", string(ret), "fuga")
		}
	}

	err = post(portlist[1], "hoge", "piyo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, port := range portlist {
		ret, err := get(port, "hoge")
		if err != nil {
			t.Fatal(err)
		}

		if string(ret) != "piyo" {
			t.Errorf("unexpected result. got [%s], want [%s]", string(ret), "piyo")
		}
	}

	err = post(portlist[2], "foo", "bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, port := range portlist {
		ret, err := get(port, "foo")
		if err != nil {
			t.Fatal(err)
		}

		if string(ret) != "bar" {
			t.Errorf("unexpected result. got [%s], want [%s]", string(ret), "bar")
		}
	}
}

func post(port int, key, value string) error {
	c := http.Client{}
	p := strconv.Itoa(port)

	resp, err := c.Post(fmt.Sprintf("http://127.0.0.1:%s/%s", p, key), "", bytes.NewBuffer([]byte(value)))
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

func get(port int, key string) (string, error) {
	c := http.Client{}
	p := strconv.Itoa(port)

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
