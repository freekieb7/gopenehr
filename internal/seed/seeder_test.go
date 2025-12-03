package seed

// func TestSeeder(t *testing.T) {
// 	databaseURL := "postgres://gopenehr:gopenehrpass@localhost:5432/gopenehr?sslmode=disable"

// 	// Setup logger and database (mock or test instance)
// 	logger := telemetry.NewLogger("test", nil)
// 	db := database.New()
// 	if err := db.Connect(context.Background(), databaseURL); err != nil {
// 		t.Fatalf("Failed to connect to database: %v", err)
// 	}
// 	defer db.Close()

// 	seeder := NewSeeder(logger, db, "../../internal/seed/fixture")

// 	// Seed a small number of EHRs for testing
// 	seeder.Seed(10)

// 	// // Verify that EHRs and Compositions were created in the database
// 	// ehrCount, err := db.CountEHRs(context.Background())
// 	// if err != nil {
// 	// 	t.Fatalf("Failed to count EHRs: %v", err)
// 	// }
// 	// if ehrCount != 10 {
// 	// 	t.Errorf("Expected 10 EHRs, got %d", ehrCount)
// 	// }

// 	// compositionCount, err := db.CountCompositions(context.Background())
// 	// if err != nil {
// 	// 	t.Fatalf("Failed to count Compositions: %v", err)
// 	// }
// 	// if compositionCount == 0 {
// 	// 	t.Errorf("Expected some Compositions to be created, got %d", compositionCount)
// 	// }
// }
