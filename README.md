url2img
=======

url2img is HTTP server with API for capturing screenshots of websites.

Example (command line):

    $ curl -s https://api.url2img.com/?url=google.com > google.jpg

Example (web browser):

    https://api.url2img.com/?url=google.com

## API

Name    | Type      | Default   | Description
----    | ----      | -------   | -----------
url     | string    |           | Target URL (**required**), http(s):// prefix is optional
output  | string    | raw       | Output format (raw, base64, html)
format  | string    | jpg       | Image format (jpg, png, webp)
quality | int       | 75        | Image quality
delay   | int       | 0         | Delay screenshot after page is loaded (milliseconds)
width   | int       | 1600      | Viewport width
height  | int       | 1200      | Viewport height
zoom    | float     | 1.0       | Zoom factor

## Usage

    Usage of url2img:

      -bind-addr string
            Bind address (default ":55888")
      -cache-dir string
            Path to cache directory, if empty caching is disabled
      -htpasswd-file string
            Path to htpasswd file, if empty auth is disabled
      -log-file string
            Path to log file, if empty logs to stdout
      -plugins-dir string
            Path to NPAPI plugins directory, if empty searches standard directories
      -max-age int
            Cache maximum age (seconds) (default 86400)
      -read-timeout int
            Read timeout (seconds) (default 5)
      -write-timeout int
            Write timeout (seconds) (default 15)

## Auth

If server is started with -htpasswd-file it will be protected with HTTP Basic Auth. Supports MD5, SHA1 and BCrypt (https://github.com/abbot/go-http-auth).

Examples:

    $ htpasswd -c .htpasswd username
    $ url2img -htpasswd-file .htpasswd

    $ curl -I 'http://localhost:55888/?url=google.com'

    HTTP/1.1 401 Unauthorized
    Content-Type: text/plain; charset=utf-8
    Www-Authenticate: Basic realm="url2img/1.0"
    X-Content-Type-Options: nosniff
    Date: Wed, 09 Nov 2016 17:50:33 GMT
    Content-Length: 17

    $ curl -u username:password -s http://localhost:55888/?url=google.com > google.jpg
    $ curl -s http://username:password@localhost:55888/?url=google.com > google.jpg
    $ curl --netrc-file my-password-file http://localhost:55888/?url=google.com > google.jpg

## Cache

If server is started with -cache-dir it will cache HTTP responses to disk according to RFC7234 (https://github.com/lox/httpcache).
You can use -max-age to control maximum age of cache file, default is 86400 seconds (1 day).

Example:

    $ curl -I 'http://localhost:55888/?url=https://reddit.com'

    HTTP/1.1 200 OK
    Accept-Ranges: bytes
    Age: 21
    Cache-Control: public,max-age=86400
    Content-Length: 21107
    Content-Type: image/jpeg
    Last-Modified: Wed, 09 Nov 2016 18:25:26 GMT
    Proxy-Date: Wed, 09 Nov 2016 17:25:26 GMT
    Server: url2img/1.0
    Via: 1.1 httpcache
    X-Cache: HIT
    Date: Wed, 09 Nov 2016 17:25:47 GMT

If you want uncached response you can request it via Cache-Control header:

    $ curl -H 'Cache-Control: no-cache' 'http://localhost:55888/?url=https://reddit.com'

Or you can send POST request with json body, POST requests are never cached:

    $ curl -X POST -d '{"url": "https://reddit.com", "format": "png"}' http://localhost:55888

## Reload

To reload server (e.g. when .htpasswd is changed or log file is rotated) send SIGHUP signal to process.
If you use one of the provided init scripts just do a reload.

## Download

Binaries are compiled with static Qt and X libraries. They should work on all recent systems (glibc >= 2.14, gcc (stdc++) >= 4.9) without installing additional dependencies.
X server must be running.

Systemd and OpenRC init scripts are included in dist/.

 - [Linux 32bit](https://url2img.com/download/url2img-1.0-32bit.tar.xz)
 - [Linux 64bit](https://url2img.com/download/url2img-1.0-64bit.tar.xz)

## Xserver

On headless servers you can use [Xvfb](https://en.wikipedia.org/wiki/Xvfb).

If you are running Debian or Ubuntu:

    # apt-get install xvfb

If you are running Fedora:

    # dnf install xorg-x11-server-Xvfb

Start server with xvfb-run:

    # xvfb-run -n 0 -s '-screen 0 1600x1200x24' url2img

Or if you use provided systemd services:

    # systemctl start xvfb@1
    # systemctl start url2img

## Compile

Install Qt bindings (https://github.com/therecipe/qt), follow instructions there and prepare installation.

Patch bindings to add QtWebKit support:

    $ go get -d github.com/therecipe/qt
    $ cd $GOPATH/src/github.com/therecipe/qt
    $ curl -s -L https://gist.github.com/gen2brain/b680573e55c9a21b9299db21041d2465/raw/5bc6489cb114275a79c02802b74766fde1925c8d/qtwebkit.patch | patch -p1

Install Qt5WebKit.

Proceed with go get -v github.com/therecipe/qt/cmd/... and qtsetup.

Install url2img to $GOPATH/bin:

    $ go get -d github.com/gen2brain/url2img
    $ go generate github.com/gen2brain/url2img/url2img
    $ go install github.com/gen2brain/url2img/url2img

## License

url2img is free/libre software released under the terms of the GNU GPL license, see the 'COPYING' file for details.
