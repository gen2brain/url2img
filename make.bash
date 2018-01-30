#!/usr/bin/env bash

mkdir -p build

CHROOT="/usr/x86_64-pc-linux-gnu-static"
INCPATH="$CHROOT/usr/include/qt5"
PLUGPATH="$CHROOT/usr/lib/qt5/plugins"
export CC=gcc CXX=g++
export PKG_CONFIG_PATH="$CHROOT/usr/lib/pkgconfig"
export PKG_CONFIG_LIBDIR="$CHROOT/usr/lib/pkgconfig"
export LIBRARY_PATH="$CHROOT/usr/lib:$CHROOT/lib"

CGO_CXXFLAGS="-I$CHROOT/usr/include -I$INCPATH -I$INCPATH/QtCore -I$INCPATH/QtGui -I$INCPATH/QtWidgets -I$INCPATH/QtNetwork -I$INCPATH/QtWebKit -I$INCPATH/QtWebKitWidgets" \
CGO_CXXFLAGS="$CGO_CXXFLAGS -pipe -O2 -std=gnu++11 -Wall -W -D_REENTRANT -DQT_NO_DEBUG -DQT_WEBKIT_LIB -DQT_CORE_LIB -DQT_GUI_LIB -DQT_NETWORK_LIB -DQT_WIDGETS_LIB" \
CGO_CXXFLAGS="$CGO_CXXFLAGS -DQT_WEBKITWIDGETS_LIB -fPIC -Wno-unused-parameter -Wno-unused-variable -Wno-return-type" \
CGO_LDFLAGS="-L$CHROOT/usr/lib -L$CHROOT/lib -L$PLUGPATH/platforms -L$PLUGPATH/imageformats -L/temp/qtwebkit/build/lib" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Core -lQt5DBus -lpthread -lz -lpcre2-16 -ldouble-conversion -lm -ldl -lgthread-2.0 -lglib-2.0 -lrt" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5WebKit -lQt5WebKitWidgets -lWTF -lWebCore -lJavaScriptCore -lANGLESupport -lWTF -lwoff2 -lsqlite3 -lbmalloc -lbrotli -lhyphen -lxslt -lxml2 -licui18n -licuuc -licudata" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Network -lQt5WebKitWidgets -lQt5Widgets -lQt5Gui -lQt5Core -lQt5Sql -lpthread -ljpeg -lwebp -lpng -lharfbuzz -lssl -lcrypto -lz -lGL" \
CGO_LDFLAGS="$CGO_LDFLAGS -lqoffscreen -lqgif -lqjpeg -lqwebp -lQt5ThemeSupport -lQt5GlxSupport -lQt5FontDatabaseSupport -lQt5ServiceSupport -lQt5EventDispatcherSupport -lQt5DBus -ldbus-1" \
CGO_LDFLAGS="$CGO_LDFLAGS -lXi -lSM -lICE -lbsd -lXxf86vm -lXcursor -lxkbcommon-x11 -lxkbcommon -lfontconfig -lfreetype -lexpat -ldl -lXrender -lXext -lX11 -lxcb -lXau -lXdmcp" \
CGO_LDFLAGS="$CGO_LDFLAGS -lm -lQt5Gui -ljpeg -lwebp -lwebpdemux -lpng -lharfbuzz -lz -lbz2 -lGL -lQt5Core -lpthread -lGL -lX11 -lXext -lnettle" \
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags 'static minimal' -o build/url2img.amd64 -v -x -ldflags "-linkmode external -s -w" github.com/gen2brain/url2img/url2img
