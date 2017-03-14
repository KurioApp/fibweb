package fibweb_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/KurioApp/fibweb/internal/mocks"
	"github.com/uudashr/fibweb"
)

func TestNumbers_validLimit(t *testing.T) {
	fibonacciService := new(mocks.FibonacciService)
	handler := fibweb.NewHTTPHandler(fibonacciService)

	limit := 7
	out := []int{0, 1, 1, 2, 3, 5, 8}

	fibonacciService.On("Seq", limit).Return(out, nil)

	req, err := http.NewRequest(http.MethodGet, "/api/fibonacci/numbers", nil)
	if err != nil {
		t.Error("err:", err)
		t.SkipNow()
	}

	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Error("got:", got, "want:", want)
		t.SkipNow()
	}

	var result []int
	err = json.NewDecoder(rec.Body).Decode(&result)
	if err != nil {
		t.Error("err:", err)
		t.SkipNow()
	}

	if got, want := len(out), len(result); got != want {
		t.Error("len(got):", got, "len(want):", want)
		t.SkipNow()
	}

	for i, count := 0, len(result); i < count; i++ {
		if got, want := result[i], out[i]; got != want {
			t.Error("result[i]:", got, "out[i]:", want, "i:", i)
		}
	}
}

func TestNumbers_noLimit(t *testing.T) {
	fibonacciService := new(mocks.FibonacciService)
	handler := fibweb.NewHTTPHandler(fibonacciService)

	out := []int{0, 1, 1, 2, 3, 5}

	fibonacciService.On("Seq", 5).Return(out, nil)

	req, err := http.NewRequest(http.MethodGet, "/api/fibonacci/numbers", nil)
	if err != nil {
		t.Error("err:", err)
		t.SkipNow()
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Error("got:", got, "want:", want)
		t.SkipNow()
	}

	var result []int
	err = json.NewDecoder(rec.Body).Decode(&result)
	if err != nil {
		t.Error("err:", err)
		t.SkipNow()
	}

	if got, want := len(out), len(result); got != want {
		t.Error("len(got):", got, "len(want):", want)
		t.SkipNow()
	}

	for i, count := 0, len(result); i < count; i++ {
		if got, want := result[i], out[i]; got != want {
			t.Error("result[i]:", got, "out[i]:", want, "i:", i)
		}
	}

	fibonacciService.AssertExpectations(t)
}

func TestNumbers_emptyLimit(t *testing.T) {
	fibonacciService := new(mocks.FibonacciService)
	handler := fibweb.NewHTTPHandler(fibonacciService)

	out := []int{0, 1, 1, 2, 3, 5}

	fibonacciService.On("Seq", 5).Return(out, nil)

	req, err := http.NewRequest(http.MethodGet, "/api/fibonacci/numbers", nil)
	if err != nil {
		t.Error("err:", err)
		t.SkipNow()
	}

	q := req.URL.Query()
	q.Add("limit", "")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Error("got:", got, "want:", want)
		t.SkipNow()
	}

	var result []int
	err = json.NewDecoder(rec.Body).Decode(&result)
	if err != nil {
		t.Error("err:", err)
		t.SkipNow()
	}

	if got, want := len(out), len(result); got != want {
		t.Error("len(got):", got, "len(want):", want)
		t.SkipNow()
	}

	for i, count := 0, len(result); i < count; i++ {
		if got, want := result[i], out[i]; got != want {
			t.Error("result[i]:", got, "out[i]:", want, "i:", i)
		}
	}

	fibonacciService.AssertExpectations(t)
}

func TestNumbers_badLimit(t *testing.T) {
	fibonacciService := new(mocks.FibonacciService)
	handler := fibweb.NewHTTPHandler(fibonacciService)

	limits := []string{"asdf", "-1", "29jh"}
	for _, limit := range limits {
		req, err := http.NewRequest(http.MethodGet, "/api/fibonacci/numbers", nil)
		if err != nil {
			t.Error("err:", err, "limit:", limit)
			t.SkipNow()
		}

		q := req.URL.Query()
		q.Add("limit", limit)
		req.URL.RawQuery = q.Encode()

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if got, want := rec.Code, http.StatusBadRequest; got != want {
			t.Error("got:", got, "want:", want, "limit:", limit)
		}
	}
}

func TestNumbers_serviceError(t *testing.T) {
	fibonacciService := new(mocks.FibonacciService)
	handler := fibweb.NewHTTPHandler(fibonacciService)

	limit := 10

	fibonacciService.On("Seq", limit).Return(nil, errors.New("Some error"))

	req, err := http.NewRequest(http.MethodGet, "/api/fibonacci/numbers", nil)
	if err != nil {
		t.Error("err:", err, "limit:", limit)
		t.SkipNow()
	}

	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusBadRequest; got != want {
		t.Error("got:", got, "want:", want, "limit:", limit)
	}
}
