package multi_fileserver

import (
	"net/http"
	"os"

	"github.com/golang/glog"
)

const directoryIndex = "index.html"

type handler struct {
	dirs []string
}

func New(dirs ...string) *handler {
	h := new(handler)
	h.dirs = reverse(dirs)
	return h
}

func reverse(dirs []string) []string {
	result := make([]string, len(dirs))
	for i, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			glog.Warningf("dir %s not found", dir)
		}
		glog.V(2).Infof("setup dir %s", dir)
		result[len(result)-i-1] = dir
	}
	return result
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	name := request.URL.Path
	if name == "" || name == "/" {
		name = directoryIndex
	}
	for _, root := range h.dirs {
		glog.V(2).Infof("search file %s in directory %s", name, root)
		f, err := http.Dir(root).Open(name)
		if err != nil {
			glog.V(2).Infof("file %s not found in directory %s", name, root)
			continue
		}
		defer f.Close()
		d, err := f.Stat()
		if err != nil {
			glog.V(2).Infof("stat file %s failed: %v", err)
			return
		}
		glog.V(2).Infof("found file %s in directory %s", name, root)
		http.ServeContent(responseWriter, request, d.Name(), d.ModTime(), f)
		return
	}
	glog.V(1).Infof("file not found %s", name)
	http.NotFound(responseWriter, request)
}
