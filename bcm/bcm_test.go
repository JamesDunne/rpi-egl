package bcm

import "testing"

func TestOpenClose(t *testing.T) {
	display, err := OpenDisplay()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("display: %dx%d", display.Width(), display.Height())
	display.Close()
}
