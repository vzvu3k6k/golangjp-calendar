package calendar

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
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
	assert.NilError(t, err)
	golden.Assert(t, out.String(), "blog_full.ical")
}
