#!/bin/bash
echo "running server on port 55888"
/usr/bin/xvfb-run -n 0 -s '-screen 0 1024x768x24' /url2img-1.0-64bit/url2img
