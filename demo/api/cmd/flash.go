package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/oullin/inertia-go/demo/api/internal/flash"
)

var db *sql.DB

const flashCookieName = "beacon_flash"

func setFlash(w http.ResponseWriter, flash flash.Message) {
	data, _ := json.Marshal(flash)

	http.SetCookie(w, &http.Cookie{
		Name:     flashCookieName,
		Value:    url.QueryEscape(string(data)),
		Path:     "/",
		HttpOnly: true,
	})
}

func consumeFlash(w http.ResponseWriter, r *http.Request) *flash.Message {
	cookie, err := r.Cookie(flashCookieName)

	if err != nil {
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:   flashCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	value, err := url.QueryUnescape(cookie.Value)

	if err != nil {
		return nil
	}

	var payload flash.Message

	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return nil
	}

	return &payload
}
