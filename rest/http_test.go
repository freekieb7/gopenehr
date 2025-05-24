package rest_test

import (
	"testing"

	"github.com/freekieb7/gopenehr/rest"
)

func TestServer(t *testing.T) {
	server := rest.NewServer("openEHR REST API")

	// Add routes
	server.Router.GET("/", rest.HandleServerInfo())

	server.Router.Group("/ehr", func(group *rest.Router) {
		group.GET("", rest.HandleListEhr())
		group.POST("", rest.HandleCreateEhr())

		group.Group("/{ehr_id}", func(group *rest.Router) {
			group.GET("", rest.HandleGetEhrById())
			group.DELETE("", rest.HandleDeleteEhrById())

			group.Group("/ehr_status", func(group *rest.Router) {
				group.GET("", rest.HandleGetEhrStatusById())
				group.PUT("", rest.HandleUpdateEhrStatusById())
			})

			group.Group("/composition", func(group *rest.Router) {
				group.GET("", rest.HandleListComposition())
				group.POST("", rest.HandleCreateComposition())

				group.Group("/{composition_id}", func(group *rest.Router) {
					group.GET("", rest.HandleGetCompositionById())
					group.PUT("", rest.HandleUpdateCompositionById())
					group.DELETE("", rest.HandleDeleteCompositionById())
				})
			})

			group.Group("/folder", func(group *rest.Router) {
				group.GET("", rest.HandleListFolder())
				group.POST("", rest.HandleCreateFolder())

				group.Group("/{folder_id}", func(group *rest.Router) {
					group.GET("", rest.HandleGetFolderById())
					group.PUT("", rest.HandleUpdateFolderById())
					group.DELETE("", rest.HandleDeleteFolderById())
				})
			})
		})
	})

	// server.ListenAndServe(":8080")
	// todo
}
