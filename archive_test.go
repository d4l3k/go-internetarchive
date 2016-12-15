package archive

import "testing"

func TestSnapshot(t *testing.T) {
	resp, err := Snapshot("http://fn.lc")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Resp: %+v", resp)
}
