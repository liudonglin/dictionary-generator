package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"dg-server/core"
	"dg-server/store"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const secret = "my token secret"

type jwtUserClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func login(c echo.Context) error {
	uiUser := &core.User{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, uiUser)

	if uiUser.Login == "" {
		return &BusinessError{Message: "用户登录名不能为空！"}
	}
	if uiUser.Password == "" {
		return &BusinessError{Message: "用户登录密码不能为空！"}
	}

	dbInstance := store.Stores()
	dbUser, _ := dbInstance.UserStore.FindLogin(uiUser.Login)

	uiUser.Password = store.EncryptionPassword(uiUser.Password)
	if uiUser.Login == dbUser.Login && uiUser.Password == dbUser.Password {

		if dbUser.Active == false {
			return &BusinessError{Message: "该用户已被禁用！"}
		}

		// Set custom claims
		claims := &jwtUserClaims{
			dbUser.Login,
			dbUser.Admin,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return err
		}

		result := StandardResult{}
		result.Message = "登录成功!"
		result.Data = t
		return c.JSON(http.StatusOK, &result)
	}
	return &BusinessError{Message: "用户登录名或密码错误！"}
}
