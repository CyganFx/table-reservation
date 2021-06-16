package http_v1

import (
	"bytes"
	"fmt"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/nosurf"
	"net/http"
	"time"
)

func (h *handler) SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Frame-Options", "deny")
		c.Next()
	}
}

func isAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)
	if session.Get("authenticatedUserID") == nil {
		return false
	}
	return true
}

func (h *handler) RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c) {
			session := sessions.Default(c)
			session.Set("redirectPathAfterLogin", c.Request.URL.Path)
			session.Save()

			http.Redirect(c.Writer, c.Request, "/api/users/login", http.StatusSeeOther)
			c.Abort()
			return
		}
		c.Writer.Header().Add("Cache-Control", "no-store")
	}
}

func (h *handler) RequireRole(roleID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("role") != roleID {
			http.Redirect(c.Writer, c.Request, "/api/users/login", http.StatusSeeOther)
			c.Abort()
			return
		}
		c.Writer.Header().Add("Cache-Control", "no-store")
	}
}

func csrfFailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", nosurf.Reason(r))
}

//front end middleware

func (h *handler) addDefaultData(td *templateData, c *gin.Context) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(c.Request)
	td.CurrentYear = time.Now().Year()
	session := sessions.Default(c)

	if session.Get("flash") != nil {
		td.Flash = session.Get("flash").(string)
		session.Delete("flash")
	}

	if session.Get("authenticatedUserID") != nil {
		if td.User == nil {
			td.User = domain.NewUser()
		}
		if td.User.ID == 0 {
			td.User.ID = session.Get("authenticatedUserID").(int)
		}
	}
	if session.Get("role") != nil {
		if td.User == nil {
			td.User = domain.NewUser()
		}
		if td.User.Role.ID == 0 {
			td.User.Role.ID = session.Get("role").(int)
		}
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

	err := ts.Execute(buf, h.addDefaultData(td, c))
	if err != nil {
		h.errors.ServerError(c, err)
		return
	}

	buf.WriteTo(c.Writer)
}
