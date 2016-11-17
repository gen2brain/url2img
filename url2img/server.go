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

// HTTP server
type Server struct {
	Bind         string
	Auth         *auth.BasicAuth
	LogFile      *os.File
	LogFilePath  string
	CacheDir     string
	Htpasswd     string
	PluginsDir   string
	MaxAge       int
	ReadTimeout  int
	WriteTimeout int
}

// Returns new Server
func NewServer() *Server {
	return &Server{}
}

// Handles requests on incoming connections
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

	// Load URL
	loader.Load(d)

	if !s.wait(p.Id) {
		msg := fmt.Sprintf("408 Request Timeout (after %d seconds)", s.ReadTimeout+s.WriteTimeout)
		http.Error(w, msg, http.StatusRequestTimeout)
		return
	}

	str, _ := loaded.Get(p.Id)
	data, err := hex.DecodeString(str.(string))
	if err != nil {
		msg := fmt.Sprintf("500 Internal Server Error (%s)", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	loaded.Remove(p.Id)

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

// Waits for url to load, timeouts after ReadTimeout+WriteTimeout
func (s *Server) wait(id string) bool {
	end := make(chan bool, 1)
	timeout := time.After(time.Duration(s.ReadTimeout+s.WriteTimeout) * time.Second)

	for {
		_, ok := loaded.Get(id)
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

	return false
}

// Opens log and htpasswd file
func (s *Server) Open() {
	if s.Htpasswd != "" {
		if _, err := os.Stat(s.Htpasswd); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		s.Auth = auth.NewBasicAuthenticator(fmt.Sprintf("%s/%s", Name, Version), auth.HtpasswdFileProvider(s.Htpasswd))
	}

	if s.LogFile != nil {
		s.LogFile.Close()
	}

	s.LogFile = os.Stderr
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

// Listens on the TCP address and serves requests
func (s *Server) ListenAndServe() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User-agent: *\nDisallow: /"))
	})

	if s.PluginsDir != "" {
		// Set NPAPI plugins path
		if st, err := os.Stat(s.PluginsDir); err == nil && st.IsDir() {
			os.Setenv("QTWEBKIT_PLUGIN_PATH", s.PluginsDir)
		}
	}

	s.Open()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-c
			s.Open()
		}
	}()

	if s.CacheDir != "" {
		cache, err := httpcache.NewDiskCache(s.CacheDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(4)
		}

		// Cache handler
		handler := httpcache.NewHandler(cache, s)

		http.Handle("/", NewHandler(handler, s.LogFile, s.Auth))
	} else {
		http.Handle("/", NewHandler(s, s.LogFile, s.Auth))
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

// Wraps handler and logs requests (apache common log format + elapsed time and cache status)
func NewHandler(handler http.Handler, file *os.File, auth *auth.BasicAuth) http.Handler {
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
		rw.Header().Set("Server", fmt.Sprintf("%s/%s", Name, Version))
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

// http.ResponseWriter with size and status
type responseWriter struct {
	http.ResponseWriter

	size   int
	status int
}

// Returns new responseWriter
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, 0, 0}
}

// Writes the data to the connection as part of an HTTP reply
func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// Sends an HTTP response header with status code
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Returns response size
func (w *responseWriter) Size() int {
	return w.size
}

// Returns response status code
func (w *responseWriter) Status() int {
	return w.status
}
