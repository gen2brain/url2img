#!/usr/bin/env bash

mkdir -p build
go generate github.com/gen2brain/url2img/url2img

# linux/amd64
CHROOT="/usr/x86_64-pc-linux-gnu-static"
INCPATH="$CHROOT/usr/include/qt5"
PLUGPATH="$CHROOT/usr/lib/qt5/plugins"
export CC=gcc CXX=g++
export PKG_CONFIG_PATH="$CHROOT/usr/lib/pkgconfig"
export PKG_CONFIG_LIBDIR="$CHROOT/usr/lib/pkgconfig"
export LIBRARY_PATH="$CHROOT/usr/lib:$CHROOT/lib"

CGO_CXXFLAGS="-I$CHROOT/usr/include -I$INCPATH -I$INCPATH/QtCore -I$INCPATH/QtGui -I$INCPATH/QtWidgets -I$INCPATH/QtNetwork -I$INCPATH/QtWebKit -I$INCPATH/QtWebKitWidgets" \
CGO_CXXFLAGS="$CGO_CXXFLAGS -pipe -O2 -std=gnu++11 -Wall -W -D_REENTRANT -DQT_NO_DEBUG -DQT_WEBKIT_LIB -DQT_CORE_LIB -DQT_GUI_LIB  -DQT_NETWORK_LIB -DQT_WIDGETS_LIB -DQT_WEBKITWIDGETS_LIB -fPIC" \
CGO_LDFLAGS="-L$CHROOT/usr/lib -L$CHROOT/lib -L$PLUGPATH/platforms -L$PLUGPATH/generic -L/temp/qtwebkit-snapshots-master/build/qt/Release/lib" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Core -lpthread -lz -lpcre16 -ldouble-conversion -lm -ldl -lgthread-2.0 -lglib-2.0 -lrt" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5WebKit -lWTF -lWebCore -lJavaScriptCore -lWTF -lwoff2 -lsqlite3 -lbmalloc -lbrotli -lhyphen -lxslt -lxml2 -licui18n -licuuc -licudata" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Network -lQt5Widgets -lQt5WebKitWidgets -lQt5Gui -lQt5Core -lpthread -lturbojpeg -lwebp -lpng -lharfbuzz  -lssl -lcrypto -lz -lGL" \
CGO_LDFLAGS="$CGO_LDFLAGS -lqxcb -lQt5XcbQpa -lQt5PlatformSupport -lX11-xcb -lXi -lxcb-render -lxcb-render-util -lSM -lICE -lbsd -lXxf86vm -lXcursor" \
CGO_LDFLAGS="$CGO_LDFLAGS -lxcb -lxcb-image -lxcb-icccm -lxcb-sync -lxcb-xfixes -lxcb-shm -lxcb-randr -lxcb-shape -lxcb-keysyms -lxcb-xinerama -lxcb-xkb -lxcb-util -lxcb-glx" \
CGO_LDFLAGS="$CGO_LDFLAGS -lxkbcommon-x11 -lxkbcommon  -lfontconfig -lfreetype -lexpat -ldl -lXrender -lXext -lX11 -lxcb -lXau -lXdmcp" \
CGO_LDFLAGS="$CGO_LDFLAGS -lm -lQt5Gui -lturbojpeg -lpng -lharfbuzz -lz -lbz2 -lGL -lQt5Core -lpthread -lGL" \
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags 'static minimal' -o build/url2img.amd64 -v -x -ldflags "-linkmode external -s -w" github.com/gen2brain/url2img/url2img

# linux/386
CHROOT="/usr/x86_64-pc-linux-gnu-static"
INCPATH="$CHROOT/usr/include/qt5"
PLUGPATH="$CHROOT/usr/lib32/qt5/plugins"
export CC=gcc CXX=g++
export PKG_CONFIG_PATH="$CHROOT/usr/lib32/pkgconfig"
export PKG_CONFIG_LIBDIR="$CHROOT/usr/lib32/pkgconfig"
export LIBRARY_PATH="$CHROOT/usr/lib32:$CHROOT/lib32"

CGO_CXXFLAGS="-I$CHROOT/usr/include -I$INCPATH -I$INCPATH/QtCore -I$INCPATH/QtGui -I$INCPATH/QtWidgets -I$INCPATH/QtNetwork -I$INCPATH/QtWebKit -I$INCPATH/QtWebKitWidgets" \
CGO_CXXFLAGS="$CGO_CXXFLAGS -pipe -O2 -std=gnu++11 -Wall -W -D_REENTRANT -DQT_NO_DEBUG -DQT_WEBKIT_LIB -DQT_CORE_LIB -DQT_GUI_LIB -DQT_NETWORK_LIB -DQT_WIDGETS_LIB -DQT_WEBKITWIDGETS_LIB -DQT_NETWORK_LIB -fPIC" \
CGO_LDFLAGS="-L$CHROOT/usr/lib -L$CHROOT/lib -L$PLUGPATH/platforms -L$PLUGPATH/generic -L/temp/qtwebkit-snapshots/build/qt/Release/lib" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Core -lpthread -lz -lpcre16 -ldouble-conversion -lm -ldl -lgthread-2.0 -lglib-2.0 -lrt" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5WebKit -lWTF -lWebCore -lJavaScriptCore -lWTF -lwoff2 -lsqlite3 -lbmalloc -lbrotli -lhyphen -lxslt -lxml2 -licui18n -licuuc -licudata" \
CGO_LDFLAGS="$CGO_LDFLAGS -lQt5Network -lQt5WebKitWidgets -lQt5WebKit -lQt5Sql -lQt5Widgets -lQt5Gui -lQt5Core -lpthread -lturbojpeg -lwebp -lpng -lharfbuzz  -lssl -lcrypto -lz -lGL" \
CGO_LDFLAGS="$CGO_LDFLAGS -lqxcb -lQt5XcbQpa -lQt5PlatformSupport -lX11-xcb -lXi -lxcb-render -lxcb-render-util -lSM -lICE -lbsd -lXxf86vm -lXcursor" \
CGO_LDFLAGS="$CGO_LDFLAGS -lxcb -lxcb-image -lxcb-icccm -lxcb-sync -lxcb-xfixes -lxcb-shm -lxcb-randr -lxcb-shape -lxcb-keysyms -lxcb-xinerama -lxcb-xkb -lxcb-util -lxcb-glx" \
CGO_LDFLAGS="$CGO_LDFLAGS -lxkbcommon-x11 -lxkbcommon  -lfontconfig -lfreetype -lexpat -ldl -lXrender -lXext -lX11 -lxcb -lXau -lXdmcp" \
CGO_LDFLAGS="$CGO_LDFLAGS -lm -lQt5Gui -lturbojpeg -lpng -lharfbuzz -lz -lbz2 -lGL -lQt5Core -lpthread -lGL" \
CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -tags 'static minimal' -o build/url2img.386 -v -x -ldflags "-linkmode external -s -w" github.com/gen2brain/url2img/url2img
