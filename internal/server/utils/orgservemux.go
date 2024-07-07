package utils

import (
	"net/http"
	"strings"

	"github.com/shashimalcse/tiny-is/internal/organization"
)

type OrgServeMux struct {
	mux                 *http.ServeMux
	organizationService organization.OrganizationService
}

func NewOrgServeMux(organizationService organization.OrganizationService) *OrgServeMux {
	return &OrgServeMux{
		mux:                 http.NewServeMux(),
		organizationService: organizationService,
	}
}

func (c *OrgServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/o/") {
		parts := strings.SplitN(r.URL.Path, "/", 4)
		if len(parts) >= 4 {
			orgName := parts[2]
			// Check if the organization exists
			ctx := r.Context()
			org, err := c.organizationService.GetOrganizationByName(ctx, orgName)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if org.Name != orgName {
				http.NotFound(w, r)
				return
			}
			r.Header.Set("org_name", orgName)
			r.Header.Set("org_id", org.Id)
			r.URL.Path = "/" + parts[3] // Update the URL path to match the handler's expected path
			c.mux.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (c *OrgServeMux) Handle(pattern string, handler http.Handler) {
	c.mux.Handle(pattern, handler)
}

func (c *OrgServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	c.mux.HandleFunc(pattern, handler)
}
