package grant_handlers

import (
	"net/http"

	"github.com/shashimalcse/tiny-is/internal/server/models"
)

type GrantHandler interface {
	HandleGrant(r *http.Request) (models.TokenResponse, error)
}
