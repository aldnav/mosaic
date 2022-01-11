package mosaic

import "testing"

func TestCreateMosaic(t *testing.T) {
	got := CreateMosaic(
		"_testdata/test.png",
		"_testdata/libimages",
		[]int{32, 32},
		"_testdata/",
	)
	want := "_testdata/test_out.png"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
