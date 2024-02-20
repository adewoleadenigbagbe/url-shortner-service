package core

func RegisterRoutes(app *BaseApp) {
	router := app.echo
	router.POST("/api/v1/auth/register", app.AuthService.RegisterUser)
	//route.POST("/api/v1/auth/sign-in", services.LoginUser)
}
