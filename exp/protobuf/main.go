package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	collectorEvent "example.com/protobuf/gen/collector/event/v1"
	resourcepb "example.com/protobuf/gen/resource/v1"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//
// ---------------------------------------------------
// Global CEL environment (CORRECT & SAFE VERSION)
// ---------------------------------------------------
//

var celEnv *cel.Env

func init() {
	env, err := cel.NewEnv(
		// Register protobuf types so CEL understands field structure
		cel.Types(&resourcepb.Resource{}),

		// IMPORTANT:
		// Use DynType for variables when passing protobuf messages.
		// This avoids ObjectType resolution failures.
		cel.Variable("resource", cel.DynType),
	)
	if err != nil {
		panic(err)
	}

	rs := resourcepb.Resource{}
	log.Println("CEL registered proto:",
		rs.ProtoReflect().Descriptor().FullName())

	celEnv = env
}

//
// ---------------------------------------------------
// Rule cache
// ---------------------------------------------------
//

type RuleEntry struct {
	Program cel.Program
	Rule    string
}

var (
	ruleCache sync.Map
	ruleGroup singleflight.Group
)

var httpClient = &http.Client{
	Timeout: 2 * time.Second,
}

func loadRule(schemaURL string) (cel.Program, error) {
	if v, ok := ruleCache.Load(schemaURL); ok {
		return v.(*RuleEntry).Program, nil
	}

	v, err, _ := ruleGroup.Do(schemaURL, func() (any, error) {
		resp, err := httpClient.Get(schemaURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(resp.Status)
		}

		var rule string
		if err := json.NewDecoder(resp.Body).Decode(&rule); err != nil {
			return nil, err
		}

		ast, issues := celEnv.Compile(rule)
		if issues != nil && issues.Err() != nil {
			return nil, issues.Err()
		}

		prg, err := celEnv.Program(
			ast,
			cel.EvalOptions(cel.OptOptimize),
		)
		if err != nil {
			return nil, err
		}

		entry := &RuleEntry{
			Program: prg,
			Rule:    rule,
		}

		ruleCache.Store(schemaURL, entry)
		log.Printf("📜 Loaded rule from %s", schemaURL)

		return prg, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(cel.Program), nil
}

//
// ---------------------------------------------------
// gRPC server
// ---------------------------------------------------
//

type server struct {
	collectorEvent.UnimplementedEventServiceServer
}

func (s *server) ExportEvents(
	ctx context.Context,
	req *collectorEvent.ExportEventsServiceRequest,
) (*collectorEvent.ExportEventsServiceResponse, error) {

	for _, re := range req.Data.ResourceEvents {
		if re.SchemaUrl == "" || re.Resource == nil {
			continue
		}

		prg, err := loadRule(re.SchemaUrl)
		if err != nil {
			log.Printf("❌ rule load failed (%s): %v", re.SchemaUrl, err)
			continue
		}

		out, _, err := prg.Eval(map[string]any{
			// Pass protobuf directly
			"resource": re.Resource,
		})
		if err != nil {
			log.Printf("❌ CEL eval failed: %v", err)
			continue
		}

		if out == types.True {
			log.Println("✅ Resource validation passed")
		} else {
			log.Println("❌ Resource validation failed")
		}
	}

	return &collectorEvent.ExportEventsServiceResponse{}, nil
}

//
// ---------------------------------------------------
// Main
// ---------------------------------------------------
//

func main() {
	// gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	collectorEvent.RegisterEventServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	go func() {
		log.Println("🚀 gRPC listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Rule endpoint
	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
		rule := `resource.attributes.exists(a, a.key == "test123" && a.value.string_value == "something")`
		jsonData, err := json.Marshal(rule)
		if err != nil {
			http.Error(w, "Failed to marshal rule", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	log.Println("📜 Rule server listening on :8080")
	panic(http.ListenAndServe(":8080", nil))
}
