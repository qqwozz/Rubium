func mainHandler() {
	vkProvider := &auth.VKProvider{
		ClientID:    os.Getenv("VK_CLIENT_ID"),
		RedirectURI: os.Getenv("VK_REDIRECT_URI"),
		Scope:       "email",
	}

	registry := auth.NewRegistry(vkProvider, googleProvider)
	authHandler := &auth.Handler{
		registry: registry,
		storage:  yourStorage,
	}

	r := gin.Default()
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/providers", authHandler.ListProviders)
		authGroup.GET("/:provider/url", authHandler.GetAuthURL)
		authGroup.POST("/:provider/callback", authHandler.Callback)
	}

	r.Run(":8080")
}