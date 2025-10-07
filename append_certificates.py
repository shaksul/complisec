from pathlib import Path

path = Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go")
content = path.read_text(encoding="utf-8")

certificates = """// Certificates -------------------------------------------------------------

func (r *TrainingRepo) CreateCertificate(ctx context.Context, certificate Certificate) error {
	query := `
		INSERT INTO certificates (
			id, tenant_id, assignment_id, user_id, material_id, course_id,
			certificate_number, issued_at, expires_at, is_valid, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	metadataJSON := marshalJSON(certificate.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		certificate.ID,
		certificate.TenantID,
		certificate.AssignmentID,
		certificate.UserID,
		certificate.MaterialID,
		certificate.CourseID,
		certificate.CertificateNumber,
		certificate.IssuedAt,
		certificate.ExpiresAt,
		certificate.IsValid,
		metadataJSON,
	)
	return err
}

func (r *TrainingRepo) GetCertificateByID(ctx context.Context, id string) (*Certificate, error) {
	return r.getCertificate(ctx, "id = $1", id)
}

func (r *TrainingRepo) GetCertificateByNumber(ctx context.Context, certificateNumber string) (*Certificate, error) {
	return r.getCertificate(ctx, "certificate_number = $1", certificateNumber)
}

func (r *TrainingRepo) GetUserCertificates(ctx context.Context, userID string, filters map[string]interface{}) ([]Certificate, error) {
	query := `
		SELECT id, tenant_id, assignment_id, user_id, material_id, course_id,
			certificate_number, issued_at, expires_at, is_valid, metadata, created_at
		FROM certificates
		WHERE user_id = $1`

	args := []interface{}{userID}
	argIdx := 2

	if v, ok := filters["tenant_id"]; ok {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["course_id"]; ok {
		query += fmt.Sprintf(" AND course_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["material_id"]; ok {
		query += fmt.Sprintf(" AND material_id = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters["is_valid"]; ok {
		query += fmt.Sprintf(" AND is_valid = $%d", argIdx)
		args = append(args, v)
		argIdx++
	}

	query += " ORDER BY issued_at DESC"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var certificates []Certificate
	for rs.Next() {
		certificate, err := scanCertificate(rs)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, *certificate)
	}

	return certificates, rs.Err()
}
"""

content += "\n\n" + certificates

path.write_text(content, encoding="utf-8")
