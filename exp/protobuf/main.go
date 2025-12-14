package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	collectorEvent "example.com/protobuf/gen/collector/event/v1"
	resource "example.com/protobuf/gen/resource/v1"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	collectorEvent.UnimplementedEventServiceServer
}

func (s *server) ExportEvents(ctx context.Context, r *collectorEvent.ExportEventsServiceRequest) (*collectorEvent.ExportEventsServiceResponse, error) {
	for _, resourceEvent := range r.Data.ResourceEvents {
		schemaUrl := resourceEvent.SchemaUrl
		if schemaUrl == "" {
			continue
		}

		// Fetch schema rules from the provided URL
		resp, err := http.Get(schemaUrl)
		if err != nil {
			log.Println("❌ Error fetching schema:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("❌ Non-OK HTTP status:", resp.Status)
			continue
		}

		// Decode the CEL rule
		var celRule string
		err = json.NewDecoder(resp.Body).Decode(&celRule)
		if err != nil {
			fmt.Println("❌ Error decoding schema response:", err)
			continue
		}

		// Create CEL environment with resource variable
		env, err := cel.NewEnv(
			cel.Variable("resource", cel.DynType),
		)
		if err != nil {
			fmt.Println("❌ Error creating CEL environment:", err)
			continue
		}

		// Compile the CEL rule
		ast, issues := env.Compile(celRule)
		if issues != nil && issues.Err() != nil {
			log.Printf("❌ Invalid rule: %v", issues.Err())
			continue
		}

		// Create a program from the AST
		prg, err := env.Program(ast)
		if err != nil {
			log.Printf("❌ Error creating program: %v", err)
			continue
		}

		// Convert resource to a map for CEL evaluation
		resourceMap := convertResourceToMap(resourceEvent.Resource)

		// Evaluate the rule
		out, _, err := prg.Eval(map[string]any{
			"resource": resourceMap,
		})
		if err != nil {
			log.Printf("❌ Error evaluating rule: %v", err)
			continue
		}

		// Check if validation passed
		if out == types.True {
			fmt.Println("✅ Resource validation passed")
		} else {
			fmt.Printf("❌ Resource validation failed for rule: %s\n", celRule)
		}
	}

	return &collectorEvent.ExportEventsServiceResponse{}, nil
}

// Helper function to convert protobuf resource to a map structure
func convertResourceToMap(resource *resource.Resource) map[string]any {
	if resource == nil {
		return map[string]any{}
	}

	attributes := make([]map[string]any, len(resource.Attributes))
	for i, attr := range resource.Attributes {
		attributes[i] = map[string]any{
			"key": attr.Key,
			"value": map[string]any{
				"string_value": attr.Value.GetStringValue(),
			},
		}
	}

	return map[string]any{
		"attributes": attributes,
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the EventService server
	collectorEvent.RegisterEventServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer) // Register reflection service on gRPC server.

	// Start serving
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`"resource.attributes[0].key == 'test123' && resource.attributes[0].value.string_value == 'something'"`))
	})
	panic(http.ListenAndServe(":8080", nil))
}
