package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	kvstore  = map[string]string{}
	myport   = flag.Int("port", 3000, "specify port number of http server")
	portList = []int{3000, 3001, 3002}
)

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)
	log.Printf("server will start on port: %d\n", *myport)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*myport), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ss := strings.Split(r.URL.Path, "/")
		key := ss[len(ss)-1]

		v, ok := kvstore[key]
		if !ok {
			q := r.URL.Query()
			if q["asked"] != nil {
				break
			}

			// ask to other server
			log.Println("missed. ask to other server")
			for _, p := range portList {
				if *myport == p {
					log.Println("skip because this is me")
					continue
				}

				v, ok := askGet(p, key)
				if !ok {
					log.Printf("missed by port: %d", p)
					continue
				}

				log.Printf("hit on other server: %d", p)
				kvstore[key] = string(v)
				log.Printf("stored: %s\n", string(v))

				w.WriteHeader(200)
				w.Write([]byte(v))
				return
			}
		} else {
			w.WriteHeader(200)
			w.Write([]byte(v))
			log.Printf("return %s", v)
			return
		}

		w.WriteHeader(404)
	case "POST":
		ss := strings.Split(r.URL.Path, "/")
		key := ss[len(ss)-1]
		v, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic("TODO: error handling")
		}

		kvstore[key] = string(v)
		log.Printf("stored: %s\n", string(v))

		q := r.URL.Query()
		if q["asked"] != nil {
			break
		}

		for _, p := range portList {
			err := askPost(p, key, string(v))
			if err != nil {
				log.Printf("failed to askPost: %v", err)
			}
		}

		w.WriteHeader(200)
	default:
		panic("TODO: implement to return 404")
	}
}

func askGet(port int, key string) (string, bool) {
	c := http.Client{}
	p := strconv.Itoa(port)
	resp, err := c.Get(fmt.Sprintf("http://127.0.0.1:%s/%s?asked=true", p, key))
	if err != nil {
		log.Println(err)
		return "", false
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", false
	}

	v, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", false
	}

	return string(v), true
}

func askPost(port int, key, value string) error {
	c := http.Client{}
	p := strconv.Itoa(port)
	resp, err := c.Post(fmt.Sprintf("http://127.0.0.1:%s/%s?asked=true", p, key), "", bytes.NewBuffer([]byte(value)))
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}
