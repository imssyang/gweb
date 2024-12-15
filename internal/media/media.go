package media

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/imssyang/gweb/internal/log"
	"github.com/imssyang/gweb/internal/media/ivf"
	"github.com/imssyang/gweb/internal/media/ogg"
)

const (
	MuxTypeIVF = "ivf"
	MuxTypeOGG = "ogg"
)

type TrackParser interface {
	MimeType() string
	NextFrame() ([]byte, any, error)
}

type Media struct {
	URL     string
	MuxType string
	Parsers []TrackParser
}

func NewMedia(uRL string) (*Media, error) {
	m := &Media{
		URL: uRL,
	}

	parsedURL, err := url.Parse(uRL)
	if err != nil {
		return nil, err
	}

	if len(parsedURL.Scheme) == 0 || parsedURL.Scheme == "file" {
		path := parsedURL.Path
		if len(parsedURL.Host) > 0 {
			path = filepath.Join(string(filepath.Separator), parsedURL.Host, path)
		}

		switch strings.ToLower(filepath.Ext(path)) {
		case ".ivf":
			ivf, err := ivf.NewIVFParser(uRL)
			if err != nil {
				return nil, err
			}

			m.MuxType = MuxTypeIVF
			m.Parsers = append(m.Parsers, *ivf)
		case ".ogg":
			ogg, err := ogg.NewOggParser(uRL)
			if err != nil {
				return nil, err
			}

			m.MuxType = MuxTypeOGG
			m.Parsers = append(m.Parsers, *ogg)
		default:
			log.Zap.Errorln("UnsupportFileExt:", uRL)
			return nil, fmt.Errorf("InvalidExt: %s", uRL)
		}
	}

	return m, nil
}
