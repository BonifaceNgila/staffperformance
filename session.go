package main

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	// Initialize session store
	store = sessions.NewCookieStore([]byte("your-secret-key-change-this-in-production"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	}
	gob.Register(User{})
}

// SetSession sets user session
func SetSession(w http.ResponseWriter, r *http.Request, user *User) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}
	session.Values["userID"] = user.ID
	session.Values["username"] = user.Username
	return session.Save(r, w)
}

// GetSession gets current user from session
func GetSession(r *http.Request) (*User, error) {
	session, err := store.Get(r, "session")
	if err != nil {
		return nil, err
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		return nil, http.ErrNoCookie
	}

	return GetUserByID(userID)
}

// ClearSession clears user session
func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// RequireAuth middleware to protect routes
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetSession(r)
		if err != nil || user == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
