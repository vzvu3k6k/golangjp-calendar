package calendar

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRun(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, loadTestdata(t, "blog_full.xml"))
		}),
	)
	defer ts.Close()

	var out bytes.Buffer
	err := Run(&out, []string{ts.URL})
	// os.WriteFile("testdata/blog_full.ical", out.Bytes(), 0644)
	assert.NilError(t, err)
	assert.Equal(t, out.String(), loadTestdata(t, "blog_full.ical"))
}
