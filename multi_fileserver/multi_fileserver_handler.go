package multi_fileserver

import (
	"net/http"
	"os"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const DIRECTORY_INDEX = "index.html"

type multiFileserverHandler struct {
	dirs []string
}

func NewMultiFileserverHandler(dirs ...string) *multiFileserverHandler {
	h := new(multiFileserverHandler)
	h.dirs = reverse(dirs)
	return h
}

func reverse(dirs []string) []string {
	result := make([]string, len(dirs))
	for i, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			logger.Warnf("dir %s not found", dir)
		}
		logger.Debugf("setup dir %s", dir)
		result[len(result)-i-1] = dir
	}
	return result
}

func (h *multiFileserverHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	name := request.URL.Path
	if name == "" || name == "/" {
		name = DIRECTORY_INDEX
	}
	for _, root := range h.dirs {
		logger.Debugf("search file %s in directory %s", name, root)
		f, err := http.Dir(root).Open(name)
		if err != nil {
			logger.Debugf("file %s not found in directory %s", name, root)
			continue
		}
		defer f.Close()
		d, err := f.Stat()
		if err != nil {
			logger.Debugf("stat file %s failed: %v", err)
			return
		}
		logger.Debugf("found file %s in directory %s", name, root)
		http.ServeContent(responseWriter, request, d.Name(), d.ModTime(), f)
		return
	}
	logger.Infof("file not found %s", name)
	http.NotFound(responseWriter, request)
}
