package bcm

import "testing"

func TestOpenClose(t *testing.T) {
	display := OpenDisplay()
	if display == nil {
		t.Fatal("display = nil")
	}
	t.Logf("display: %dx%d", display.Width(), display.Height())
	display.Close()
}
