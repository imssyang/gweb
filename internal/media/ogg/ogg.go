package ogg

import (
	"os"

	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/oggreader"
)

type OggParser struct {
	Path   string
	Header *oggreader.OggHeader
	Reader *oggreader.OggReader
}

func NewOggParser(path string) (*OggParser, error) {
	ogg := &OggParser{
		Path: path,
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader, header, err := oggreader.NewWith(file)
	if err != nil {
		return nil, err
	}

	ogg.Header = header
	ogg.Reader = reader
	return ogg, nil
}

func (p OggParser) MimeType() string {
	return webrtc.MimeTypeOpus
}

func (p OggParser) NextFrame() ([]byte, any, error) {
	return p.Reader.ParseNextPage()
}
