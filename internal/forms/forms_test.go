package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("Form shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "b")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("Shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	form := New(postedData)
	form.MinLength("a", 2)
	if form.Valid() {
		t.Error("Form shows min length for field a when data is shorter")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "john")

	form = New(postedData)

	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("Shows min length of 3 is met when data is longer")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("Should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "psanodiya@gmail.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Got an invalid email when it should have been valid")
	}

	postedData = url.Values{}
	postedData.Add("email", "x")
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Got a valid email for invalid email")
	}
}
