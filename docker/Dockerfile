FROM debian:jessie

# Set environment variables
ENV DEBIAN_FRONTEND noninteractive

# Install dependencies
RUN apt-get update -q && apt-get install -y xvfb xz-utils

# add url2img binary
WORKDIR /
ADD https://url2img.com/download/url2img-1.0-64bit.tar.xz /
RUN tar xf url2img-1.0-64bit.tar.xz

# add entrypoint
ADD docker-entrypoint.sh /

# start daemon
CMD ["/docker-entrypoint.sh"]
