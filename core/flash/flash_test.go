package flash

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCookieStoreDefaults(t *testing.T) {
	s := NewCookieStore()

	if s.cookieName != "inertia_flash" {
		t.Errorf("expected cookie name %q, got %q", "inertia_flash", s.cookieName)
	}

	if s.path != "/" {
		t.Errorf("expected path %q, got %q", "/", s.path)
	}

	if !s.httpOnly {
		t.Error("expected httpOnly to be true")
	}

	if s.secure {
		t.Error("expected secure to be false by default")
	}
}

func TestNewCookieStoreWithOptions(t *testing.T) {
	s := NewCookieStore(
		WithCookieName("my_flash"),
		WithPath("/app"),
		WithSecure(true),
		WithHTTPOnly(false),
		WithSameSite(http.SameSiteStrictMode),
	)

	if s.cookieName != "my_flash" {
		t.Errorf("expected cookie name %q, got %q", "my_flash", s.cookieName)
	}

	if s.path != "/app" {
		t.Errorf("expected path %q, got %q", "/app", s.path)
	}

	if !s.secure {
		t.Error("expected secure to be true")
	}

	if s.httpOnly {
		t.Error("expected httpOnly to be false")
	}

	if s.sameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSiteStrictMode, got %v", s.sameSite)
	}
}

func TestSetAndConsume(t *testing.T) {
	s := NewCookieStore(WithCookieName("test_flash"))
	msg := Message{Kind: "success", Title: "Done", Message: "Contact created."}

	rec := httptest.NewRecorder()
	s.Set(rec, msg)

	cookies := rec.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected a cookie to be set")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookies[0])

	rec2 := httptest.NewRecorder()
	got := s.Consume(rec2, req)

	if got == nil {
		t.Fatal("expected a flash message, got nil")
	}

	if got.Kind != msg.Kind {
		t.Errorf("expected Kind %q, got %q", msg.Kind, got.Kind)
	}

	if got.Title != msg.Title {
		t.Errorf("expected Title %q, got %q", msg.Title, got.Title)
	}

	if got.Message != msg.Message {
		t.Errorf("expected Message %q, got %q", msg.Message, got.Message)
	}

	deleteCookies := rec2.Result().Cookies()

	if len(deleteCookies) == 0 {
		t.Fatal("expected a delete cookie to be set")
	}

	if deleteCookies[0].MaxAge != -1 {
		t.Errorf("expected MaxAge -1 (delete), got %d", deleteCookies[0].MaxAge)
	}
}

func TestConsumeWithNoCookie(t *testing.T) {
	s := NewCookieStore()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	got := s.Consume(rec, req)

	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestConsumeWithInvalidJSON(t *testing.T) {
	s := NewCookieStore(WithCookieName("test_flash"))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "test_flash", Value: "not-json"})

	rec := httptest.NewRecorder()
	got := s.Consume(rec, req)

	if got != nil {
		t.Errorf("expected nil for invalid JSON, got %+v", got)
	}
}

func TestConsumeAfterConsumeReturnsNil(t *testing.T) {
	s := NewCookieStore(WithCookieName("test_flash"))
	msg := Message{Kind: "info", Title: "Hey", Message: "Hello."}

	rec := httptest.NewRecorder()
	s.Set(rec, msg)

	cookies := rec.Result().Cookies()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookies[0])

	rec2 := httptest.NewRecorder()
	first := s.Consume(rec2, req)

	if first == nil {
		t.Fatal("expected first consume to return message")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec3 := httptest.NewRecorder()
	second := s.Consume(rec3, req2)

	if second != nil {
		t.Errorf("expected second consume to return nil, got %+v", second)
	}
}
