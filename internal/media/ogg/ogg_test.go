package ogg_test

import (
	"testing"

	"github.com/imssyang/gweb/internal/media/ogg"
)

func TestOggFile(t *testing.T) {
	path := "/opt/app/gweb/tests/play-from-disk/output.ogg"
	ogg, err := ogg.NewOggParser(path)
	if err != nil {
		t.Fatalf("path: %s err: %v", path, err)
	}
	t.Logf("%+v", ogg.Header)
}
