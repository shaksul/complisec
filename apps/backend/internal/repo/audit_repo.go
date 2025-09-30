package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          string
	TenantID    string
	ActorID     *string
	Action      string
	Entity      string
	EntityID    *string
	PayloadJSON *string
	Timestamp   time.Time
}

type AuditRepo struct {
	db *DB
}

func NewAuditRepo(db *DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Log(ctx context.Context, log AuditLog) error {
	_, err := r.db.Exec(`
		INSERT INTO audit_log (id, tenant_id, actor_id, action, entity, entity_id, payload_json, ts)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, log.ID, log.TenantID, log.ActorID, log.Action, log.Entity, log.EntityID, log.PayloadJSON, log.Timestamp)
	return err
}

func (r *AuditRepo) LogAction(ctx context.Context, tenantID, actorID, action, entity string, entityID *string, payload interface{}) error {
	var payloadJSON *string
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		jsonStr := string(jsonBytes)
		payloadJSON = &jsonStr
	}

	var actorIDPtr *string
	if actorID != "" {
		actorIDPtr = &actorID
	}

	log := AuditLog{
		ID:          generateUUID(),
		TenantID:    tenantID,
		ActorID:     actorIDPtr,
		Action:      action,
		Entity:      entity,
		EntityID:    entityID,
		PayloadJSON: payloadJSON,
		Timestamp:   time.Now(),
	}

	return r.Log(ctx, log)
}

func generateUUID() string {
	return uuid.New().String()
}

// GetAuditLogs получает журнал аудита с фильтрацией
func (r *AuditRepo) GetAuditLogs(ctx context.Context, tenantID string, limit, offset int, actorID, action, entity *string) ([]AuditLog, error) {
	query := `
		SELECT id, tenant_id, actor_id, action, entity, entity_id, payload_json, ts
		FROM audit_log 
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argIndex := 2

	if actorID != nil {
		query += fmt.Sprintf(" AND actor_id = $%d", argIndex)
		args = append(args, *actorID)
		argIndex++
	}

	if action != nil {
		query += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, *action)
		argIndex++
	}

	if entity != nil {
		query += fmt.Sprintf(" AND entity = $%d", argIndex)
		args = append(args, *entity)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY ts DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(&log.ID, &log.TenantID, &log.ActorID, &log.Action, &log.Entity, &log.EntityID, &log.PayloadJSON, &log.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
