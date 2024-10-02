package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"w4w/handlers"
	"w4w/store"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const LogLevel = slog.LevelDebug

func init() {
	godotenv.Load()
	db := ConnectToDb()
	store.SetupProductsStore(db)
}

func ConnectToDb() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Error when opening connection to database")
	}

	return db
}

func SetupLogging() {
	opts := &slog.HandlerOptions{Level: LogLevel}
	var handler slog.Handler = slog.NewTextHandler(os.Stdin, opts)
	if os.Getenv("APP_ENV ") == "production" {
		handler = slog.NewJSONHandler(os.Stdin, opts)
	}
	slog.SetDefault(slog.New(handler))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newSessId() string {
	b := make([]byte, 32)
	if _, err := rand.Reader.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessIdCookie, err := c.Cookie("session_id")

		var sessIdStr string

		if err == http.ErrNoCookie {
			sessIdStr := newSessId()
			c.SetCookie(&http.Cookie{
				Name:     "session_id",
				Value:    sessIdStr,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				MaxAge:   86400 * 7,
			})
		} else {
			sessIdStr = sessIdCookie.Value
		}

		c.Set("session_id", sessIdStr)
		return next(c)
	}
}

func main() {
	e := echo.New()

	e.Use(SessionMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	t := &Template{
		template.Must(template.ParseGlob("html/*.html")),
	}

	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/products/:id", handlers.ProductDetails)
	e.DELETE("/products/:id", handlers.DeleteProduct)
	e.PUT("/products/:id", handlers.UpdateProduct)
	e.POST("/products", handlers.NewProduct)
	e.GET("/products", handlers.GetAllProducts)

	e.GET("/products/categories/:id", handlers.GetCategories)

	e.DELETE("/cart/:id", handlers.DeleteFromCart)
	e.POST("/cart/:id", handlers.AddToCart)
	e.DELETE("/cart", handlers.ClearCart)
	e.GET("/cart", handlers.ViewCart)

	e.GET("/admin/viewproducts", handlers.AdminGetProductsList)
	e.GET("/admin", func(c echo.Context) error {
		return c.Render(http.StatusOK, "admin", nil)
	})
	e.GET("/admin/newproduct", func(c echo.Context) error {
		return c.Render(http.StatusOK, "newProduct", nil)
	})
	e.GET("/admin/edit/:id", handlers.EditProduct)

	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":8080"))
}
