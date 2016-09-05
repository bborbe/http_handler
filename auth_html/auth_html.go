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
		Title      string
	}{
		CookieName: h.cookieName,
		Title:      "Login",
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
<title>{{.Title}}</title>
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta http-equiv="Content-Language" content="en">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="author" content="Benjamin Borbe">
<meta name="description" content="Booking App">
<link rel="stylesheet" type="text/css" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
<link rel="stylesheet" type="text/css" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap-theme.min.css">
<style>
html {
	position: relative;
	min-height: 100%;
}

body {
	margin-top: 60px;
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
<div class="view-container">
	<div class="container">
		<div class="starter-template">

			<form name="loginForm" class="form-horizontal" action="javascript:login()">
				<fieldset>

					<legend>Login required</legend>

					<div class="form-group">
						<label class="col-md-3 control-label" for="login">Login</label>

						<div class="col-md-3">
							<input type="text" id="login" name="login" min="1" max="255" required="" placeholder="login" class="form-control input-md">
						</div>
					</div>

					<div class="form-group">
						<label class="col-md-3 control-label" for="password">Password</label>

						<div class="col-md-3">
							<input type="password" id="password" name="password" min="1" max="255" required="" placeholder="password" class="form-control input-md">
						</div>
					</div>

					<div class="form-group">
						<label class="col-md-3 control-label" for="singlebutton"></label>

						<div class="col-md-3">
							<input type="submit" id="singlebutton" name="singlebutton" class="btn btn-primary" value="login">
						</div>
					</div>

				</fieldset>
			</form>
		</div>
	</div>
</div>
</body>
</html>`
