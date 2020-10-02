#!/usr/bin/env bash

mkdir -p build

CHROOT="/usr/x86_64-pc-linux-musl"
INCPATH="$CHROOT/usr/include/qt5"
PLUGPATH="$CHROOT/usr/plugins"

CC="x86_64-pc-linux-musl-gcc" \
CXX="x86_64-pc-linux-musl-g++" \
LIBRARY_PATH="$CHROOT/usr/lib:$CHROOT/lib" \
PKG_CONFIG_PATH="$CHROOT/usr/lib/pkgconfig" \
PKG_CONFIG_LIBDIR="$CHROOT/usr/lib/pkgconfig" \
CGO_CXXFLAGS="-I$CHROOT/usr/include -I$INCPATH -I$INCPATH/QtCore -I$INCPATH/QtGui -I$INCPATH/QtWidgets -I$INCPATH/QtNetwork -I$INCPATH/QtWebKit -I$INCPATH/QtWebKitWidgets" \
CGO_CXXFLAGS="$CGO_CXXFLAGS -pipe -O2 -std=gnu++11 -Wall -W -D_REENTRANT -DQT_NO_DEBUG -DQT_WEBKIT_LIB -DQT_CORE_LIB -DQT_GUI_LIB -DQT_NETWORK_LIB -DQT_WIDGETS_LIB" \
CGO_CXXFLAGS="$CGO_CXXFLAGS -DQT_WEBKITWIDGETS_LIB -fPIC -Wno-unused-parameter -Wno-unused-variable -Wno-return-type" \
CGO_LDFLAGS="-L$CHROOT/usr/lib -L$CHROOT/lib -L$PLUGPATH/platforms -L$PLUGPATH/imageformats -L/temp/t/qtwebkit-5.212.0-alpha4/build/lib" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Core -lpthread -lz -lqtpcre2 -ldouble-conversion -lm -ldl -lrt" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5WebKit -lQt5WebKitWidgets -lWTF -lWebCore -lWebCoreTestSupport -lJavaScriptCore -lWTF -lwoff2 -lsqlite3 -lbmalloc -lbrotli -lhyphen -lxslt -lxml2 -licui18n -licuuc -licudata" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Network -lQt5WebKitWidgets -lQt5Widgets -lQt5Gui -lQt5Core -lQt5PrintSupport -lQt5Sql -lpthread -ljpeg -lwebp -lpng -lssl -lcrypto -lz" \
CGO_LDFLAGS="$CGO_LDFLAGS -lqoffscreen -lqgif -lqjpeg -lQt5ThemeSupport -lQt5FontDatabaseSupport -lQt5ServiceSupport -lQt5EventDispatcherSupport" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5WebKit -lQt5WebKitWidgets -lWTF -lWebCore -lWebCoreTestSupport -lJavaScriptCore -lWTF -lwoff2 -lsqlite3 -lbmalloc -lbrotli -lhyphen -lxslt -lxml2 -licui18n -licuuc -licudata" \
CGO_LDFLAGS="$CGO_LDFLAGS -lm -lQt5Gui -lfontconfig -lfreetype -luuid -ljpeg -lpng -lqtharfbuzz -lz -lQt5Core -lpthread" \
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -tags 'static minimal' -o build/url2img -v -x -ldflags "-linkmode external -s -w '-extldflags=-static'" github.com/gen2brain/url2img/cmd/url2img
