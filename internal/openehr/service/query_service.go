package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/freekieb7/gopenehr/internal/audit"
	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/openehr/aql"
	"github.com/google/uuid"
)

var (
	ErrQueryNotFound      = fmt.Errorf("AQL query not found")
	ErrQueryAlreadyExists = fmt.Errorf("AQL query with the given name already exists")
)

type StoredQuery struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Query   string    `json:"q"`
	Type    string    `json:"type"`
	Saved   time.Time `json:"saved"`
}

type QueryService struct {
	Logger       *slog.Logger
	DB           *database.Database
	AuditService *audit.Service
}

func NewQueryService(logger *slog.Logger, db *database.Database, auditService *audit.Service) QueryService {
	return QueryService{
		Logger:       logger,
		DB:           db,
		AuditService: auditService,
	}
}

func (s *QueryService) QueryAndCopyTo(ctx context.Context, w io.Writer, aqlQuery string, aqlParams map[string]any) error {
	if aqlParams == nil {
		aqlParams = make(map[string]any)
	}

	sqlQuery, _, err := aql.ToSQL(aqlQuery, aqlParams)
	if err != nil {
		s.Logger.Error("internal error", "error", err)
		return err
	}

	s.Logger.DebugContext(ctx, "query error", "error", err, "aql", aqlQuery, "sql", strings.ReplaceAll(strings.ReplaceAll(sqlQuery, "\n", " "), "\t", " "))

	rows, err := s.DB.Query(ctx, sqlQuery)
	if err != nil {
		s.Logger.ErrorContext(ctx, "query error", "error", err, "aql", aqlQuery, "sql", strings.ReplaceAll(strings.ReplaceAll(sqlQuery, "\n", " "), "\t", " "))
		return err
	}

	// Stream results as JSON array
	_, _ = w.Write([]byte(`{"rows":[`))

	first := true
	for rows.Next() {
		var jsonData []byte
		if err := rows.Scan(&jsonData); err != nil {
			s.Logger.Error("scan error", "error", err)
			continue
		}

		if !first {
			_, _ = w.Write([]byte(","))
		}
		_, _ = w.Write(jsonData)
		first = false

		// Flush each row so client receives data progressively
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	_, _ = w.Write([]byte("]}"))
	return nil
}

func (s *QueryService) ListStoredQueries(ctx context.Context, filterName string) ([]StoredQuery, error) {
	var query strings.Builder
	var args []any

	query.WriteString(`
		SELECT COALESCE(jsonb_agg(
			jsonb_build_object(
				'name', name,
				'type', 'AQL', 
				'version', version,
				'saved', to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MSTZH:TZM'),
				'q', query
			)), '[]'::jsonb) as queries
		FROM openehr.tbl_query
	`)

	if filterName != "" {
		namespaceSeperator := strings.LastIndex(filterName, "::")
		if namespaceSeperator != -1 {
			query.WriteString(` WHERE name = $1`)
		} else {
			query.WriteString(` WHERE name LIKE '%' || $1`)
		}
		args = append(args, filterName)
	}

	var queries []StoredQuery
	if err := s.DB.QueryRow(ctx, query.String(), args...).Scan(&queries); err != nil {
		return nil, fmt.Errorf("error querying stored AQL queries: %w", err)
	}

	return queries, nil
}

func (s *QueryService) GetQueryByName(ctx context.Context, name string, filterVersion string) (StoredQuery, error) {
	var query strings.Builder
	var args []any

	query.WriteString(`SELECT name, version, query, created_at FROM openehr.tbl_query WHERE name = $1 `)
	args = append(args, name)

	if filterVersion != "" {
		query.WriteString(`AND version = $2 `)
		args = append(args, filterVersion)
	}

	query.WriteString(`ORDER BY created_at DESC LIMIT 1`)

	var storedQuery StoredQuery
	if err := s.DB.QueryRow(ctx, query.String(), args...).Scan(&storedQuery.Name, &storedQuery.Version, &storedQuery.Query, &storedQuery.Saved); err != nil {
		if err == database.ErrNoRows {
			return StoredQuery{}, ErrQueryNotFound
		}
		return StoredQuery{}, fmt.Errorf("error retrieving AQL query by name: %w", err)
	}

	storedQuery.Type = "AQL"
	return storedQuery, nil
}

func (s *QueryService) StoreQuery(ctx context.Context, name, version, aqlQuery string) error {
	// Store the new query
	_, err := s.DB.Exec(ctx, `INSERT INTO openehr.tbl_query (id, name, version, query) VALUES ($1, $2, $3, $4) ON CONFLICT (name, version) DO UPDATE SET query = EXCLUDED.query`,
		uuid.New(),
		name,
		version,
		aqlQuery,
	)
	if err != nil {
		return fmt.Errorf("error storing AQL query: %w", err)
	}

	return nil
}
