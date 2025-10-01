package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freekieb7/gopenehr/database"
	"github.com/freekieb7/gopenehr/rest"
	restv2 "github.com/freekieb7/gopenehr/rest/v2"
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Initialize database connection
	pg := database.NewPostgres()
	if err := pg.Connect(ctx, "postgres://gopenehr:gopenehr@localhost:5432/gopenehr?sslmode=disable"); err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer pg.Close()

	// Initialize database
	db := database.NewDatabase(&pg)

	// Setup router and server
	router := rest.NewRouter()

	// Add routes
	router.Group("/openehr/v1", func(group *rest.Router) {
		router.GET("/", rest.HandleServerInfo())

		router.Group("/ehr", func(group *rest.Router) {
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

		router.Group("/template", func(group *rest.Router) {
			group.GET("", rest.HandleListTemplates())
			group.POST("", rest.HandleCreateTemplate())

			group.Group("/{template_id}", func(group *rest.Router) {
				group.GET("", rest.HandleGetTemplateById())
				group.PUT("", rest.HandleUpdateTemplateById())
				group.DELETE("", rest.HandleDeleteTemplateById())
			})
		})

		router.Group("/agent", func(group *rest.Router) {
			group.GET("", rest.HandleListAgent())
			group.POST("", rest.HandleCreateAgent())

			group.Group("/{agent_id}", func(group *rest.Router) {
				group.GET("", rest.HandleGetAgentById())
				group.PUT("", rest.HandleUpdateAgentById())
				group.DELETE("", rest.HandleDeleteAgentById())
			})
		})

		router.Group("/group", func(group *rest.Router) {
			group.GET("", rest.HandleListGroup())
			group.POST("", rest.HandleCreateGroup())

			group.Group("/{group_id}", func(group *rest.Router) {
				group.GET("", rest.HandleGetGroupById())
				group.PUT("", rest.HandleUpdateGroupById())
				group.DELETE("", rest.HandleDeleteGroupById())
			})
		})

		router.Group("/organisation", func(group *rest.Router) {
			group.GET("", rest.HandleListOrganisation())
			group.POST("", rest.HandleCreateOrganisation())

			group.Group("/{organisation_id}", func(group *rest.Router) {
				group.GET("", rest.HandleGetOrganisationById())
				group.PUT("", rest.HandleUpdateOrganisationById())
				group.DELETE("", rest.HandleDeleteOrganisationById())
			})
		})

		router.Group("/person", func(group *rest.Router) {
			group.GET("", rest.HandleListPerson())
			group.POST("", rest.HandleCreatePerson())

			group.Group("/{person_id}", func(group *rest.Router) {
				group.GET("", rest.HandleGetPersonById())
				group.PUT("", rest.HandleUpdatePersonById())
				group.DELETE("", rest.HandleDeletePersonById())
			})
		})

		router.Group("/role", func(group *rest.Router) {
			group.GET("", rest.HandleListRole())
			group.POST("", rest.HandleCreateRole())

			group.Group("/{role_id}", func(group *rest.Router) {
				group.GET("", rest.HandleGetRoleById())
				group.PUT("", rest.HandleUpdateRoleById())
				group.DELETE("", rest.HandleDeleteRoleById())
			})
		})

		router.GET("/query", rest.HandleExecuteQuery())
	})

	router.Group("/chp/v1", func(group *rest.Router) {
		group.Group("/users", func(group *rest.Router) {
			group.GET("", restv2.HandleListUsers())
			group.POST("", restv2.HandleCreateUser(&db))

			group.Group("/{user_id}", func(group *rest.Router) {
				group.GET("", restv2.HandleGetUserById())
				group.PATCH("", restv2.HandleUpdateUserById())
				group.DELETE("", restv2.HandleDeleteUserById())
			})
		})

		group.Group("/folders", func(group *rest.Router) {
			group.GET("", restv2.HandleListFolders())
			group.POST("", restv2.HandleCreateFolder())

			group.Group("/{folder_id}", func(group *rest.Router) {
				group.GET("", restv2.HandleGetFolderById())
				group.PATCH("", restv2.HandleUpdateFolderById())
				group.DELETE("", restv2.HandleDeleteFolderById())
			})
		})

		group.Group("/documents", func(group *rest.Router) {
			group.GET("", restv2.HandleListDocuments())
			group.POST("", restv2.HandleCreateDocument())

			group.Group("/{document_id}", func(group *rest.Router) {
				group.GET("", restv2.HandleGetDocumentById())
				group.PATCH("", restv2.HandleUpdateDocumentById())
				group.DELETE("", restv2.HandleDeleteDocumentById())

				group.Group("/revisions", func(group *rest.Router) {
					group.GET("", restv2.HandleListDocumentRevisions())
					group.Group("/{revision_id}", func(group *rest.Router) {
						group.GET("", restv2.HandleGetDocumentRevisionById())
						group.PATCH("", restv2.HandleUpdateDocumentRevisionById())
						group.DELETE("", restv2.HandleDeleteDocumentRevisionById())
					})
				})
			})
		})

		group.Group("/audit-logs", func(group *rest.Router) {
			group.GET("", restv2.HandleGetAuditLog())
		})

		group.Group("/webhooks", func(group *rest.Router) {
			group.GET("", restv2.HandleListWebhooks())
			group.POST("", restv2.HandleCreateWebhook())

			group.Group("/{webhook_id}", func(group *rest.Router) {
				group.GET("", restv2.HandleGetWebhookById())
				group.PATCH("", restv2.HandleUpdateWebhookById())
				group.DELETE("", restv2.HandleDeleteWebhookById())
			})
		})

		group.Group("/calendar", func(group *rest.Router) {
			group.GET("", restv2.HandleGetCalendar())

			group.Group("/events", func(group *rest.Router) {
				group.GET("", restv2.HandleListCalendarEvents())
				group.POST("", restv2.HandleCreateCalendarEvent())

				group.Group("/{event_id}", func(group *rest.Router) {
					group.GET("", restv2.HandleGetCalendarEventById())
					group.PATCH("", restv2.HandleUpdateCalendarEventById())
					group.DELETE("", restv2.HandleDeleteCalendarEventById())
				})
			})
		})
	})

	server := rest.NewServer(logger, router)

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
