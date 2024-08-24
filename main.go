package main

import (
	"fmt"
	"github.com/aaanger/lenslocked/controllers"
	"github.com/aaanger/lenslocked/models"
	"github.com/aaanger/lenslocked/templates"
	"github.com/aaanger/lenslocked/views"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", usersC.New)

	r.Post("/signup", usersC.Create)

	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signin", usersC.SignIn)

	r.Post("/signin", usersC.ProcessSignIn)

	r.Get("/users/me", usersC.CurrentUser)

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "gallery.gohtml"))
	r.Get("/gallery", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting server on :3000...")

	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect([]byte(csrfKey),
		//fix before deploying
		csrf.Secure(false))
	http.ListenAndServe(":3000", csrfMw(r))
}
