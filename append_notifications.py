from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

notifications = """// Notifications ------------------------------------------------------------

func (r *TrainingRepo) CreateNotification(ctx context.Context, notification TrainingNotification) error {
	query := `
		INSERT INTO training_notifications (
			id, tenant_id, assignment_id, user_id, type, title, message,
			is_read, sent_at, read_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query,
		notification.ID,
		notification.TenantID,
		notification.AssignmentID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.IsRead,
		notification.SentAt,
		notification.ReadAt,
	)
	return err
}

func (r *TrainingRepo) GetUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]TrainingNotification, error) {
	query := `
		SELECT id, tenant_id, assignment_id, user_id, type, title, message,
			is_read, sent_at, read_at
		FROM training_notifications
		WHERE user_id = $1`

	args := []interface{}{userID}

	if unreadOnly {
		query += " AND is_read = false"
	}

	query += " ORDER BY sent_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var notifications []TrainingNotification
	for rs.Next() {
		notification, err := scanNotification(rs)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, *notification)
	}

	return notifications, rs.Err()
}

func (r *TrainingRepo) MarkNotificationAsRead(ctx context.Context, notificationID, userID string) error {
	query := `
		UPDATE training_notifications
		SET is_read = true, read_at = COALESCE(read_at, CURRENT_TIMESTAMP)
		WHERE id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, notificationID, userID)
	return err
}
"""

content += "\n\n" + notifications

path.write_text(content, encoding="utf-8")
