package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	postFormData := url.Values{}
	form := New(postFormData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when it should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postFormData := url.Values{}
	form := New(postFormData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postFormData = url.Values{}
	postFormData.Add("a", "a")
	postFormData.Add("b", "b")
	postFormData.Add("c", "c")

	form = New(postFormData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid even when required fields are present")
	}
}

func TestForm_Has(t *testing.T) {
	postFormData := url.Values{}
	postFormData.Add("a", "a")
	postFormData.Add("b", "")

	form := New(postFormData)

	if !form.Has("a") {
		t.Error("shows value does not exist on form data when it does")
	}

	if form.Has("b") {
		t.Error("shows value does exist on form data when it does not")
	}
}

func TestForm_MinLength(t *testing.T) {
	postFormData := url.Values{}
	form := New(postFormData)

	form.MinLength("x", 3)
	if form.Valid() {
		t.Error("form shows minlength for non existing field")
	}

	isErr := form.Errors.Get("x")
	if isErr == "" {
		t.Error("should have an error but did not get one")
	}

	postFormData = url.Values{}
	postFormData.Add("some_field", "some value")
	form = New(postFormData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows minlength of 100 met when data length is shorter")
	}

	postFormData = url.Values{}
	postFormData.Add("another_field", "abc123")
	form = New(postFormData)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows minlength of 1 is not met when it is")
	}

	isErr = form.Errors.Get("another_field")
	if isErr != "" {
		t.Error("should not have an error but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postFormData := url.Values{}
	form := New(postFormData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non existent field")
	}

	postFormData = url.Values{}
	postFormData.Add("email", "me@here.com")
	form = New(postFormData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postFormData = url.Values{}
	postFormData.Add("email_invalid", "me@here.")
	form = New(postFormData)

	form.IsEmail("email_invalid")
	if form.Valid() {
		t.Error("got a valid email for an invalid email address")
	}
}
