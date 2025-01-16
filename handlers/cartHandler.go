package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"w4w/models"
	"w4w/services"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AddToCart(c echo.Context) error {
	idStr := c.Param("id")
	productId, err := strconv.Atoi(idStr)

	if err != nil {
		return err
	}

	session, err := session.Get("session", c)

	if err != nil {
		logSessErr(err)
		return err
	}

	cart, ok := session.Values["cart"].(*models.Cart)

	slog.Info(fmt.Sprintf("cart: %+v", cart))

	if !ok {
		slog.Error("Error converting cart session value to models.Cart")
		return c.NoContent(http.StatusInternalServerError)
	}

	for _, cartProductId := range cart.Items {
		if productId == cartProductId {
			return c.Render(http.StatusOK, "cartDupeItem", nil)
		}
	}

	cart.Items = append(cart.Items, productId)

	session.Values["cart"] = cart

	err = session.Save(c.Request(), c.Response())

	if err != nil {
		slog.Error("Error saving session", "Error", err)
	}

	return c.Render(http.StatusOK, "cartAddSuccess", nil)
}

func ViewCart(c echo.Context) error {
	cart, err := getCartFromContext(c)

	if err != nil {
		logSessErr(err)
		return err
	}

	products := make([]models.Product, 0)

	for _, productId := range cart.Items {
		product, err := services.GetProductById(productId)

		if err != nil {
			slog.Error("Error getting product from service", "Error", err)
			return err
		}

		products = append(products, product)
	}

	return c.Render(http.StatusOK, "cart", products)
}

func DeleteFromCart(c echo.Context) error {
	idStr := c.Param("id")
	idToDelete, err := strconv.Atoi(idStr)

	if err != nil {
		slog.Error("Error converting id to int", "Error", err)
		return err
	}

	session, err := session.Get("session", c)

	if err != nil {
		slog.Error("Error getting session data", "Error", err)
		return err
	}

	cart, ok := session.Values["cart"].(*models.Cart)

	if !ok {
		slog.Error("Error getting cart from session")
		return c.NoContent(http.StatusInternalServerError)
	}

	deleted := false
	for i, productId := range cart.Items {
		if idToDelete == productId {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			deleted = true
		}
	}

	if !deleted {
		slog.Error("Error deleting product from cart")
		return c.NoContent(http.StatusInternalServerError)
	}

	session.Values["cart"] = cart

	err = session.Save(c.Request(), c.Response())

	if err != nil {
		slog.Error("Error saving session data", "Error", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func ClearCart(c echo.Context) error {
	session, err := session.Get("session", c)

	if err != nil {
		logSessErr(err)
		return err
	}

	session.Values["cart"] = new(models.Cart)

	err = session.Save(c.Request(), c.Response())

	if err != nil {
		slog.Error("Error saving session data", "Error", err)
		return err
	}

	return c.NoContent(http.StatusOK)
}

func getCartFromContext(c echo.Context) (*models.Cart, error) {
	session, err := session.Get("session", c)

	if err != nil {
		return nil, err
	}

	cart, ok := session.Values["cart"].(*models.Cart)

	if !ok {
		err = fmt.Errorf("Error getting cart models from session")
		return nil, err
	}

	return cart, err
}

func logSessErr(err error) {
	slog.Error("Error getting session data", "Error", err)
}
