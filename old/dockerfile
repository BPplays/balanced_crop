FROM golang:1.22

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive


# install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

# Install ImageMagick deps
RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev \
    libgif-dev libx11-dev --no-install-recommends

ENV IMAGEMAGICK_VERSION=7.1.1-34

RUN cd && \
	wget https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
	tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib

WORKDIR /go/projects/base
COPY . .
RUN go install src.techknowlogick.com/xgo@latest


# Build the Go application for Windows
# CMD GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CGO_FLAGS= go build -o /out_bin/balanced_crop.exe
CMD xgo /go/projects/base