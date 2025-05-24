package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freekieb7/gopenehr/rest"
)

func main() {
	ctx := context.Background()

	if err := Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func Run(ctx context.Context) error {
	// Add gracefull shutdown support by listening for interuptions
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	port := 8080
	addr := fmt.Sprintf("0.0.0.0:%d", port)

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

			// todo contribution (just not sure what the difference is with demographics)
		})
	})

	server.Router.Group("/template", func(group *rest.Router) {
		group.GET("", rest.HandleListTemplates())
		group.POST("", rest.HandleCreateTemplate())

		group.Group("/{template_id}", func(group *rest.Router) {
			group.GET("", rest.HandleGetTemplateById())
			group.PUT("", rest.HandleUpdateTemplateById())
			group.DELETE("", rest.HandleDeleteTemplateById())
		})
	})

	server.Router.Group("/agent", func(group *rest.Router) {
		group.GET("", rest.HandleListAgent())
		group.POST("", rest.HandleCreateAgent())

		group.Group("/{agent_id}", func(group *rest.Router) {
			group.GET("", rest.HandleGetAgentById())
			group.PUT("", rest.HandleUpdateAgentById())
			group.DELETE("", rest.HandleDeleteAgentById())
		})
	})

	server.Router.Group("/group", func(group *rest.Router) {
		group.GET("", rest.HandleListGroup())
		group.POST("", rest.HandleCreateGroup())

		group.Group("/{group_id}", func(group *rest.Router) {
			group.GET("", rest.HandleGetGroupById())
			group.PUT("", rest.HandleUpdateGroupById())
			group.DELETE("", rest.HandleDeleteGroupById())
		})
	})

	server.Router.Group("/organisation", func(group *rest.Router) {
		group.GET("", rest.HandleListOrganisation())
		group.POST("", rest.HandleCreateOrganisation())

		group.Group("/{organisation_id}", func(group *rest.Router) {
			group.GET("", rest.HandleGetOrganisationById())
			group.PUT("", rest.HandleUpdateOrganisationById())
			group.DELETE("", rest.HandleDeleteOrganisationById())
		})
	})

	server.Router.Group("/person", func(group *rest.Router) {
		group.GET("", rest.HandleListPerson())
		group.POST("", rest.HandleCreatePerson())

		group.Group("/{person_id}", func(group *rest.Router) {
			group.GET("", rest.HandleGetPersonById())
			group.PUT("", rest.HandleUpdatePersonById())
			group.DELETE("", rest.HandleDeletePersonById())
		})
	})

	server.Router.Group("/role", func(group *rest.Router) {
		group.GET("", rest.HandleListRole())
		group.POST("", rest.HandleCreateRole())

		group.Group("/{role_id}", func(group *rest.Router) {
			group.GET("", rest.HandleGetRoleById())
			group.PUT("", rest.HandleUpdateRoleById())
			group.DELETE("", rest.HandleDeleteRoleById())
		})
	})

	server.Router.GET("/query", rest.HandleExecuteQuery())

	// Serve app
	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Listening and serving on: %s", addr)
		srvErr <- server.ListenAndServe(addr)
	}()

	// Wait for interruption.
	select {
	case err := <-srvErr:
		// Error when starting HTTP server.
		return err
	case <-ctx.Done():
		// Cleanup afer shutdown signalled
		log.Println("Shutdown signal received")

		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// if err := http.Shut.Shutdown(ctx); err != nil {
		// 	return err
		// }

		log.Println("Shutdown completed")
	}

	return nil
}
