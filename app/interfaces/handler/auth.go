package handler

import (
	"net/http"
	"todo-app/usecase"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func (h *handler) signUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	}

	generateHTML(w, r, nil, "layout", "signup")
}

func (h *handler) signupUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		errorResponse(w, r, err)
	}

	if err := h.userUsecase.SignUp(r.PostFormValue("email"), r.PostFormValue("password")); err != nil {
		errorResponse(w, r, err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	}

	generateHTML(w, r, nil, "layout", "login")
}

func (h *handler) authenticate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		errorResponse(w, r, err)
	}

	id, err := h.userUsecase.Authenticate(r.PostFormValue("email"), r.PostFormValue("password"))
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	session, err := h.sessionUsecase.CreateSession(id)
	if err != nil {
		errorResponse(w, r, err)
	}

	cookie := &http.Cookie{
		Name:     "todo_cookie",
		Value:    string(session.ID),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/tasks", http.StatusFound)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := h.session(r)
	if err != nil {
		errorResponse(w, r, err)
	} else if s == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if err := h.sessionUsecase.DeleteSession(s.UserID); err != nil {
		errorResponse(w, r, err)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *handler) session(r *http.Request) (*usecase.Session, error) {
	cookie, err := r.Cookie("todo_cookie")
	if err != nil {
		return nil, nil
	}

	session, err := h.sessionUsecase.Verify(usecase.SessionID(cookie.Value))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to verify cookie")
	}

	return session, nil
}
