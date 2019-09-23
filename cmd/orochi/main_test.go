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

	time.Sleep(500 * time.Millisecond)

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

	c := http.Client{}
	p := strconv.Itoa(portlist[0])
	key := "hoge"
	value := "fuga"
	resp, err := c.Post(fmt.Sprintf("http://127.0.0.1:%s/%s", p, key), "", bytes.NewBuffer([]byte(value)))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	resp, err = c.Get(fmt.Sprintf("http://127.0.0.1:%s/%s", p, key))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	retVal, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(retVal) != "fuga" {
		t.Errorf("unexpected result. got [%s], want [%s]", string(retVal), "fugaa")
	}
}
