package middleware

import (
	"net/http"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/gin-gonic/gin"
)

const CookieToken string = "token"
const ContextUUID string = "uuid"
const ContextUser string = "user"

type authServer interface {
	CreateUser(email string, password string) (uuid string, err error)
	VerifyToken(token string) (*models.User, error)
}

type userServer interface {
	GetUserByUUID(uuid string) (*models.User, error)
}

type authMiddleware struct {
	authServer
	userServer
}

func NewAuthMiddleware(authServer authServer, userServer userServer) *authMiddleware {
	return &authMiddleware{
		authServer: authServer,
		userServer: userServer,
	}
}

func (m *authMiddleware) RegisterNewUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		err := c.Bind(&user)
		if err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"err": err})
			c.Abort()
			return
		}

		if user.Email == "" || user.Password == "" || user.Username == "" || user.Name == "" {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"err": "All data are required"})
			c.Abort()
			return
		}

		uuid, err := m.CreateUser(user.Email, user.Password)
		if err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{"err": err})
			c.Abort()
			return
		}

		c.Set(ContextUUID, uuid)
		//TODO: make secure
		c.SetCookie(CookieToken, uuid, 1000, "", "", false, false)
		c.Next()

	}
}

/*
Extract the user information (uuid) from the request's cookies and add it to the request context.
If not valid token (or no token at all) is found, nothing is done to the request.
This middleware is useful to get the user info in pages where people can visit as guests as well.
For example, the home page can be visited by guests, but logged in users will see different content in the same page.
*/
func (m *authMiddleware) ExtractUserFromToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(CookieToken)

		tokenUser, _ := m.VerifyToken(token)

		//If there is no user (no valid token) then go to the next handler without adding a user (guest mode)
		if tokenUser == nil {
			c.Next()
			return
		}

		user, err := m.userServer.GetUserByUUID(tokenUser.UUID)
		if err != nil {
			c.HTML(http.StatusTemporaryRedirect, "error.html", gin.H{"err": err})
			c.Abort()
			return
		}

		//Inject the user into the request context and move on to the next handler
		c.Set(ContextUser, user)
		c.Next()

	}
}

/*
Ensure the user is authenticated before moving to the next handler.
Unauthenticated users are redirected to the login page.
Use this middleware in endpoints that require authenticated users.
*/
func (m *authMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get(ContextUser)
		//If the user is not authenticated, redirect to login page.
		if !ok {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
