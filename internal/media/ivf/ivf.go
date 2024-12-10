package ivf

import (
	"os"

	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/ivfreader"
)

type IVFParser struct {
	Path   string
	Header *ivfreader.IVFFileHeader
	Reader *ivfreader.IVFReader
}

func NewIVFParser(path string) (*IVFParser, error) {
	ivf := &IVFParser{
		Path: path,
	}

	file, err := os.Open(ivf.Path)
	if err != nil {
		return nil, err
	}

	reader, header, err := ivfreader.NewWith(file)
	if err != nil {
		return nil, err
	}

	ivf.Header = header
	ivf.Reader = reader
	return ivf, nil
}

func (p *IVFParser) WebRTCMimeType() string {
	switch p.Header.FourCC {
	case "AV01":
		return webrtc.MimeTypeAV1
	case "VP90":
		return webrtc.MimeTypeVP9
	case "VP80":
		return webrtc.MimeTypeVP8
	default:
		return ""
	}
}
