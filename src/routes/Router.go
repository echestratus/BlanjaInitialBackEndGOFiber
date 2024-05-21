package routes

import (
	"BlanjaInitialBackEndGOFiber/src/controllers"
	"BlanjaInitialBackEndGOFiber/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	//Product Routes
	app.Get("/products", controllers.GetAllProducts)
	app.Get("/product/:id", controllers.GetDetailProduct)
	app.Post("/product", middlewares.JwtMiddleware(), controllers.CreateNewProduct)
	app.Put("/product/:id", controllers.UpdateProduct)
	app.Delete("/product/:id", controllers.DeleteProduct)
	//Category Routes
	app.Get("/categories", controllers.GetAllCategories)
	app.Get("/category/:id", controllers.GetDetailCategory)
	app.Post("/category", controllers.CreateCategory)
	app.Put("/category/:id", controllers.UpdateCategory)
	app.Delete("/category/:id", controllers.DeleteCategory)
	//Customer Routes
	app.Get("/customers", controllers.GetAllCustomers)
	app.Get("/customer/:id", controllers.GetCustomerById)
	app.Get("/customer", middlewares.JwtMiddleware(), controllers.GetDetailCustomer)
	app.Post("/customer", controllers.CreateNewCustomer)
	app.Put("/customer/:id", controllers.UpdateCustomerById)
	app.Put("/customer", middlewares.JwtMiddleware(), controllers.UpdateCustomer)
	app.Delete("/customer/:id", controllers.DeleteCustomer)
	//Seller Routes
	app.Get("/sellers", controllers.GetAllSellers)
	app.Get("/seller/:id", controllers.GetSellerById)
	app.Get("/seller", middlewares.JwtMiddleware(), controllers.GetDetailSeller)
	app.Post("/seller", controllers.CreateNewSeller)
	app.Put("/seller/:id", controllers.UpdateSellerById)
	app.Put("/seller", middlewares.JwtMiddleware(), controllers.UpdateSeller)
	app.Delete("/seller/:id", controllers.DeleteSeller)
	//Auth Routes
	app.Post("/customer/login", controllers.AuthLogin)
	app.Post("/refreshToken", controllers.RefreshToken)
	//Upload Routes
	app.Post("/upload", controllers.UploadLocal)
	app.Post("/uploadServer", controllers.UploadFileServer)
}
