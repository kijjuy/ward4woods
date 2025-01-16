package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"w4w/handlers"
	"w4w/models"
	"w4w/store"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
)

const (
	LogLevel         = slog.LevelDebug
	DayInSeconds     = 86400
	layoutName       = "_layout.html"
	templateDir      = "html"
	bootstrapCssPath = "html/bootstrap/css/bootstrap.css"
	bootstrapJsPath  = "html/bootstrap/js/bootstrap.js"
	jqueryPath       = "html/jquery.js"
	indexCssPath     = "html/index.css"
)

func init() {
	godotenv.Load()
	db := ConnectToDb()
	store.SetupProductsStore(db)
	gob.Register(new(models.Cart))
}

func ConnectToDb() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Error when opening connection to database")
	}

	err = db.Ping()

	if err != nil {
		slog.Error("Error Pinging database.", "Error", err)
		panic("Error pinging database. Shutting down now.")
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
	layout    *template.Template
	templates map[string]*template.Template
}

func NewTemplate(layoutPath, templatesDir string) (*Template, error) {
	layout, err := template.ParseGlob(layoutPath)
	if err != nil {
		return nil, err
	}

	templates := make(map[string]*template.Template)

	files, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := filepath.Base(file)
		if name != layoutName {
			tmpl, err := template.Must(layout.Clone()).ParseFiles(file)
			if err != nil {
				return nil, err
			}

			for _, template := range tmpl.Templates() {
				if template.Name() != "content" && template.Name() != "title" {
					templates[template.Name()] = template
				}
			}
			templates[strings.TrimSuffix(name, ".html")] = tmpl

		}
	}

	return &Template{
		layout:    layout,
		templates: templates,
	}, nil

}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found.", name)
	}

	if c.Request().Header.Get("Hx-Request") == "true" {
		return tmpl.ExecuteTemplate(w, name, data)
	} else {

		return tmpl.Execute(w, data)
	}
}

func newSessId() string {
	b := make([]byte, 32)
	if _, err := rand.Reader.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func CreateCartMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := session.Get("session", c)

		if err != nil {
			slog.Error("Error getting session", "Error", err)
			return err
		}

		if _, ok := session.Values["cart"].(*models.Cart); !ok {
			slog.Info("Creating new cart")
			session.Values["cart"] = new(models.Cart)

			session.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   DayInSeconds * 7,
				HttpOnly: true,
			}

			err = session.Save(c.Request(), c.Response())
			if err != nil {
				slog.Error("Error saving session", "Error", err)
				return err
			}
		}
		return next(c)
	}

}

func main() {
	e := echo.New()

	templateName := "html/" + layoutName
	t, err := NewTemplate(templateName, templateDir)
	if err != nil {
		slog.Error("Could not find template")
		panic("Error setting up templates.")
	}
	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(5))))

	sessionSecret := []byte(os.Getenv("SESSION_STORE_KEY"))
	sessionStore := sessions.NewCookieStore(sessionSecret)

	e.Use(session.Middleware(sessionStore))
	e.Use(CreateCartMiddleware)

	admin := e.Group("/admin")

	e.File("/bootstrap/css/bootstrap.css", bootstrapCssPath)
	e.File("/bootstrap/js/bootstrap.js", bootstrapJsPath)
	e.File("/jquery.js", jqueryPath)
	e.File("/index.css", indexCssPath)
	e.Static("/images/", "uploads")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/products/:id", handlers.ProductDetails)
	e.GET("/products", handlers.GetAllProducts)
	e.GET("/products/categories/:id", handlers.GetCategories)

	e.DELETE("/cart/:id", handlers.DeleteFromCart)
	e.POST("/cart/:id", handlers.AddToCart)
	e.DELETE("/cart", handlers.ClearCart)
	e.GET("/cart", handlers.ViewCart)

	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("ADMIN_USER") && password == os.Getenv("ADMIN_PASS") {
			return true, nil
		}
		slog.Info("", "Username", os.Getenv("ADMIN_USER"))
		slog.Info("", "Password", os.Getenv("ADMIN_PASS"))
		slog.Info("typed:", "Username", username)
		slog.Info("typed:", "Password", password)
		return false, nil
	}))

	admin.GET("", func(c echo.Context) error {
		return c.Render(http.StatusOK, "admin", nil)
	})
	admin.GET("/viewproducts", handlers.AdminGetProductsList)
	admin.GET("/newproduct", func(c echo.Context) error {
		return c.Render(http.StatusOK, "newProduct", nil)
	})
	admin.GET("/newImage", func(c echo.Context) error {
		return c.Render(http.StatusOK, "imageUpload", nil)
	})
	admin.GET("/products/edit/:id", handlers.EditProduct)
	admin.DELETE("/products/:id", handlers.DeleteProduct)
	admin.PUT("/products/:id", handlers.UpdateProduct)
	admin.POST("/products", handlers.NewProduct)

	e.Logger.Fatal(e.Start(":8080"))
}
