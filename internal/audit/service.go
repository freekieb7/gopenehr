package audit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/audit"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/google/uuid"
)

type Service struct {
	Logger *telemetry.Logger
	DB     *database.Database
}

func NewService(logger *telemetry.Logger, db *database.Database) *Service {
	return &Service{
		Logger: logger,
		DB:     db,
	}
}

// Pagination types for audit log listing
type ListEventsRequest struct {
	PageSize int
	Token    string
}

type ListEventsResponse struct {
	Events    []audit.Event
	NextToken utils.Optional[string]
	PrevToken utils.Optional[string]
}

type ListEventsCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Direction string    `json:"direction"`
}

func (s *Service) ListEventsPaginated(ctx context.Context, req ListEventsRequest) (ListEventsResponse, error) {
	// Set default page size if invalid
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 25
	}

	cursor, err := decodeListEventsCursor(req.Token)
	if err != nil {
		s.Logger.Warn("Failed to decode page token", "token", req.Token, "error", err)
		return ListEventsResponse{}, errors.New("invalid page token")
	}

	// Extract direction from cursor, default to "next" for first page
	direction := "next"
	if cursor.E {
		direction = cursor.V.Direction
	}

	query := `SELECT id, actor_id, actor_type, resource, action, success, ip_address, user_agent, details, created_at FROM audit.tbl_audit_log `
	args := []any{}
	argIdx := 1

	order := "DESC"
	cmp := "<"

	if direction == "prev" {
		order = "ASC"
		cmp = ">"
	}

	if cursor.E {
		query += fmt.Sprintf("WHERE (created_at, id) %s ($%d, $%d) ", cmp, argIdx, argIdx+1)
		args = append(args, cursor.V.CreatedAt, cursor.V.ID)
		argIdx += 2
	}

	query += fmt.Sprintf("ORDER BY created_at %s, id %s ", order, order)
	query += fmt.Sprintf("LIMIT $%d", argIdx)
	args = append(args, pageSize+1) // Fetch one extra to detect more pages

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		s.Logger.Warn("Failed to query audit logs", "error", err)
		return ListEventsResponse{}, err
	}
	defer rows.Close()

	events := []audit.Event{}
	for rows.Next() {
		var event audit.Event
		if err := rows.Scan(&event.ID, &event.ActorID, &event.ActorType, &event.Resource, &event.Action, &event.Success, &event.IPAddress, &event.UserAgent, &event.Details, &event.CreatedAt); err != nil {
			s.Logger.Warn("Failed to scan event", "error", err)
			return ListEventsResponse{}, err
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		return ListEventsResponse{
			Events: events,
		}, nil
	}

	// Check if we have more results than requested (indicates more pages)
	hasMore := len(events) > pageSize
	if hasMore {
		events = events[:pageSize] // Remove the extra record
	}

	// If fetching previous page, reverse results to maintain order
	if direction == "prev" {
		for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
			events[i], events[j] = events[j], events[i]
		}
	}

	// Generate next/prev tokens based on actual data availability
	var nextCursor, prevCursor utils.Optional[string]

	if len(events) > 0 {
		// Generate next token (for older records)
		showNext := (direction == "next" && hasMore) || (direction == "prev")
		if showNext {
			if nextCursorStr, err := encodeListEventsCursor(ListEventsCursor{
				ID:        events[len(events)-1].ID,
				CreatedAt: events[len(events)-1].CreatedAt,
				Direction: "next",
			}); err == nil {
				nextCursor = utils.Some(nextCursorStr)
			}
		}

		// Generate prev token (for newer records)
		showPrev := (direction == "prev" && hasMore) || (direction == "next" && req.Token != "")
		if showPrev {
			if prevCursorStr, err := encodeListEventsCursor(ListEventsCursor{
				ID:        events[0].ID,
				CreatedAt: events[0].CreatedAt,
				Direction: "prev",
			}); err == nil {
				prevCursor = utils.Some(prevCursorStr)
			}
		}
	}

	return ListEventsResponse{
		Events:    events,
		NextToken: nextCursor,
		PrevToken: prevCursor,
	}, nil
}

func encodeListEventsCursor(cursor ListEventsCursor) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func decodeListEventsCursor(token string) (utils.Optional[ListEventsCursor], error) {
	if token == "" {
		return utils.None[ListEventsCursor](), nil
	}
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return utils.None[ListEventsCursor](), err
	}
	var cursor ListEventsCursor
	err = json.Unmarshal(data, &cursor)
	return utils.Some(cursor), err
}
