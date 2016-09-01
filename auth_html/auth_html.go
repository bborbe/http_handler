package auth_html

import (
	"fmt"
	"net/http"

	"html/template"

	"github.com/bborbe/http/header"
	"github.com/golang/glog"
)

type Check func(username string, password string) (bool, error)

type handler struct {
	handler    http.HandlerFunc
	check      Check
	cookieName string
}

func New(subhandler http.HandlerFunc, check Check) *handler {
	h := new(handler)
	h.handler = subhandler
	h.check = check
	h.cookieName = "auth_html_login"
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(2).Infof("check html auth")
	if err := h.serveHTTP(responseWriter, request); err != nil {
		if err := h.printForm(responseWriter); err != nil {
			glog.Warningf("print login form failed: %v", err)
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *handler) printForm(responseWriter http.ResponseWriter) error {
	responseWriter.WriteHeader(http.StatusUnauthorized)
	var t = template.Must(template.New("loginForm").Parse(HTML))
	data := struct {
		CookieName string
	}{
		CookieName: h.cookieName,
	}
	return t.Execute(responseWriter, data)
}

func (h *handler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	glog.V(2).Infof("check html auth")
	cookie, err := request.Cookie(h.cookieName)
	if err != nil {
		glog.V(2).Infof("get cookie %v failed: %v", h.cookieName, err)
		return err
	}
	user, pass, err := header.ParseAuthorizationToken(cookie.Value)
	if err != nil {
		glog.Warningf("parse header failed: %v", err)
		return err
	}
	result, err := h.check(user, pass)
	if err != nil {
		glog.Warningf("check auth for user %v failed: %v", user, err)
		return err
	}
	if !result {
		glog.V(1).Infof("auth invalid for user %v", user)
		return fmt.Errorf("auth invalid for user %v", user)
	}
	h.handler(responseWriter, request)
	return nil
}

const HTML = `<!DOCTYPE html>
<html>
<head>
<style>
label {
	width: 100px;
	display: inline-block;
}
</style>
<script>
function login() {
	var value = btoa(document.getElementById("login").value+":"+document.getElementById("password").value);
	document.cookie = "{{.CookieName}}=" + value + "; path=/";
	document.location.reload();
}
</script>
</head>
<body>
<h1>Login required</h1>
<form action="javascript:login()">
<div><label for="login">Login:</label><input type="text" name="login" id="login"></div>
<div><label for="password">Password:</label><input type="password" name="password" id="password"></div>
<input type="submit" value="login">
</form>
</body>
<html>
`
