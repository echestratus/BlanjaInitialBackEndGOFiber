package routes

import (
	"BlanjaInitialBackEndGOFiber/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	//Product Routes
	app.Get("/products", controllers.GetAllProducts)
	app.Get("/product/:id", controllers.GetDetailProduct)
	app.Post("/product", controllers.CreateNewProduct)
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
	app.Get("/customer/:id", controllers.GetDetailCustomer)
	app.Post("/customer", controllers.CreateNewCustomer)
	app.Put("/customer/:id", controllers.UpdateCustomer)
	app.Delete("/customer/:id", controllers.DeleteCustomer)
}
