package http_v1

import (
	"bytes"
	"fmt"
	"github.com/CyganFx/table-reservation/tasty-food-cafe/pkg/domain"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Frame-Options", "deny")
		c.Next()
	}
}

//Front-end middleware

func addDefaultData(td *templateData, c *gin.Context) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	session := sessions.Default(c)

	if session.Get("flash") != nil {
		td.Flash = session.Get("flash").(string)
		session.Delete("flash")
	}

	if session.Get("role") != nil {
		td.User = domain.NewUser()
		td.User.Role.ID = session.Get("role").(int)
	}

	td.IsAuthenticated = isAuthenticated(c)

	session.Save()
	return td
}

func (h *handler) render(c *gin.Context, name string, td *templateData) {
	ts, ok := h.templateCache[name]
	if !ok {
		h.errors.ServerError(c, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, addDefaultData(td, c))
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	buf.WriteTo(c.Writer)
}

func isAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)
	if session.Get("authenticatedUserID") == nil {
		return false
	}
	return true
}

//func RequireAuthentication() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if !isAuthenticated(c) {
//			session := sessions.Default(c)
//			session.Set("redirectPathAfterLogin", c.Request.URL.Path)
//			session.Save()
//
//			http.Redirect(c.Writer, c.Request, "/api/user/login", http.StatusSeeOther)
//			return
//		}
//
//		c.Writer.Header().Add("Cache-Control", "no-store")
//		c.Next()
//	}
//}
