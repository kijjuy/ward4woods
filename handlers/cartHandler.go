package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"w4w/services"

	"github.com/labstack/echo/v4"
)

func AddToCart(c echo.Context) error {
	idStr := c.Param("id")
	productId, err := strconv.Atoi(idStr)
	slog.Info("Hello from cart handler")

	if err != nil {
		return err
	}

	sessId := c.Get("session_id").(string)

	err = services.AddToCart(sessId, productId)

	if err != nil {
		slog.Error("couldnt add dupe item", "Error", err)
		return c.Render(http.StatusOK, "cartDupeItem", nil)
	}

	return c.Render(http.StatusOK, "cartAddSuccess", nil)
}

func ViewCart(c echo.Context) error {
	sessId := c.Get("session_id").(string)

	products, err := services.GetCartItems(sessId)

	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "cart", products)
}

func DeleteFromCart(c echo.Context) error {
	sessId := c.Get("session_id").(string)
	cartItemIdStr := c.Param("id")
	cartItemId, err := strconv.Atoi(cartItemIdStr)

	if err != nil {
		return err
	}

	err = services.DeleteFromCart(sessId, cartItemId)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func ClearCart(c echo.Context) error {
	sessId := c.Get("session_id").(string)

	err := services.ClearCart(sessId)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
