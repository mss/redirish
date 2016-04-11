package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func CFail(w http.ResponseWriter, msg string) {
	Fail(w, 400, msg, nil)
}

func SFail(w http.ResponseWriter, msg string, err error) {
	Fail(w, 500, msg, err)
}

func Fail(w http.ResponseWriter, code int, msg string, err error) {
	w.Header().Add("Warning", msg)
	w.WriteHeader(code)
	if err != nil {
		msg += " (" + err.Error() + ")"
	}
	log.Print(msg)
}

func main() {
	addr := ":8080"
	argv := os.Args[1:]
	argc := len(argv)
	switch argc {
	default:
		os.Exit(1)
	case 1:
		addr = argv[0]
	case 0:
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "redirish")

		host := strings.Split(r.Host, ":")[0]
		if host == "" {
			CFail(w, "missing required Host header")
			return
		}
		if net.ParseIP(host) != nil {
			CFail(w, "missing hostname in Host header")
			return
		}
		host = strings.Trim(host, ".")
		if !strings.Contains(host, ".") {
			CFail(w, "not enough tokens in Host header")
			return
		}

		host += "."
		for i := 0; ; i++ {
			cname, _ := net.LookupCNAME(host)
			if cname == "" {
				break
			}
			host = cname
			if i == 10 {
				SFail(w, "excessive CNAME recursion", nil)
			}
		}

		host = "_http._tcp." + host
		txts, err := net.LookupTXT(host)
		if err != nil {
			SFail(w, "lookup failure", err)
			return
		}

		res := strings.SplitN(txts[rand.Intn(len(txts))], " ", 2)
		switch len(res) {
		case 2:
			base := res[1]
			w.Header().Add("Location", base+r.URL.Path)
			fallthrough
		case 1:
			code, err := strconv.Atoi(res[0])
			if err != nil {
				SFail(w, "lookup failure", err)
				return
			}
			w.WriteHeader(code)
		default:
			SFail(w, "lookup failure", nil)
		}
	})
	http.ListenAndServe(addr, nil)
}
