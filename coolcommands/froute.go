package coolcommands

//go:generate genfor-interp-a $GOFILE

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/KernelDeimos/anything-gos/interp_a"

	"github.com/gin-gonic/gin"
)

func InstallFroute(ii interp_a.HybridEvaluator) {
	i_froute := interp_a.InterpreterFactoryA{}.MakeEmpty()

	i_froute.AddOperation("host", CmdFrouteHost)
	i_froute.AddOperation("proxy", CmdFrouteProxy)

	ii.AddOperation("froute", i_froute.OpEvaluate)
}

func CmdFrouteHost(args []interface{}) ([]interface{}, error) {
	//::gen verify-args froute-host addr string dir string
	if len(args) < 2 {
		return nil, errors.New("froute-host requires at least 2 arguments")
	}

	var addr string
	var dir string
	{
		var ok bool
		addr, ok = args[0].(string)
		if !ok {
			return nil, errors.New("froute-host: argument 0: addr; must be type string")
		}
		dir, ok = args[1].(string)
		if !ok {
			return nil, errors.New("froute-host: argument 1: dir; must be type string")
		}
	}
	//::end

	router := gin.Default()
	router.Static("/", dir)
	err := router.Run(addr)
	if err != nil {
		return nil, err
	}

	return []interface{}{}, nil
}

func CmdFrouteProxy(args []interface{}) ([]interface{}, error) {
	//::gen verify-args froute-host local string remote string
	if len(args) < 2 {
		return nil, errors.New("froute-host requires at least 2 arguments")
	}

	var local string
	var remote string
	{
		var ok bool
		local, ok = args[0].(string)
		if !ok {
			return nil, errors.New("froute-host: argument 0: local; must be type string")
		}
		remote, ok = args[1].(string)
		if !ok {
			return nil, errors.New("froute-host: argument 1: remote; must be type string")
		}
	}
	//::end

	target, err := url.Parse(remote)
	if err != nil {
		return []interface{}{}, err
	}

	targetQuery := target.RawQuery

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		// REWRITE THE HEADERS
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", target.Host)
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(local, nil)
	if err != nil {
		return []interface{}{}, err
	}

	return []interface{}{}, nil
}

/*
# Concept:

# Host static files
froute host localhost:8004 .

# Relay hosted files
froute proxy localhost:8005 http://10.8.0.4:8004
*/

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
