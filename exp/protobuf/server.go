package grpc

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/google/cel-go/cel"
// 	"github.com/google/cel-go/common/types"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// type EventService struct {
// 	collectorEvent.UnimplementedEventServiceServer
// 	client *http.Client
// }

// func NewEventService() collectorEvent.EventServiceServer {
// 	return &EventService{
// 		client: &http.Client{},
// 	}
// }

// func NewServer() *grpc.Server {
// 	var opts []grpc.ServerOption
// 	grpcServer := grpc.NewServer(opts...)
// 	collectorEvent.RegisterEventServiceServer(grpcServer, NewEventService())
// 	reflection.Register(grpcServer) // Register reflection service on gRPC server.
// 	return grpcServer
// }

// // ExportEvents implements v1.EventServiceServer.
// func (e *EventService) ExportEvents(ctx context.Context, r *collectorEvent.ExportEventsServiceRequest) (*collectorEvent.ExportEventsServiceResponse, error) {
// 	for _, resourceEvent := range r.Data.ResourceEvents {
// 		schemaUrl := resourceEvent.SchemaUrl
// 		if schemaUrl == "" {
// 			continue
// 		}

// 		// Fetch schema rules based on resourceEvent.SchemaUrl
// 		resp, err := e.client.Get(resourceEvent.SchemaUrl)
// 		if err != nil {
// 			fmt.Println("❌ Error fetching schema:", err)
// 			continue
// 		}
// 		defer resp.Body.Close()
// 		if resp.StatusCode != http.StatusOK {
// 			fmt.Println("❌ Non-OK HTTP status:", resp.Status)
// 			continue
// 		}

// 		var schemaRules []string
// 		err = json.NewDecoder(resp.Body).Decode(&schemaRules)
// 		if err != nil {
// 			fmt.Println("❌ Error decoding schema response:", err)
// 			continue
// 		}

// 		if err := validateResource(resourceEvent, schemaRules); err != nil {
// 			fmt.Println("❌ Validation error:", err)
// 		} else {
// 			fmt.Println("✅ Resource valid")
// 		}
// 	}
// 	return &collectorEvent.ExportEventsServiceResponse{
// 		PartialSuccess: &collectorEvent.ExportEventsPartialSuccess{RejectedDataPoints: 0},
// 	}, nil
// }

// func validateResource(r *event.ResourceEvent, rules []string) error {
// 	env, err := cel.NewEnv(
// 		cel.Variable("resource", cel.DynType),
// 		cel.Variable("schema_url", cel.StringType),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	resourceMap := resourceToMap(r.Resource)
// 	vars := map[string]any{
// 		"resource":   resourceMap,
// 		"schema_url": r.SchemaUrl,
// 	}

// 	for _, rule := range rules {
// 		ast, issues := env.Compile(rule)
// 		if issues != nil && issues.Err() != nil {
// 			return fmt.Errorf("invalid rule: %v", issues.Err())
// 		}

// 		prg, _ := env.Program(ast)
// 		out, _, _ := prg.Eval(vars)

// 		if out != types.True {
// 			return fmt.Errorf("validation failed for rule: %s", rule)
// 		}
// 	}

// 	return nil
// }

// func resourceToMap(r *resource.Resource) map[string]any {
// 	m := make(map[string]any)

// 	for _, attr := range r.Attributes {
// 		if attr == nil || attr.Key == "" || attr.Value == nil {
// 			continue
// 		}

// 		k := attr.Key
// 		switch v := attr.Value.Value.(type) {
// 		case *common.AnyValue_StringValue:
// 			m[k] = v.StringValue
// 		case *common.AnyValue_IntValue:
// 			m[k] = v.IntValue
// 		case *common.AnyValue_BoolValue:
// 			m[k] = v.BoolValue
// 		case *common.AnyValue_DoubleValue:
// 			m[k] = v.DoubleValue
// 			// add other cases if needed (array, kvlist, bytes)
// 		}
// 	}
// 	return m
// }
