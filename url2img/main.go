package main

import (
	"flag"

	"github.com/orcaman/concurrent-map"
)

const (
	Name    = "url2img"
	Version = "1.0"
)

var (
	loader *Loader
	loaded cmap.ConcurrentMap
)

func main() {
	loader = NewLoader()
	loaded = cmap.New()

	server := NewServer()
	flag.StringVar(&server.Bind, "bind-addr", ":55888", "Bind address")
	flag.StringVar(&server.LogFilePath, "log-file", "", "Path to log file, if empty logs to stdout")
	flag.StringVar(&server.CacheDir, "cache-dir", "", "Path to cache directory, if empty caching is disabled")
	flag.StringVar(&server.Htpasswd, "htpasswd-file", "", "Path to htpasswd file, if empty auth is disabled")
	flag.StringVar(&server.PluginsDir, "plugins-dir", "", "Path to NPAPI plugins directory, if empty search standard directories")
	flag.IntVar(&server.MaxAge, "max-age", 86400, "Cache maximum age (seconds)")
	flag.IntVar(&server.ReadTimeout, "read-timeout", 5, "Read timeout (seconds)")
	flag.IntVar(&server.WriteTimeout, "write-timeout", 15, "Write timeout (seconds)")
	flag.Parse()

	server.ListenAndServe()
	defer server.LogFile.Close()

	loader.Exec()
}
