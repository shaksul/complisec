package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	TenantID     string
	Email        string
	PasswordHash string
	FirstName    *string
	LastName     *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Printf("DEBUG: user_repo.Create inserting user tenant=%s email=%s", u.TenantID, u.Email)
	_, err = r.db.ExecContext(ctx, `
        INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, u.ID, u.TenantID, u.Email, string(hashedPassword), u.FirstName, u.LastName, u.IsActive)
	if err != nil {
		log.Printf("ERROR: user_repo.Create insert failed: %v", err)
	}
	return err
}

func (r *UserRepo) GetByEmail(ctx context.Context, tenantID, email string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1 AND email = $2
    `, tenantID, email)

	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("DEBUG: user_repo.GetByEmail no rows tenant=%s email=%s", tenantID, email)
			return nil, nil
		}
		log.Printf("ERROR: user_repo.GetByEmail query failed tenant=%s email=%s: %v", tenantID, email, err)
		return nil, err
	}
	log.Printf("DEBUG: user_repo.GetByEmail found user id=%s", u.ID)
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`, id)

	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByIDAndTenant(ctx context.Context, id, tenantID string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE id = $1 AND tenant_id = $2
	`, id, tenantID)

	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmailAndTenant(ctx context.Context, email, tenantID string) (*User, error) {
	row := r.db.QueryRow(`
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE email = $1 AND tenant_id = $2
	`, email, tenantID)

	var u User
	err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, tenantID string) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1 ORDER BY created_at DESC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) ListPaginated(ctx context.Context, tenantID string, page, pageSize int) ([]User, int64, error) {
	offset := (page - 1) * pageSize

	// Получаем общее количество записей
	var total int64
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE tenant_id = $1
	`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Получаем данные с пагинацией
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, tenantID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

// GetUsersByTenant retrieves all users for a tenant
func (r *UserRepo) GetUsersByTenant(ctx context.Context, tenantID string) ([]User, error) {
	return r.List(ctx, tenantID)
}

func (r *UserRepo) Update(ctx context.Context, u User) error {
	_, err := r.db.Exec(`
		UPDATE users SET first_name = $1, last_name = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND tenant_id = $5
	`, u.FirstName, u.LastName, u.IsActive, u.ID, u.TenantID)
	return err
}

func (r *UserRepo) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	log.Printf("DEBUG: GetUserRoles called for userID: %s", userID)

	rows, err := r.db.Query(`
		SELECT r.name FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		JOIN users u ON ur.user_id = u.id
		WHERE ur.user_id = $1 AND u.tenant_id = r.tenant_id
	`, userID)
	if err != nil {
		log.Printf("ERROR: GetUserRoles query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	roles := make([]string, 0)
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			log.Printf("ERROR: GetUserRoles scan failed: %v", err)
			return nil, err
		}
		log.Printf("DEBUG: GetUserRoles found role: %s", role)
		roles = append(roles, role)
	}

	log.Printf("DEBUG: GetUserRoles returning roles: %v", roles)
	return roles, nil
}

// GetUserPermissions получает все права пользователя через его роли
func (r *UserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT p.code FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		JOIN users u ON ur.user_id = u.id
		WHERE ur.user_id = $1 AND u.tenant_id = (
			SELECT tenant_id FROM roles WHERE id = ur.role_id
		)
		ORDER BY p.code
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]string, 0)
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (r *UserRepo) GetUserRoleIDs(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT r.id FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		JOIN users u ON ur.user_id = u.id
		WHERE ur.user_id = $1 AND u.tenant_id = r.tenant_id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleIDs := make([]string, 0)
	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		roleIDs = append(roleIDs, roleID)
	}
	return roleIDs, nil
}

func (r *UserRepo) SetUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	// Get user's tenant ID first
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Remove existing roles
	_, err = r.db.Exec("DELETE FROM user_roles WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	// Add new roles (only if they belong to the same tenant)
	for _, roleID := range roleIDs {
		_, err = r.db.Exec(`
			INSERT INTO user_roles (user_id, role_id) 
			SELECT $1, $2 
			WHERE EXISTS (
				SELECT 1 FROM roles r 
				WHERE r.id = $2 AND r.tenant_id = $3
			)
		`, userID, roleID, user.TenantID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepo) GetUserWithRoles(ctx context.Context, userID string) (*UserWithRoles, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	roles, err := r.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserWithRoles{
		User:  *user,
		Roles: roles,
	}, nil
}

func (r *UserRepo) GetUserWithRolesByTenant(ctx context.Context, userID, tenantID string) (*UserWithRoles, error) {
	user, err := r.GetByIDAndTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	roleIDs, err := r.GetUserRoleIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserWithRoles{
		User:  *user,
		Roles: roleIDs,
	}, nil
}

func (r *UserRepo) SearchUsers(ctx context.Context, tenantID string, search, role string, isActive *bool, sortBy, sortDir string, page, pageSize int) ([]User, int64, error) {
	offset := (page - 1) * pageSize

	// Build WHERE clause
	whereClause := "WHERE u.tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	if search != "" {
		whereClause += " AND (u.email ILIKE $" + fmt.Sprintf("%d", argIndex) + " OR u.first_name ILIKE $" + fmt.Sprintf("%d", argIndex+1) + " OR u.last_name ILIKE $" + fmt.Sprintf("%d", argIndex+2) + ")"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	if isActive != nil {
		whereClause += " AND u.is_active = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	// Build ORDER BY clause
	orderBy := "ORDER BY u.created_at DESC"
	if sortBy != "" {
		switch sortBy {
		case "email":
			orderBy = "ORDER BY u.email"
		case "first_name":
			orderBy = "ORDER BY u.first_name"
		case "last_name":
			orderBy = "ORDER BY u.last_name"
		case "created_at":
			orderBy = "ORDER BY u.created_at"
		case "updated_at":
			orderBy = "ORDER BY u.updated_at"
		}
		if sortDir == "asc" {
			orderBy += " ASC"
		} else {
			orderBy += " DESC"
		}
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM users u
		` + whereClause

	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT u.id, u.tenant_id, u.email, u.password_hash, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM users u
		` + whereClause + `
		` + orderBy + `
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *UserRepo) GetUserStats(ctx context.Context, userID string) (map[string]int, error) {
	// Get user's tenant ID first
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	stats := make(map[string]int)

	// Count documents
	var docCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM documents WHERE created_by = $1 AND tenant_id = $2
	`, userID, user.TenantID).Scan(&docCount)
	if err != nil {
		// Если таблица не существует, возвращаем 0
		stats["documents_count"] = 0
	} else {
		stats["documents_count"] = docCount
	}

	// Count risks
	var riskCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM risks WHERE owner_user_id = $1 AND tenant_id = $2
	`, userID, user.TenantID).Scan(&riskCount)
	if err != nil {
		stats["risks_count"] = 0
	} else {
		stats["risks_count"] = riskCount
	}

	// Count incidents
	var incidentCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM incidents WHERE reported_by = $1 AND tenant_id = $2
	`, userID, user.TenantID).Scan(&incidentCount)
	if err != nil {
		stats["incidents_count"] = 0
	} else {
		stats["incidents_count"] = incidentCount
	}

	// Count assets
	var assetCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM assets WHERE responsible_user_id = $1 AND tenant_id = $2
	`, userID, user.TenantID).Scan(&assetCount)
	if err != nil {
		stats["assets_count"] = 0
	} else {
		stats["assets_count"] = assetCount
	}

	// Count sessions
	var sessionCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM user_sessions WHERE user_id = $1 AND is_active = true
	`, userID).Scan(&sessionCount)
	if err != nil {
		stats["sessions_count"] = 0
	} else {
		stats["sessions_count"] = sessionCount
	}

	// Count successful logins
	var loginCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM login_attempts WHERE user_id = $1 AND success = true
	`, userID).Scan(&loginCount)
	if err != nil {
		stats["login_count"] = 0
	} else {
		stats["login_count"] = loginCount
	}

	// Calculate activity score based on recent activity
	var activityCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM user_activities WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
	`, userID).Scan(&activityCount)
	if err != nil {
		stats["activity_score"] = 0
	} else {
		// Simple scoring: more activity = higher score (max 100)
		activityScore := activityCount * 5
		if activityScore > 100 {
			activityScore = 100
		}
		stats["activity_score"] = activityScore
	}

	return stats, nil
}

type UserWithRoles struct {
	User
	Roles []string
}

// UserActivity represents user activity log
type UserActivity struct {
	ID          string
	UserID      string
	Action      string
	Description string
	IPAddress   *string
	UserAgent   *string
	CreatedAt   time.Time
	Metadata    map[string]interface{}
}

// GetUserActivity retrieves user activity with pagination
func (r *UserRepo) GetUserActivity(ctx context.Context, userID string, page, pageSize int) ([]UserActivity, int64, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM user_activities WHERE user_id = $1
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get activities with pagination
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, action, description, ip_address, user_agent, created_at, metadata
		FROM user_activities 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	activities := make([]UserActivity, 0)
	for rows.Next() {
		var activity UserActivity
		var metadataJSON sql.NullString
		err := rows.Scan(&activity.ID, &activity.UserID, &activity.Action, &activity.Description,
			&activity.IPAddress, &activity.UserAgent, &activity.CreatedAt, &metadataJSON)
		if err != nil {
			return nil, 0, err
		}

		// Parse metadata JSON if present
		if metadataJSON.Valid && metadataJSON.String != "" {
			// For now, we'll leave metadata as nil since we don't have a JSON parser
			// In a real implementation, you'd parse the JSON here
			activity.Metadata = make(map[string]interface{})
		}

		activities = append(activities, activity)
	}

	return activities, total, nil
}

// GetUserActivityStats retrieves aggregated activity statistics
func (r *UserRepo) GetUserActivityStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get daily activity for the last 30 days
	rows, err := r.db.QueryContext(ctx, `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM user_activities 
		WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, userID)
	if err != nil {
		stats["daily_activity"] = []map[string]interface{}{}
	} else {
		defer rows.Close()
		dailyActivity := make([]map[string]interface{}, 0)
		for rows.Next() {
			var date string
			var count int
			if err := rows.Scan(&date, &count); err == nil {
				dailyActivity = append(dailyActivity, map[string]interface{}{
					"date":  date,
					"count": count,
				})
			}
		}
		stats["daily_activity"] = dailyActivity
	}

	// Get top actions
	rows, err = r.db.QueryContext(ctx, `
		SELECT action, COUNT(*) as count
		FROM user_activities 
		WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY action
		ORDER BY count DESC
		LIMIT 10
	`, userID)
	if err != nil {
		stats["top_actions"] = []map[string]interface{}{}
	} else {
		defer rows.Close()
		topActions := make([]map[string]interface{}, 0)
		for rows.Next() {
			var action string
			var count int
			if err := rows.Scan(&action, &count); err == nil {
				topActions = append(topActions, map[string]interface{}{
					"action": action,
					"count":  count,
				})
			}
		}
		stats["top_actions"] = topActions
	}

	// Get login history (assuming we have a login_attempts table)
	rows, err = r.db.QueryContext(ctx, `
		SELECT ip_address, user_agent, created_at, success
		FROM login_attempts 
		WHERE user_id = $1 
		ORDER BY created_at DESC
		LIMIT 10
	`, userID)
	if err != nil {
		stats["login_history"] = []map[string]interface{}{}
	} else {
		defer rows.Close()
		loginHistory := make([]map[string]interface{}, 0)
		for rows.Next() {
			var ipAddress, userAgent string
			var createdAt time.Time
			var success bool
			if err := rows.Scan(&ipAddress, &userAgent, &createdAt, &success); err == nil {
				loginHistory = append(loginHistory, map[string]interface{}{
					"ip_address": ipAddress,
					"user_agent": userAgent,
					"created_at": createdAt,
					"success":    success,
				})
			}
		}
		stats["login_history"] = loginHistory
	}

	return stats, nil
}

// LogUserActivity logs a user activity
func (r *UserRepo) LogUserActivity(ctx context.Context, userID, action, description, ipAddress, userAgent string, metadata map[string]interface{}) error {
	// Get user's tenant ID first
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		log.Printf("ERROR: LogUserActivity GetByID: %v", err)
		return err
	}
	if user == nil {
		log.Printf("WARN: LogUserActivity user not found: %s", userID)
		return errors.New("user not found")
	}

	// Convert metadata to JSON string
	var metadataJSON string
	if len(metadata) > 0 {
		// For now, we'll store as simple JSON string
		// In a real implementation, you'd use a proper JSON library
		metadataJSON = "{}" // Placeholder
	}

	// Insert activity into database
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO user_activities (user_id, tenant_id, action, description, ip_address, user_agent, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, user.TenantID, action, description, ipAddress, userAgent, metadataJSON)

	if err != nil {
		log.Printf("ERROR: LogUserActivity insert failed: %v", err)
		return err
	}

	log.Printf("User Activity logged: %s - %s: %s (IP: %s)", userID, action, description, ipAddress)
	return nil
}

// LogLoginAttempt logs a login attempt
func (r *UserRepo) LogLoginAttempt(ctx context.Context, userID, tenantID, email, ipAddress, userAgent string, success bool, failureReason string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO login_attempts (user_id, tenant_id, email, ip_address, user_agent, success, failure_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, tenantID, email, ipAddress, userAgent, success, failureReason)

	if err != nil {
		log.Printf("ERROR: LogLoginAttempt insert failed: %v", err)
		return err
	}

	log.Printf("Login attempt logged: %s - %s (IP: %s, Success: %v)", email, failureReason, ipAddress, success)
	return nil
}
