package handler

import (
	"dg-server/core"
	"dg-server/store"
	"dg-server/tpl"
	"encoding/json"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func listTemplete(c echo.Context) error {
	q := &core.TempleteQuery{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, q)

	tplStore := store.Stores().TempleteStore
	list, total, err := tplStore.List(q)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Data: &PageResult{Total: total, List: list},
	})
}

func loadTemplete(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtUserClaims)

	q := &core.TempleteLoadReq{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, q)

	content, _ := tpl.GetTableScript(q, claims.Name)
	return c.JSON(http.StatusOK, &StandardResult{
		Data: content,
	})
}
