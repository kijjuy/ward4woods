package store

import (
	"w4w/models"
)

func AddToCart(sessId string, productId int) error {
	_, err := db.Exec("INSERT INTO cart_items (cart_id, product_id) VALUES("+
		"(SELECT cart_id FROM carts WHERE session_id = $1), "+
		"(SELECT product_id FROM products WHERE product_id = $2))", sessId, productId)

	return err
}

func GetCartIdBySessId(sessId string) (int, error) {
	row := db.QueryRow("SELECT cart_id FROM carts WHERE session_id = $1", sessId)

	var cartId int

	err := row.Scan(&cartId)

	return cartId, err
}

func GetCartDisplayModelBySessId(sessId string) (models.CartDisplayProducts, error) {
	rows, err := db.Query("SELECT products.*, cart_item_id FROM cart_items "+
		"JOIN products ON cart_items.product_id = products.product_id "+
		"WHERE cart_id = (SELECT cart_id FROM carts WHERE session_id = $1)", sessId)

	if err != nil {
		return nil, err
	}

	cartDisplayProducts := models.NewCartDisplayProducts()

	for rows.Next() {
		product := models.NewProduct()
		var cartItemId int
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Description, &product.Category, &cartItemId)

		if err != nil {
			return nil, err
		}

		cartDisplayProduct := models.CartDisplayProduct{
			Product:    product,
			CartItemId: cartItemId,
		}

		cartDisplayProducts = append(cartDisplayProducts, cartDisplayProduct)
	}

	return cartDisplayProducts, nil
}

func GetCartItemsBySessId(sessId string) (models.Products, error) {
	displayModels, err := GetCartDisplayModelBySessId(sessId)
	if err != nil {
		return nil, err
	}

	products := models.NewProducts()

	for _, viewModel := range displayModels {
		products = append(products, viewModel.Product)
	}
	return products, err
}

func CreateCart(sessId string) error {
	_, err := db.Exec("INSERT INTO carts (session_id) Values($1)", sessId)
	return err
}

func CartItemOwenershipIsValid(sessId string, cartItemId int) (bool, error) {
	row := db.QueryRow("SELECT session_id FROM cart_items JOIN carts ON cart_items.cart_id = carts.cart_id WHERE cart_item_id = $1", cartItemId)

	var cartOwnerSessId string

	err := row.Scan(&cartOwnerSessId)

	if err != nil || sessId != cartOwnerSessId {
		return false, err
	}

	return true, err
}

func DeleteFromCart(sessId string, cartItemId int) error {
	_, err := db.Exec("DELETE FROM cart_items WHERE cart_item_id = $1", cartItemId)
	return err
}

func DeleteAllCartItems(sessId string) error {
	_, err := db.Exec("DELETE FROM cart_items WHERE cart_id = (SELECT cart_id FROM carts WHERE session_id = $1)", sessId)

	return err
}
