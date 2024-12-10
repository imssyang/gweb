package ogg

import (
	"os"

	"github.com/pion/webrtc/v4/pkg/media/oggreader"
)

type OggParser struct {
	Path   string
	Header *oggreader.OggHeader
	Reader *oggreader.OggReader
}

func NewOggParser(path string) (*OggParser, error) {
	ivf := &OggParser{
		Path: path,
	}

	file, err := os.Open(ivf.Path)
	if err != nil {
		return nil, err
	}

	reader, header, err := oggreader.NewWith(file)
	if err != nil {
		return nil, err
	}

	ivf.Header = header
	ivf.Reader = reader
	return ivf, nil
}
