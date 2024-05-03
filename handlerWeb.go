package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func loginHandler(c *gin.Context) {
	var form = struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Redirect string `form:"redirect" binding:"omitempty"`
		Retry    string `form:"retry" binding:"omitempty"`
	}{
		Retry:    "/login",
		Redirect: config.ProxyPath,
	}

	// see if a form request with login data
	contentType := c.GetHeader("Content-Type")

	// only for web form encoded login
	if contentType != "application/x-www-form-urlencoded" {
		c.Redirect(http.StatusFound, form.Retry)
		return
	}

	if err := c.ShouldBind(&form); err != nil {
		c.Redirect(http.StatusFound, form.Retry)
		return
	}

	user, err := identityForUser(form.Username, form.Password)
	if err != nil {
		c.Redirect(http.StatusFound, form.Retry)
		return
	}

	// set authToken as HttpOnly cookie
	maxAge := int(config.expireDuration.Seconds())
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"authToken",
		user.token,
		maxAge,
		"/",
		"",
		config.Release,
		true,
	)
	c.Redirect(http.StatusFound, form.Redirect)
}

// logoutHandler clears the auth token cookie and redirects the user if a redirect URL is provided.
// If web login is disabled, it returns an empty JSON response.
// This handler does not require authentication.
func logoutHandler(c *gin.Context) {
	// Clear the auth token cookie
	c.SetCookie(
		"authToken",
		"",
		-1,
		"/",
		"",
		config.Release,
		true,
	)
	err := c.Request.ParseForm()
	redirect := "."
	if err == nil {
		redirect = c.Request.Form.Get("redirect")
	}
	c.Redirect(http.StatusFound, redirect)
}

func webFormHandler(c *gin.Context) {
	if _, err := c.Cookie("authToken"); err == nil {
		// allow logout if user if there is an authToken
		c.HTML(http.StatusOK, "logout.html", gin.H{
			"redirect": "/login",
		})
		return
	}

	// otherwise show login
	c.HTML(http.StatusOK, "login.html", gin.H{
		"redirect": config.ProxyPath,
	})
}
