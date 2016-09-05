package auth_html

import (
	"net/http"

	"html/template"

	"github.com/bborbe/http/header"
	"github.com/golang/glog"
	"time"
)

const (
	fieldNameLogin    = "login"
	fieldNamePassword = "password"
	cookieName        = "auth-http-proxy-token"
)

var expiration = time.Now().Add(24 * time.Hour)

type Check func(username string, password string) (bool, error)

type handler struct {
	handler http.HandlerFunc
	check   Check
}

func New(subhandler http.HandlerFunc, check Check) *handler {
	h := new(handler)
	h.handler = subhandler
	h.check = check
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(2).Infof("check html auth")
	if err := h.serveHTTP(responseWriter, request); err != nil {
		glog.Warningf("check html auth failed: %v", err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	glog.V(2).Infof("check html auth")
	valid, err := h.validateLoginCookie(request)
	if err != nil {
		glog.V(2).Infof("validate login failed: %v", err)
		return err
	}
	if valid {
		glog.V(2).Infof("login is valid, forward request")
		h.handler(responseWriter, request)
	}
	return h.validateLoginParams(responseWriter, request)
}

func (h *handler) validateLoginParams(responseWriter http.ResponseWriter, request *http.Request) error {
	glog.V(2).Infof("validate login via params")
	login := request.FormValue(fieldNameLogin)
	password := request.FormValue(fieldNamePassword)
	valid, err := h.check(login, password)
	if err != nil {
		glog.V(2).Infof("check login failed: %v", err)
		return err
	}
	if !valid {
		glog.V(2).Infof("login failed, show login form")
		return h.loginForm(responseWriter)
	}
	glog.V(2).Infof("login success, set cookie")
	http.SetCookie(responseWriter, &http.Cookie{
		Name:    cookieName,
		Value:   header.CreateAuthorizationToken(login, password),
		Expires: expiration,
		Path:    "/",
		Domain:  request.URL.Host,
	},
	)
	http.Redirect(responseWriter, request, "/", http.StatusTemporaryRedirect)
	return nil
}

func (h *handler) validateLoginCookie(request *http.Request) (bool, error) {
	glog.V(2).Infof("validate login via cookie")
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		glog.V(2).Infof("get cookie %v failed: %v", cookieName, err)
		return false, nil
	}
	user, pass, err := header.ParseAuthorizationToken(cookie.Value)
	if err != nil {
		glog.V(2).Infof("parse header failed: %v", err)
		return false, nil
	}
	return h.check(user, pass)
}

func (h *handler) loginForm(responseWriter http.ResponseWriter) error {
	glog.V(2).Infof("login form")
	responseWriter.WriteHeader(http.StatusUnauthorized)
	var t = template.Must(template.New("loginForm").Parse(HTML))
	data := struct {
		CookieName string
		Title      string
	}{
		CookieName: cookieName,
		Title:      "Login",
	}
	responseWriter.Header().Add("Content-Type", "text/html")
	return t.Execute(responseWriter, data)
}

const HTML = `<!DOCTYPE html>
<html>
<title>{{.Title}}</title>
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta http-equiv="Content-Language" content="en">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="author" content="Benjamin Borbe">
<meta name="description" content="Login Form">
<link rel="icon" href="data:;base64,=">
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
</script>
</head>
<body>
<div class="view-container">
	<div class="container">
		<div class="starter-template">

			<form name="loginForm" class="form-horizontal" action="" method="post">
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
