from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

helpers = """// Helpers ------------------------------------------------------------------

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanAssignment(scanner rowScanner) (*TrainingAssignment, error) {
	var assignment TrainingAssignment
	var materialID sql.NullString
	var courseID sql.NullString
	var dueAt sql.NullTime
	var completedAt sql.NullTime
	var assignedBy sql.NullString
	var lastAccessed sql.NullTime
	var reminderSent sql.NullTime
	var metadataJSON []byte

	if err := scanner.Scan(
		&assignment.ID,
		&assignment.TenantID,
		&materialID,
		&courseID,
		&assignment.UserID,
		&assignment.Status,
		&dueAt,
		&completedAt,
		&assignedBy,
		&assignment.Priority,
		&assignment.ProgressPercentage,
		&assignment.TimeSpentMinutes,
		&lastAccessed,
		&reminderSent,
		&metadataJSON,
		&assignment.CreatedAt,
	); err != nil {
		return nil, err
	}

	assignment.MaterialID = stringPointer(materialID)
	assignment.CourseID = stringPointer(courseID)
	assignment.DueAt = timePointer(dueAt)
	assignment.CompletedAt = timePointer(completedAt)
	assignment.AssignedBy = stringPointer(assignedBy)
	assignment.LastAccessedAt = timePointer(lastAccessed)
	assignment.ReminderSentAt = timePointer(reminderSent)
	assignment.Metadata = unmarshalJSONMap(metadataJSON)

	return &assignment, nil
}

func scanCertificate(scanner rowScanner) (*Certificate, error) {
	var certificate Certificate
	var materialID sql.NullString
	var courseID sql.NullString
	var expiresAt sql.NullTime
	var metadataJSON []byte

	if err := scanner.Scan(
		&certificate.ID,
		&certificate.TenantID,
		&certificate.AssignmentID,
		&certificate.UserID,
		&materialID,
		&courseID,
		&certificate.CertificateNumber,
		&certificate.IssuedAt,
		&expiresAt,
		&certificate.IsValid,
		&metadataJSON,
		&certificate.CreatedAt,
	); err != nil {
		return nil, err
	}

	certificate.MaterialID = stringPointer(materialID)
	certificate.CourseID = stringPointer(courseID)
	certificate.ExpiresAt = timePointer(expiresAt)
	certificate.Metadata = unmarshalJSONMap(metadataJSON)

	return &certificate, nil
}

func scanNotification(scanner rowScanner) (*TrainingNotification, error) {
	var notification TrainingNotification
	var readAt sql.NullTime

	if err := scanner.Scan(
		&notification.ID,
		&notification.TenantID,
		&notification.AssignmentID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.IsRead,
		&notification.SentAt,
		&readAt,
	); err != nil {
		return nil, err
	}

	notification.ReadAt = timePointer(readAt)
	return &notification, nil
}

func (r *TrainingRepo) getCertificate(ctx context.Context, predicate string, arg interface{}) (*Certificate, error) {
	query := fmt.Sprintf(
		"SELECT id, tenant_id, assignment_id, user_id, material_id, course_id, certificate_number, issued_at, expires_at, is_valid, metadata, created_at FROM certificates WHERE %s",
		predicate,
	)

	row := r.db.QueryRowContext(ctx, query, arg)
	certificate, err := scanCertificate(row)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}

func marshalJSON(data map[string]any) []byte {
	if len(data) == 0 {
		return nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return b
}

func unmarshalJSONMap(data []byte) map[string]any {
	if len(data) == 0 {
		return nil
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}

func patternForSearch(value interface{}) string {
	text := strings.ToLower(fmt.Sprint(value))
	return \"%\" + text + \"%\"
}

func cloneStringSlice(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func stringPointer(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func timePointer(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	t := nt.Time
	return &t
}

func intPointer(ni sql.NullInt64) *int {
	if !ni.Valid {
		return nil
	}
	v := int(ni.Int64)
	return &v
}
"""

content += "\n\n" + helpers

path.write_text(content, encoding="utf-8")
