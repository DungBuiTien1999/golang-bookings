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
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	formData := url.Values{}
	formData.Add("a", "a")
	formData.Add("b", "b")
	formData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = formData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("show does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	formData := url.Values{}
	form := New(formData)

	isHas := form.Has("a")
	if isHas {
		t.Error("form shows valid when required field missing")
	}

	formData.Add("a", "a")

	form = New(formData)
	isHas = form.Has("a")
	if !isHas {
		t.Error("form shows invalid when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	formData1 := url.Values{}
	formData1.Add("a", "a")

	form := New(formData1)
	isValid := form.MinLength("a", 3)
	if isValid {
		t.Error("len of field only 1 but show valid")
	}

	formData2 := url.Values{}
	formData2.Add("a", "abc")

	form = New(formData2)
	isValid = form.MinLength("a", 3)
	if !isValid {
		t.Error("len of field valid but show invalid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// case invalid 1:
	formData1 := url.Values{}
	formData1.Add("email", "abc")

	form := New(formData1)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("(case 1) show valid when invalid email address")
	}

	isError := form.Errors.Get("email")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	// case invalid 2:
	formData2 := url.Values{}
	formData2.Add("email", "abc@a")

	form = New(formData2)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("(case 2) show valid when invalid email address")
	}

	// case valid:
	formData3 := url.Values{}
	formData3.Add("email", "abc@a.com")

	form = New(formData3)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("show invalid when email address is valid")
	}

	isError = form.Errors.Get("email")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}
