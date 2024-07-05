package authn

import (
	"github.com/a-h/templ"
	"github.com/shashimalcse/tiny-is/internal/authn/screens"
	"github.com/shashimalcse/tiny-is/internal/cache"
	"github.com/shashimalcse/tiny-is/internal/session"
	"github.com/shashimalcse/tiny-is/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type Authn struct {
	CacheService *cache.CacheService
	SessionStore *session.SessionStore
	userService  *user.UserService
}

func NewAuthn(cacheService *cache.CacheService, sessionStore *session.SessionStore, userService *user.UserService) *Authn {
	return &Authn{
		CacheService: cacheService,
		userService:  userService,
		SessionStore: sessionStore,
	}
}

func (authn Authn) GetLoginPage(sessionDataKey string) templ.Component {
	return screens.LoginPage(sessionDataKey)
}

func (authn Authn) ValidateUser(username, password string) (bool, error) {
	hashedPassword, err := authn.userService.GetHashedPasswordByUsername(username)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (authn Authn) GetUserIdByUsername(username string) (string, error) {
	return authn.userService.GetUserIdByUsername(username)
}

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(bytes), nil
// }
