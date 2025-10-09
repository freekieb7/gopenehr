package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freekieb7/gopenehr/internal/grpc"
	"github.com/freekieb7/gopenehr/internal/http"
	"github.com/freekieb7/gopenehr/internal/storage"
)

func main() {
	ctx := context.Background()

	if err := Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func Run(ctx context.Context) error {
	var (
		apiPort  = flag.Int("port", 8080, "Port for the OpenEHR API")
		grpcPort = flag.Int("grpc-port", 8081, "Port for the gRPC server")
	)
	flag.Parse()

	// Add gracefull shutdown support by listening for interuptions
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	addr := fmt.Sprintf("0.0.0.0:%d", *apiPort)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Initialize database
	db := storage.NewDatabase()
	if err := db.Connect(ctx, "postgres://postgres:example@localhost:5432/postgres?sslmode=disable"); err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// HTTP Handlers
	openEHRHandler := http.NewOpenEHRHandler(&db)
	queryHandler := http.NewQueryHandler(logger, &db)

	// Setup router and server
	router := http.NewRouter()

	// Add routes
	router.Group("/schemas", func(group *http.Router) {
		schemaHandler := http.NewSchemaHandler()
		group.GET("", schemaHandler.GetSchemas)
		group.GET("/{name}", schemaHandler.GetSchema)
	})

	// OpenEHR API routes
	router.Group("/openehr/v1", func(group *http.Router) {
		group.GET("/", openEHRHandler.ServerInfo)

		group.Group("/ehr", func(group *http.Router) {
			group.GET("", openEHRHandler.ListEhr)
			group.POST("", openEHRHandler.CreateEhr)

			group.Group("/{ehr_id}", func(group *http.Router) {
				group.GET("", openEHRHandler.GetEhrById)
				group.DELETE("", openEHRHandler.DeleteEhrById)

				group.Group("/ehr_status", func(group *http.Router) {
					group.GET("", openEHRHandler.GetEhrStatusById)
					group.PUT("", openEHRHandler.UpdateEhrStatusById)
				})

				group.Group("/composition", func(group *http.Router) {
					group.GET("", openEHRHandler.ListComposition)
					group.POST("", openEHRHandler.CreateComposition)

					group.Group("/{composition_id}", func(group *http.Router) {
						group.GET("", openEHRHandler.GetCompositionById)
						group.PUT("", openEHRHandler.UpdateCompositionById)
						group.DELETE("", openEHRHandler.DeleteCompositionById)
					})
				})

				group.Group("/folder", func(group *http.Router) {
					group.GET("", openEHRHandler.ListFolder)
					group.POST("", openEHRHandler.CreateFolder)

					group.Group("/{folder_id}", func(group *http.Router) {
						group.GET("", openEHRHandler.GetFolderById)
						group.PUT("", openEHRHandler.UpdateFolderById)
						group.DELETE("", openEHRHandler.DeleteFolderById)
					})
				})

				// todo contribution (just not sure what the difference is with demographics)
			})
		})

		group.Group("/template", func(group *http.Router) {
			group.GET("", openEHRHandler.ListTemplates)
			group.POST("", openEHRHandler.CreateTemplate)

			group.Group("/{template_id}", func(group *http.Router) {
				group.GET("", openEHRHandler.GetTemplateById)
				group.PUT("", openEHRHandler.UpdateTemplateById)
				group.DELETE("", openEHRHandler.DeleteTemplateById)
			})
		})

		group.Group("/agent", func(group *http.Router) {
			group.GET("", openEHRHandler.ListAgent)
			group.POST("", openEHRHandler.CreateAgent)

			group.Group("/{agent_id}", func(group *http.Router) {
				group.GET("", openEHRHandler.GetAgentById)
				group.PUT("", openEHRHandler.UpdateAgentById)
				group.DELETE("", openEHRHandler.DeleteAgentById)
			})
		})

		group.Group("/group", func(group *http.Router) {
			group.GET("", openEHRHandler.ListGroup)
			group.POST("", openEHRHandler.CreateGroup)

			group.Group("/{group_id}", func(group *http.Router) {
				group.GET("", openEHRHandler.GetGroupById)
				group.PUT("", openEHRHandler.UpdateGroupById)
				group.DELETE("", openEHRHandler.DeleteGroupById)
			})
		})

		group.Group("/organisation", func(group *http.Router) {
			group.GET("", openEHRHandler.ListOrganisation)
			group.POST("", openEHRHandler.CreateOrganisation)

			group.Group("/{organisation_id}", func(group *http.Router) {
				group.GET("", openEHRHandler.GetOrganisationById)
				group.PUT("", openEHRHandler.UpdateOrganisationById)
				group.DELETE("", openEHRHandler.DeleteOrganisationById)
			})
		})

		group.Group("/person", func(group *http.Router) {
			group.GET("", openEHRHandler.ListPerson)
			group.POST("", openEHRHandler.CreatePerson)

			group.Group("/{person_id}", func(group *http.Router) {
				group.GET("", openEHRHandler.GetPersonById)
				group.PUT("", openEHRHandler.UpdatePersonById)
				group.DELETE("", openEHRHandler.DeletePersonById)
			})
		})

		group.Group("/role", func(group *http.Router) {
			group.GET("", openEHRHandler.ListRole)
			group.POST("", openEHRHandler.CreateRole)

			group.Group("/{role_id}", func(group *http.Router) {
				group.GET("", openEHRHandler.GetRoleById)
				group.PUT("", openEHRHandler.UpdateRoleById)
				group.DELETE("", openEHRHandler.DeleteRoleById)
			})
		})

		group.Group("/query", func(group *http.Router) {
			group.POST("", queryHandler.ExecuteQuery)
			group.POST("/prepared", queryHandler.CreatePreparedTable)
			group.POST("/prepared/{name}/sync", queryHandler.SyncPreparedTable)
		})
	})

	server := http.NewServer(logger, router)

	// Serve app
	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Listening and serving on: %s", addr)
		srvErr <- server.ListenAndServe(addr)
	}()

	// Run grpc server
	grpcServer := grpc.NewServer()
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *grpcPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Printf("gRPC server listening on %s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			srvErr <- err
		}
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
