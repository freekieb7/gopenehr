package webhook

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/freekieb7/gopenehr/internal/database"
	"github.com/freekieb7/gopenehr/internal/telemetry"
	"github.com/freekieb7/gopenehr/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrInvalidEventType = errors.New("invalid event type")
)

type Service struct {
	Logger *telemetry.Logger
	DB     *database.Database
}

func NewService(logger *telemetry.Logger, db *database.Database) *Service {
	return &Service{Logger: logger, DB: db}
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
