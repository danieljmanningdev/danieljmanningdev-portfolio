package http

import (
	"net/http"
	"os"
	"strings"
)

// StaticFiles serves public files without directory listings or dotfiles.
// The directory must contain only trusted, deployable public assets.
func StaticFiles(directory string) http.Handler {
	return http.FileServer(publicAssetFS{root: http.Dir(directory)})
}

type publicAssetFS struct{ root http.FileSystem }

func (f publicAssetFS) Open(name string) (http.File, error) {
	for _, segment := range strings.Split(name, "/") {
		if strings.HasPrefix(segment, ".") {
			return nil, os.ErrNotExist
		}
	}
	file, err := f.root.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if !info.Mode().IsRegular() {
		_ = file.Close()
		return nil, os.ErrNotExist
	}
	return file, nil
}
