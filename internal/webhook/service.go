package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidEventType = errors.New("invalid event type")
)

type Service struct {
	Logger *slog.Logger
	DB     *database.Database
}

func NewService(logger *slog.Logger, db *database.Database) Service {
	return Service{Logger: logger, DB: db}
}

type EventType string

const (
	EventTypeEHRCreated EventType = "ehr.created"
	EventTypeEHRDeleted EventType = "ehr.deleted"

	EventTypeEHRStatusUpdated EventType = "ehr_status.updated"

	EventTypeCompositionCreated EventType = "composition.created"
	EventTypeCompositionUpdated EventType = "composition.updated"
	EventTypeCompositionDeleted EventType = "composition.deleted"

	EventTypeDirectoryCreated EventType = "directory.created"
	EventTypeDirectoryUpdated EventType = "directory.updated"
	EventTypeDirectoryDeleted EventType = "directory.deleted"

	EventTypePersonCreated EventType = "person.created"
	EventTypePersonUpdated EventType = "person.updated"
	EventTypePersonDeleted EventType = "person.deleted"

	EventTypeAgentCreated EventType = "agent.created"
	EventTypeAgentUpdated EventType = "agent.updated"
	EventTypeAgentDeleted EventType = "agent.deleted"

	EventTypeGroupCreated EventType = "group.created"
	EventTypeGroupUpdated EventType = "group.updated"
	EventTypeGroupDeleted EventType = "group.deleted"

	EventTypeOrganisationCreated EventType = "organisation.created"
	EventTypeOrganisationUpdated EventType = "organisation.updated"
	EventTypeOrganisationDeleted EventType = "organisation.deleted"

	EventTypeRoleCreated EventType = "role.created"
	EventTypeRoleUpdated EventType = "role.updated"
	EventTypeRoleDeleted EventType = "role.deleted"

	EventTypeQueryExecuted EventType = "query.executed"
	EventTypeQueryStored   EventType = "query.stored"
)

var EventTypes = map[EventType]string{
	EventTypeEHRCreated:          "EHR Created",
	EventTypeEHRDeleted:          "EHR Deleted",
	EventTypeEHRStatusUpdated:    "EHR Status Updated",
	EventTypeCompositionCreated:  "Composition Created",
	EventTypeCompositionDeleted:  "Composition Deleted",
	EventTypePersonCreated:       "Person Created",
	EventTypePersonUpdated:       "Person Updated",
	EventTypePersonDeleted:       "Person Deleted",
	EventTypeAgentCreated:        "Agent Created",
	EventTypeAgentUpdated:        "Agent Updated",
	EventTypeAgentDeleted:        "Agent Deleted",
	EventTypeGroupCreated:        "Group Created",
	EventTypeGroupUpdated:        "Group Updated",
	EventTypeGroupDeleted:        "Group Deleted",
	EventTypeOrganisationCreated: "Organisation Created",
	EventTypeOrganisationUpdated: "Organisation Updated",
	EventTypeOrganisationDeleted: "Organisation Deleted",
	EventTypeRoleCreated:         "Role Created",
	EventTypeRoleUpdated:         "Role Updated",
	EventTypeRoleDeleted:         "Role Deleted",
	EventTypeQueryExecuted:       "Query Executed",
	EventTypeQueryStored:         "Query Stored",
}

func IsValidEventType(event EventType) bool {
	_, exists := EventTypes[event]
	return exists
}

type Subscription struct {
	ID              uuid.UUID
	URL             string
	Secret          string
	EventTypes      []EventType
	IsActive        bool
	LastDeliveredAt utils.Optional[time.Time]
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *Service) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	rows, err := s.DB.Query(ctx, `
		SELECT id, url, secret, event_types, is_active, created_at
		FROM webhook.tbl_subscription
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query webhook subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		var subscription Subscription
		if err := rows.Scan(&subscription.ID, &subscription.URL, &subscription.Secret, &subscription.EventTypes, &subscription.IsActive, &subscription.CreatedAt); err != nil {
			return subscriptions, fmt.Errorf("failed to scan webhook subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}
	if rows.Err() != nil {
		return subscriptions, fmt.Errorf("error iterating webhook subscriptions: %w", err)
	}

	return subscriptions, nil
}

func (s *Service) ExistsSubscription(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	row := s.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM webhook.tbl_subscription WHERE id = $1)`, id)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check existence of webhook subscription: %w", err)
	}
	return exists, nil
}

func (s *Service) ExistsSubscriptionWithURL(ctx context.Context, url string) (bool, error) {
	var exists bool
	row := s.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM webhook.tbl_subscription WHERE url = $1)`, url)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check existence of webhook subscription: %w", err)
	}
	return exists, nil
}

func (s *Service) Subscribe(ctx context.Context, url string, eventTypes []string) (Subscription, error) {
	if url == "" {
		return Subscription{}, fmt.Errorf("subscription URL is required")
	}

	eventTypesConverted := make([]EventType, len(eventTypes))
	for i, et := range eventTypes {
		eventType := EventType(et)
		if !IsValidEventType(eventType) {
			return Subscription{}, fmt.Errorf("%w: %s", ErrInvalidEventType, et)
		}
		eventTypesConverted[i] = eventType
	}

	randomSecret, err := utils.GenerateRandomString(32)
	if err != nil {
		return Subscription{}, fmt.Errorf("failed to generate random string: %w", err)
	}

	subscription := Subscription{
		URL:        url,
		Secret:     randomSecret,
		EventTypes: eventTypesConverted,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}

	row := s.DB.QueryRow(ctx, `INSERT INTO webhook.tbl_subscription (url, secret, event_types, is_active)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, subscription.URL, subscription.Secret, subscription.EventTypes, subscription.IsActive)
	if err := row.Scan(&subscription.ID); err != nil {
		return Subscription{}, fmt.Errorf("failed to insert webhook subscription: %w", err)
	}

	return subscription, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, subscriptionID uuid.UUID, eventTypes utils.Optional[[]string]) error {
	if eventTypes.E {
		eventTypesConverted := make([]EventType, len(eventTypes.V))
		for i, et := range eventTypes.V {
			eventType := EventType(et)
			if !IsValidEventType(eventType) {
				return fmt.Errorf("%w: %s", ErrInvalidEventType, et)
			}
			eventTypesConverted[i] = eventType
		}

		_, err := s.DB.Exec(ctx, `UPDATE webhook.tbl_subscription SET event_types = $1, updated_at = NOW() WHERE id = $2`, eventTypesConverted, subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to update webhook subscription: %w", err)
		}
	}

	return nil
}

func (s *Service) Unsubscribe(ctx context.Context, subscriptionID uuid.UUID) error {
	_, err := s.DB.Exec(ctx, `DELETE FROM webhook.tbl_subscription WHERE id = $1`, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook subscription: %w", err)
	}
	return nil
}

func (s *Service) RegisterEvent(ctx context.Context, eventType EventType, data map[string]any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook event data: %w", err)
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, database.ErrTxClosed) {
			s.Logger.ErrorContext(ctx, "failed to rollback transaction", "error", err)
		}
	}()

	var eventID uuid.UUID
	row := tx.QueryRow(ctx, `INSERT INTO webhook.tbl_event (type, payload) VALUES ($1, $2) RETURNING id`, string(eventType), dataBytes)
	err = row.Scan(&eventID)
	if err != nil {
		return fmt.Errorf("failed to insert webhook event: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO webhook.tbl_delivery (event_id, subscription_id, attempt_count, status, created_at)
		SELECT e.id, s.id, 0, 'pending', NOW()
		FROM webhook.tbl_event e
		CROSS JOIN webhook.tbl_subscription s
		WHERE s.is_active = TRUE AND e.id = $1 AND $2 = ANY(s.event_types)
	`, eventID, string(eventType))
	if err != nil {
		return fmt.Errorf("failed to insert webhook deliveries: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
