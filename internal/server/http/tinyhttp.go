package tinyhttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/shashimalcse/tiny-is/internal/organization"
)

type contextKey string

const SERVER_URL contextKey = "server_url"
const SERVER_SCHEME contextKey = "server_scheme"

type TinyServeMux struct {
	mux                 *http.ServeMux
	organizationService organization.OrganizationService
}

func NewTinyServeMux(organizationService organization.OrganizationService) *TinyServeMux {
	return &TinyServeMux{
		mux:                 http.NewServeMux(),
		organizationService: organizationService,
	}
}

func (c *TinyServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	ctx := context.WithValue(r.Context(), SERVER_URL, r.Host)
	ctx = context.WithValue(ctx, SERVER_SCHEME, scheme)
	r = r.WithContext(ctx)

	if strings.HasPrefix(r.URL.Path, "/o/") {
		parts := strings.SplitN(r.URL.Path, "/", 4)
		if len(parts) >= 4 {
			orgName := parts[2]
			// Check if the organization exists
			ctx := r.Context()
			org, err := c.organizationService.GetOrganizationByName(ctx, orgName)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					http.NotFound(w, r)
					return
				}
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if org.Name != orgName {
				http.NotFound(w, r)
				return
			}
			r.Header.Set("org_name", orgName)
			r.Header.Set("org_id", org.Id)
			r.URL.Path = "/" + parts[3]
			c.mux.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (c *TinyServeMux) Handle(pattern string, handler http.Handler) {
	c.mux.Handle(pattern, handler)
}

func (c *TinyServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	c.mux.HandleFunc(pattern, handler)
}

func (c *TinyServeMux) getFullURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.RequestURI())
}
