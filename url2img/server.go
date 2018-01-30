package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/abbot/go-http-auth"
	"github.com/lox/httpcache"
)

// Server represents HTTP server
type Server struct {
	Bind         string
	Auth         *auth.BasicAuth
	LogFile      *os.File
	LogFilePath  string
	CacheDir     string
	Htpasswd     string
	MaxAge       int
	ReadTimeout  int
	WriteTimeout int
}

// NewServer returns new Server
func NewServer() *Server {
	return &Server{}
}

// ServeHTTP handles requests on incoming connections
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := NewParams()

	switch r.Method {
	case "GET", "HEAD":
		err := p.FormValues(r)
		if err != nil {
			msg := fmt.Sprintf("400 Bad Request (%s)", err.Error())
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	case "POST":
		err := p.BodyValues(r)
		if err != nil {
			msg := fmt.Sprintf("400 Bad Request (%s)", err.Error())
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	default:
		msg := fmt.Sprintf("405 Method Not Allowed (%s)", r.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	d, err := p.Marshal()
	if err != nil {
		msg := fmt.Sprintf("500 Internal Server Error (%s)", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	loader.Load(d)

	if !s.wait(p.Id) {
		msg := fmt.Sprintf("408 Request Timeout (after %d seconds)", s.ReadTimeout+s.WriteTimeout)
		http.Error(w, msg, http.StatusRequestTimeout)
		return
	}

	str, _ := loaded.Load(p.Id)
	loaded.Delete(p.Id)

	data, err := hex.DecodeString(str.(string))
	if err != nil {
		msg := fmt.Sprintf("500 Internal Server Error (%s)", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(string(data), "Err") {
		msg := fmt.Sprintf("500 Internal Server Error (%s)", string(data))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if s.CacheDir != "" {
		w.Header().Set("Cache-Control", fmt.Sprintf("public,max-age=%d", s.MaxAge))
		w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	}

	switch p.Output {
	case "raw":
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.%s\"", p.Url, p.Format))
		w.Write(data)
	case "base64":
		b64 := base64.StdEncoding.EncodeToString(data)
		w.Write([]byte(b64))
	case "html":
		html := "<!DOCTYPE html><html><body><img src=\"data:image/%s;base64,%s\" download\"%s\"/></body></html>"
		html = fmt.Sprintf(html, p.Format, base64.StdEncoding.EncodeToString(data), p.Url+"."+p.Format)
		w.Write([]byte(html))
	}
}

// ListenAndServe listens on the TCP address and serves requests
func (s *Server) ListenAndServe() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User-agent: *\nDisallow: /"))
	})

	s.open()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-c
			s.open()
		}
	}()

	if s.CacheDir != "" {
		cache, err := httpcache.NewDiskCache(s.CacheDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(4)
		}

		handler := httpcache.NewHandler(cache, s)
		http.Handle("/", newHandler(handler, s.LogFile, s.Auth))
	} else {
		http.Handle("/", newHandler(s, s.LogFile, s.Auth))
	}

	srv := &http.Server{
		ReadTimeout:  time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.WriteTimeout) * time.Second,
	}

	listener, err := net.Listen("tcp4", s.Bind)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(5)
	}

	go srv.Serve(listener)
}

// wait waits for url to load, timeouts after ReadTimeout+WriteTimeout
func (s *Server) wait(id string) bool {
	end := make(chan bool, 1)
	timeout := time.After(time.Duration(s.ReadTimeout+s.WriteTimeout) * time.Second)

	for {
		_, ok := loaded.Load(id)
		select {
		case <-end:
			return true
		case <-timeout:
			return false
		default:
			if ok {
				end <- true
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// open opens log and htpasswd file
func (s *Server) open() {
	if s.Htpasswd != "" {
		if _, err := os.Stat(s.Htpasswd); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		s.Auth = auth.NewBasicAuthenticator(fmt.Sprintf("%s/%s", name, version), auth.HtpasswdFileProvider(s.Htpasswd))
	}

	if s.LogFile != nil {
		s.LogFile.Close()
	}

	s.LogFile = os.Stdout
	if s.LogFilePath != "" {
		var err error
		if _, err = os.Stat(s.LogFilePath); err == nil {
			s.LogFile, err = os.Open(s.LogFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(2)
			}
		} else {
			s.LogFile, err = os.Create(s.LogFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(3)
			}
		}
	}
}

// newHandler wraps handler and logs requests (apache common log format + elapsed time and cache status)
func newHandler(handler http.Handler, file *os.File, auth *auth.BasicAuth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Recovered: %v\n", r)
			}
		}()

		startTime := time.Now()

		ip := r.RemoteAddr
		if c := strings.LastIndex(ip, ":"); c != -1 {
			ip = ip[:c]
		}

		pattern := "%s - %s [%s] \"%s %d %d\" %f %s\n"
		request := fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)
		timeFormat := "02/Jan/2006:15:04:05 -0700"

		rw := NewResponseWriter(w)
		rw.Header().Set("Server", fmt.Sprintf("%s/%s", name, version))
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		userid := "-"
		if auth != nil {
			rw.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", auth.Realm))

			userid = auth.CheckAuth(r)
			if userid == "" {
				http.Error(rw, "401 Unauthorized", http.StatusUnauthorized)

				endTime := time.Now()
				elapsedTime := endTime.Sub(startTime)
				formattedTime := endTime.Format(timeFormat)

				fmt.Fprintf(file, pattern, ip, userid, formattedTime, request, rw.Status(), rw.Size(), elapsedTime.Seconds(), "")
				return
			}
		}

		handler.ServeHTTP(rw, r)

		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		formattedTime := endTime.Format(timeFormat)

		cacheStatus := rw.Header().Get(httpcache.CacheHeader)

		fmt.Fprintf(file, pattern, ip, userid, formattedTime, request, rw.Status(), rw.Size(), elapsedTime.Seconds(), cacheStatus)
	})
}
