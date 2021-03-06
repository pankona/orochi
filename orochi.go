package orochi

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Orochi is a struct to launch an Orochi server
type Orochi struct {
	PortList []int

	server  *http.Server
	port    int
	kvstore map[string]string
}

// Serve starts to serve a http server with specified port
func (o *Orochi) Serve(port int) error {
	o.port = port
	o.kvstore = map[string]string{}

	mux := http.NewServeMux()
	mux.Handle("/", o)

	o.server = &http.Server{
		Addr:    ":" + strconv.Itoa(o.port),
		Handler: mux,
	}
	log.Printf("server will start on port: %d\n", o.port)
	return o.server.ListenAndServe()
}

// Port returns its port to listen
func (o *Orochi) Port() int {
	return o.port
}

// Shutdown stops served http server in graceful
func (o *Orochi) Shutdown() error {
	return o.server.Shutdown(context.Background())
}

func (o *Orochi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, err := extractKeyFromPath(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		v, ok := o.kvstore[key]
		if ok {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(v))
			if err != nil {
				log.Printf("failed to write response: %v", err)
				return
			}
			log.Printf("return %s", v)
			return
		}

		q := r.URL.Query()
		if q["asked"] != nil {
			break
		}

		// ask to other server
		log.Println("missed. ask to other server")
		for _, p := range o.PortList {
			if o.port == p {
				log.Println("skip because this is me")
				continue
			}

			v, ok := o.askGet(p, key)
			if !ok || v == "" {
				log.Printf("missed by port: %d", p)
				continue
			}

			log.Printf("hit on other server: %d", p)
			o.kvstore[key] = v
			log.Printf("stored: %s\n", v)

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(v))
			if err != nil {
				log.Printf("failed to write response: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)

	case "POST":
		v, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		o.kvstore[key] = string(v)
		log.Printf("stored: %s\n", v)

		q := r.URL.Query()
		if q["asked"] != nil {
			break
		}

		for _, p := range o.PortList {
			err := o.askPost(p, key, string(v))
			if err != nil {
				log.Printf("failed to askPost: %v", err)
			}
		}
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func extractKeyFromPath(path string) (string, error) {
	ss := strings.Split(path, "/")
	if len(ss) != 2 {
		return "", fmt.Errorf("[%s] is invalid. number of element must be 1 (but %d specified)", path, len(ss)-1)
	}
	return ss[1], nil
}

func (o *Orochi) askGet(port int, key string) (string, bool) {
	c := http.Client{}
	p := strconv.Itoa(port)

	resp, err := c.Get(fmt.Sprintf("http://127.0.0.1:%s/%s?asked=1", p, key))
	if err != nil {
		log.Println(err)
		return "", false
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed close response body: %v", err)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return "", false
	}

	v, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", false
	}

	return string(v), true
}

func (o *Orochi) askPost(port int, key, value string) error {
	c := http.Client{}
	p := strconv.Itoa(port)

	resp, err := c.Post(fmt.Sprintf("http://127.0.0.1:%s/%s?asked=1", p, key), "", bytes.NewBuffer([]byte(value)))
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed close response body: %v", err)
		}
	}()

	return nil
}
