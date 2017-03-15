package httpfib_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/uudashr/fibweb/httpfib"
)

func TestSeq_limitAdded(t *testing.T) {
	limit := 10

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		limitS := q.Get("limit")

		limitI, err := strconv.Atoi(limitS)
		if err != nil {
			t.Error("err:", err)
			t.SkipNow()
		}

		if got, want := limitI, limit; got != want {
			t.Error("got:", got, "want:", want)
		}
	}))
	defer ts.Close()

	svc := httpfib.NewFibonacciService(ts.URL)
	svc.Seq(10)
}

func TestSeq_decodedResult(t *testing.T) {
	result := []int{1, 2, 3}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(result)
	}))
	defer ts.Close()

	svc := httpfib.NewFibonacciService(ts.URL)
	out, err := svc.Seq(10)
	if err != nil {
		t.Error("err:", err)
		t.SkipNow()
	}

	if got, want := len(out), len(result); got != want {
		t.Error("got:", got, "want:", want)
		t.SkipNow()
	}

	for i, count := 0, len(out); i < count; i++ {
		if got, want := out[i], result[i]; got != want {
			t.Error("out[i]:", got, "result[i]:", want, "i:", i)
			t.SkipNow()
		}
	}
}
