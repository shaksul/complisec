from pathlib import Path

parts = []

parts.append("""package repo

import (
	\"context\"
	\"database/sql\"
	\"encoding/json\"
	\"fmt\"
	\"math\"
	\"strings\"
	\"time\"

	\"github.com/lib/pq\"
)

type TrainingRepo struct {
	db DBInterface
}

func NewTrainingRepo(db DBInterface) *TrainingRepo {
	return &TrainingRepo{db: db}
}
""")

parts.append("""// Materials ----------------------------------------------------------------

func (r *TrainingRepo) CreateMaterial(ctx context.Context, material Material) error {
	query := `
		INSERT INTO materials (
			id, tenant_id, title, description, uri, type, material_type,
			duration_minutes, tags, is_required, passing_score, attempts_limit,
			metadata, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14
		)`

	metadataJSON := marshalJSON(material.Metadata)

	_, err := r.db.ExecContext(
		ctx,
		query,
		material.ID,
		material.TenantID,
		material.Title,
		material.Description,
		material.URI,
		material.Type,
		material.MaterialType,
		material.DurationMinutes,
		pq.Array(material.Tags),
		material.IsRequired,
		material.PassingScore,
		material.AttemptsLimit,
		metadataJSON,
		material.CreatedBy,
	)

	return err
}

func (r *TrainingRepo) GetMaterialByID(ctx context.Context, id string) (*Material, error) {
	query := `
		SELECT id, tenant_id, title, description, uri, type, material_type,
			duration_minutes, tags, is_required, passing_score, attempts_limit,
			metadata, created_by, created_at, updated_at
		FROM materials
		WHERE id = $1`

	var material Material
	var tags pq.StringArray
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&material.ID,
		&material.TenantID,
		&material.Title,
		&material.Description,
		&material.URI,
		&material.Type,
		&material.MaterialType,
		&material.DurationMinutes,
		(*pq.StringArray)(&tags),
		&material.IsRequired,
		&material.PassingScore,
		&material.AttemptsLimit,
		&metadataJSON,
		&material.CreatedBy,
		&material.CreatedAt,
		&material.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	material.Tags = cloneStringSlice([]string(tags))
	material.Metadata = unmarshalJSONMap(metadataJSON)

	return &material, nil
}

func (r *TrainingRepo) ListMaterials(ctx context.Context, tenantID string, filters map[string]interface{}) ([]Material, error) {
	query := `
		SELECT id, tenant_id, title, description, uri, type, material_type,
			duration_minutes, tags, is_required, passing_score, attempts_limit,
			metadata, created_by, created_at, updated_at
		FROM materials
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	argIdx := 2

	if v, ok := filters[\"material_type\"]; ok {
		query += fmt.Sprintf(\" AND material_type = $%d\", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters[\"is_required\"]; ok {
		query += fmt.Sprintf(\" AND is_required = $%d\", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters[\"search\"]; ok {
		query += fmt.Sprintf(\" AND (LOWER(title) LIKE $%d OR LOWER(COALESCE(description, '')) LIKE $%d)\", argIdx, argIdx)
		args = append(args, patternForSearch(v))
		argIdx++
	}

	query += \" ORDER BY created_at DESC\"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var materials []Material
	for rs.Next() {
		var material Material
		var tags pq.StringArray
		var metadataJSON []byte

		if err := rs.Scan(
			&material.ID,
			&material.TenantID,
			&material.Title,
			&material.Description,
			&material.URI,
			&material.Type,
			&material.MaterialType,
			&material.DurationMinutes,
			(*pq.StringArray)(&tags),
			&material.IsRequired,
			&material.PassingScore,
			&material.AttemptsLimit,
			&metadataJSON,
			&material.CreatedBy,
			&material.CreatedAt,
			&material.UpdatedAt,
		); err != nil {
			return nil, err
		}

		material.Tags = cloneStringSlice([]string(tags))
		material.Metadata = unmarshalJSONMap(metadataJSON)
		materials = append(materials, material)
	}

	return materials, rs.Err()
}

func (r *TrainingRepo) UpdateMaterial(ctx context.Context, material Material) error {
	query := `
		UPDATE materials SET
			tenant_id = $2,
			title = $3,
			description = $4,
			uri = $5,
			type = $6,
			material_type = $7,
			duration_minutes = $8,
			tags = $9,
			is_required = $10,
			passing_score = $11,
			attempts_limit = $12,
			metadata = $13,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	metadataJSON := marshalJSON(material.Metadata)

	_, err := r.db.ExecContext(
		ctx,
		query,
		material.ID,
		material.TenantID,
		material.Title,
		material.Description,
		material.URI,
		material.Type,
		material.MaterialType,
		material.DurationMinutes,
		pq.Array(material.Tags),
		material.IsRequired,
		material.PassingScore,
		material.AttemptsLimit,
		metadataJSON,
	)
	return err
}

func (r *TrainingRepo) DeleteMaterial(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, \"DELETE FROM materials WHERE id = $1\", id)
	return err
}
""")

parts.append("""// Courses -------------------------------------------------------------------

func (r *TrainingRepo) CreateCourse(ctx context.Context, course TrainingCourse) error {
	query := `
		INSERT INTO training_courses (
			id, tenant_id, title, description, is_active, created_by
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.TenantID,
		course.Title,
		course.Description,
		course.IsActive,
		course.CreatedBy,
	)
	return err
}

func (r *TrainingRepo) GetCourseByID(ctx context.Context, id string) (*TrainingCourse, error) {
	query := `
		SELECT id, tenant_id, title, description, is_active, created_by, created_at, updated_at
		FROM training_courses
		WHERE id = $1`

	var course TrainingCourse
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.TenantID,
		&course.Title,
		&course.Description,
		&course.IsActive,
		&course.CreatedBy,
		&course.CreatedAt,
		&course.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *TrainingRepo) ListCourses(ctx context.Context, tenantID string, filters map[string]interface{}) ([]TrainingCourse, error) {
	query := `
		SELECT id, tenant_id, title, description, is_active, created_by, created_at, updated_at
		FROM training_courses
		WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	argIdx := 2

	if v, ok := filters[\"is_active\"]; ok {
		query += fmt.Sprintf(\" AND is_active = $%d\", argIdx)
		args = append(args, v)
		argIdx++
	}

	if v, ok := filters[\"search\"]; ok {
		query += fmt.Sprintf(\" AND LOWER(title) LIKE $%d\", argIdx)
		args = append(args, patternForSearch(v))
		argIdx++
	}

	query += \" ORDER BY created_at DESC\"

	rs, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var courses []TrainingCourse
	for rs.Next() {
		var course TrainingCourse
		if err := rs.Scan(
			&course.ID,
			&course.TenantID,
			&course.Title,
			&course.Description,
			&course.IsActive,
			&course.CreatedBy,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, rs.Err()
}

func (r *TrainingRepo) UpdateCourse(ctx context.Context, course TrainingCourse) error {
	query := `
		UPDATE training_courses SET
			tenant_id = $2,
			title = $3,
			description = $4,
			is_active = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.TenantID,
		course.Title,
		course.Description,
		course.IsActive,
	)
	return err
}

func (r *TrainingRepo) DeleteCourse(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, \"DELETE FROM training_courses WHERE id = $1\", id)
	return err
}

func (r *TrainingRepo) AddMaterialToCourse(ctx context.Context, cm CourseMaterial) error {
	query := `
		INSERT INTO course_materials (id, course_id, material_id, order_index, is_required)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (course_id, material_id) DO UPDATE
		SET order_index = EXCLUDED.order_index,
			is_required = EXCLUDED.is_required`

	_, err := r.db.ExecContext(ctx, query,
		cm.ID,
		cm.CourseID,
		cm.MaterialID,
		cm.OrderIndex,
		cm.IsRequired,
	)
	return err
}

func (r *TrainingRepo) RemoveMaterialFromCourse(ctx context.Context, courseID, materialID string) error {
	_, err := r.db.ExecContext(ctx, \"DELETE FROM course_materials WHERE course_id = $1 AND material_id = $2\", courseID, materialID)
	return err
}

func (r *TrainingRepo) GetCourseMaterials(ctx context.Context, courseID string) ([]CourseMaterial, error) {
	query := `
		SELECT id, course_id, material_id, order_index, is_required, created_at
		FROM course_materials
		WHERE course_id = $1
		ORDER BY order_index ASC, created_at ASC`

	rs, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var materials []CourseMaterial
	for rs.Next() {
		var cm CourseMaterial
		if err := rs.Scan(
			&cm.ID,
			&cm.CourseID,
			&cm.MaterialID,
			&cm.OrderIndex,
			&cm.IsRequired,
			&cm.CreatedAt,
		); err != nil {
			return nil, err
		}
		materials = append(materials, cm)
	}

	return materials, rs.Err()
}
""")

# Additional sections would continue here

content = "\n".join(parts)

Path(r"d:/CompliSec/apps/backend/internal/repo/training_repo.go").write_text(content, encoding="utf-8")
