package media

import (
	"net/url"
	"path/filepath"

	"github.com/imssyang/gweb/internal/log"
)

func Parse(urlPath string) error {
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		return err
	}

	if len(parsedURL.Scheme) == 0 || parsedURL.Scheme == "file" {
		path := parsedURL.Path
		if len(parsedURL.Host) > 0 {
			path = filepath.Join(string(filepath.Separator), parsedURL.Host, path)
		}

		switch filepath.Ext(path) {
		case ".ivf":
		case ".ogg":
		default:
			log.Zap.Errorln("UnsupportFileExt:", urlPath)
			return nil
		}
	}

	return nil
}
