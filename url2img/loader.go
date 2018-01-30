package main

import (
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/webkit"
	"github.com/therecipe/qt/widgets"
)

//go:generate qtmoc

// Object struct
type Object struct {
	core.QObject

	_ func(data string)     `signal:"load"`
	_ func(id, data string) `signal:"loadFinished"`
}

// Loader represents image loader
type Loader struct {
	*Object
	*widgets.QWidget
}

// NewLoader returns new loader
func NewLoader() *Loader {
	os.Setenv("QT_QPA_PLATFORM", "offscreen")

	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetApplicationName(name)
	app.SetApplicationVersion(version)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetAttribute(core.Qt__WA_DontShowOnScreen, true)
	widget.Show()

	l := &Loader{NewObject(nil), widget}

	l.ConnectLoad(func(data string) {
		p := NewParams()
		err := p.Unmarshal(data)
		if err == nil {
			l.LoadPage(p.Url, p.Id, p.Format, p.UA, p.Quality, p.Delay, p.Width, p.Height, p.Zoom, p.Full)
		}
	})

	l.ConnectLoadFinished(func(id, data string) {
		loaded.Store(id, data)
	})

	return l
}

// LoadPage loads page
func (l *Loader) LoadPage(url, id, format, ua string, quality, delay, width, height int, zoom float64, full bool) {
	view := webkit.NewQWebView(l.QWidget_PTR())
	view.SetAttribute(core.Qt__WA_DontShowOnScreen, true)
	view.Resize2(width, width)

	page := webkit.NewQWebPage(view.QWidget_PTR())
	view.SetPage(page)

	page.MainFrame().SetZoomFactor(zoom)
	page.MainFrame().SetScrollBarPolicy(core.Qt__Horizontal, core.Qt__ScrollBarAlwaysOff)
	page.MainFrame().SetScrollBarPolicy(core.Qt__Vertical, core.Qt__ScrollBarAlwaysOff)

	page.ConnectUserAgentForUrl(func(url *core.QUrl) string {
		if ua != "" {
			return ua
		}

		return page.UserAgentForUrlDefault(url)
	})

	l.setAttributes(page.Settings())
	l.setPath(page.Settings(), os.TempDir())

	page.ConnectLoadFinished(func(bool) {
		if delay > 0 {
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		if full {
			js := `var d=document;
				Math.max(Math.max(d.body.scrollHeight, d.documentElement.scrollHeight),
				Math.max(d.body.offsetHeight, d.documentElement.offsetHeight),
				Math.max(d.body.clientHeight, d.documentElement.clientHeight));`
			height = page.MainFrame().EvaluateJavaScript(js).ToInt(true)

			if height == 0 {
				height = defHeight
			} else if height > 32768 {
				height = 32768
			}

			page.MainFrame().EvaluateJavaScript(`window.scrollTo(0, ` + strconv.Itoa(height) + `);`)

			page.SetViewportSize(core.NewQSize2(width, height))
			view.Resize2(width, height)
		}

		image := gui.NewQImage3(width, height, gui.QImage__Format_RGB888)
		if image.IsNull() {
			l.LoadFinished(id, "ErrIsNull")
			view.DeleteLater()
			return
		}

		painter := gui.NewQPainter()
		painter.Begin(gui.NewQPaintDeviceFromPointer(image.Pointer()))
		if !painter.IsActive() {
			l.LoadFinished(id, "ErrIsActive")
			view.DeleteLater()
			return
		}

		painter.SetRenderHint(gui.QPainter__Antialiasing, true)
		painter.SetRenderHint(gui.QPainter__TextAntialiasing, true)
		painter.SetRenderHint(gui.QPainter__HighQualityAntialiasing, true)
		painter.SetRenderHint(gui.QPainter__SmoothPixmapTransform, true)
		page.MainFrame().Render(painter, gui.NewQRegion())
		painter.End()

		buff := core.NewQBuffer(view)
		buff.Open(core.QIODevice__ReadWrite)
		if !buff.IsWritable() {
			l.LoadFinished(id, "ErrIsWritable")
			view.DeleteLater()
			return
		}

		ok := image.Save2(buff, strings.ToUpper(format), quality)
		data := []byte(buff.Data().ConstData())
		if !ok {
			data = []byte("ErrSave2")
		}

		l.LoadFinished(id, hex.EncodeToString(data))

		image.DestroyQImage()

		buff.Close()
		buff.DeleteLater()

		view.DeleteLater()
	})

	view.Show()
	view.Load(core.NewQUrl3(url, core.QUrl__TolerantMode))
}

// setAttributes sets web page attributes
func (l *Loader) setAttributes(settings *webkit.QWebSettings) {
	settings.SetAttribute(webkit.QWebSettings__AutoLoadImages, true)
	settings.SetAttribute(webkit.QWebSettings__JavascriptEnabled, true)
	settings.SetAttribute(webkit.QWebSettings__JavascriptCanOpenWindows, false)
	settings.SetAttribute(webkit.QWebSettings__JavascriptCanCloseWindows, false)
	settings.SetAttribute(webkit.QWebSettings__JavascriptCanAccessClipboard, false)
	settings.SetAttribute(webkit.QWebSettings__LocalContentCanAccessFileUrls, true)
	settings.SetAttribute(webkit.QWebSettings__LocalContentCanAccessRemoteUrls, true)
	settings.SetAttribute(webkit.QWebSettings__SiteSpecificQuirksEnabled, true)
	settings.SetAttribute(webkit.QWebSettings__PrivateBrowsingEnabled, true)

	settings.SetAttribute(webkit.QWebSettings__PluginsEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__JavaEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__WebGLEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__WebAudioEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__NotificationsEnabled, false)

	settings.SetAttribute(webkit.QWebSettings__Accelerated2dCanvasEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__AcceleratedCompositingEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__TiledBackingStoreEnabled, false)

	settings.SetAttribute(webkit.QWebSettings__LocalStorageEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__LocalStorageDatabaseEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__OfflineStorageDatabaseEnabled, false)
	settings.SetAttribute(webkit.QWebSettings__OfflineWebApplicationCacheEnabled, false)
}

// setPath sets storage path
func (l *Loader) setPath(settings *webkit.QWebSettings, path string) {
	settings.SetIconDatabasePath(path)
	settings.SetLocalStoragePath(path)
	settings.SetOfflineStoragePath(path)
	settings.SetOfflineWebApplicationCachePath(path)
}

// Exec starts Qt main loop
func (l *Loader) Exec() {
	widgets.QApplication_Exec()
}
