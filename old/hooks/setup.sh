#!/usr/bin/bash

apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev libgif-dev libx11-dev --no-install-recommends

wget https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
cd ImageMagick* && \
./configure \
	--without-magick-plus-plus \
	--without-perl \
	--disable-openmp \
	--with-gvc=no \
	--disable-docs && 