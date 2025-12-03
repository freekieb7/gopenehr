package seed

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"runtime"
	"sync"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/openehr"
	"github.com/freekieb7/gopenehr/internal/openehr/rm"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/google/uuid"
)

type Seeder struct {
	Logger     *telemetry.Logger
	DB         *database.Database
	FixtureDir string
}

func NewSeeder(logger *telemetry.Logger, db *database.Database, fixtureDir string) *Seeder {
	return &Seeder{
		Logger:     logger,
		DB:         db,
		FixtureDir: fixtureDir,
	}
}

func (s *Seeder) Seed(count int) {
	var wg sync.WaitGroup
	workerCount := runtime.GOMAXPROCS(0)

	s.Logger.Info("Starting seeding process", "total_count", count, "workers", workerCount)
	for range workerCount {
		wg.Go(func() {
			s.seedEHRs(context.Background(), count)
		})
	}
	wg.Wait()
}

func (s *Seeder) seedEHRs(ctx context.Context, count int) {
	data, err := os.ReadFile(s.FixtureDir + "/t0016_rapportage.v1.json")
	if err != nil {
		panic(err)
	}

	var composition rm.COMPOSITION
	err = json.Unmarshal(data, &composition)
	if err != nil {
		panic(err)
	}

	openehrService := openehr.NewService(s.Logger, s.DB)

	s.Logger.Info("Starting seeding EHRs", "count", count)

	var compositionsCreated int
	var compositionsUpdated int

	for i := range count {
		compositionsToCreate := rand.Intn(10)
		compositionsToUpdate := rand.Intn(10)

		ehrID := uuid.New()
		_, err := openehrService.CreateEHR(ctx, ehrID, openehr.NewEHRStatus(uuid.New()))
		if err != nil {
			s.Logger.Error("Failed to create EHR", "error", err)
			continue
		}

		for range compositionsToCreate {
			s.RandomizeRapportage(&composition)
			newComposition, err := openehrService.CreateComposition(ctx, ehrID, composition)
			if err != nil {
				s.Logger.Error("Failed to create Composition", "error", err)
				continue
			}

			for range compositionsToUpdate {
				s.RandomizeRapportage(&composition)
				newComposition, err = openehrService.UpdateComposition(ctx, ehrID, newComposition, composition)
				if err != nil {
					s.Logger.Error("Failed to create Composition", "error", err)
					continue
				}
			}
		}

		compositionsCreated += compositionsToCreate
		compositionsUpdated += compositionsToUpdate

		if i%100 == 0 {
			s.Logger.Info("Seeded EHRs", "count", i, "compositions_created", compositionsCreated, "compositions_updated", compositionsUpdated)
		}
	}
	s.Logger.Info("Seeding complete", "total", count, "compositions_created", compositionsCreated, "compositions_updated", compositionsUpdated)
}

func (s *Seeder) RandomizeRapportage(composition *rm.COMPOSITION) {
	composition.Category.DefiningCode.CodeString = RandFromSlice([]string{"at0031", "at0032", "at0033", "at0034", "at0035", "at0036", "at0051"})
	composition.Context.V.OtherContext.V.ITEM_TREE().Items.V[0].ELEMENT().Value.V.DV_CODED_TEXT().DefiningCode.CodeString = RandFromSlice([]string{"at0011", "at0012", "at0013", "at0014", "at0015"})
	composition.Context.V.OtherContext.V.ITEM_TREE().Items.V[1].ELEMENT().Value.V.DV_CODED_TEXT().DefiningCode.CodeString = RandFromSlice([]string{"at0021", "at0022", "at0023", "at0024"})
	composition.Context.V.OtherContext.V.ITEM_TREE().Items.V[2].ELEMENT().Value.V.DV_CODED_TEXT().DefiningCode.CodeString = RandFromSlice([]string{"at0031", "at0032", "at0033", "at0034", "at0035", "at0036", "at0051"})
	composition.Context.V.OtherContext.V.ITEM_TREE().Items.V[3].ELEMENT().Value.V.DV_CODED_TEXT().DefiningCode.CodeString = RandFromSlice([]string{"at0041", "at0042", "at0043", "at0044"})
	composition.Context.V.OtherContext.V.ITEM_TREE().Items.V[4].ELEMENT().Value.V.DV_EHR_URI().Value = "ehr://ehr_id/value"
	composition.Content.V[0].EVALUATION().Data.ITEM_TREE().Items.V[0].ELEMENT().Value.V.DV_TEXT().Value = "1235"
}

func RandFromSlice[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}
