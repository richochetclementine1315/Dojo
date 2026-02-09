package repository

import (
	"dojo/internal/models"
	"strings"

	"gorm.io/gorm"
)

type ProblemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) *ProblemRepository {
	return &ProblemRepository{db: db}
}

// Create creates a new problem
func (r *ProblemRepository) Create(problem *models.Problem) error {
	return r.db.Create(problem).Error
}

// FindByID retrieves a problem by ID
func (r *ProblemRepository) FindByID(id string) (*models.Problem, error) {
	var problem models.Problem
	err := r.db.First(&problem, "id = ?", id).Error
	return &problem, err
}

// FindAll retrieves all problems with optional filters
func (r *ProblemRepository) FindAll(filters map[string]interface{}, page, limit int) ([]models.Problem, int64, error) {
	var problems []models.Problem
	var total int64

	query := r.db.Model(&models.Problem{})

	// Apply filters
	if difficulty, ok := filters["difficulty"].(string); ok && difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	if platform, ok := filters["platform"].(string); ok && platform != "" {
		query = query.Where("platform = ?", platform)
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern)
	}

	if tags, ok := filters["tags"].([]string); ok && len(tags) > 0 {
		query = query.Where("tags @> ?", tags)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&problems).Error; err != nil {
		return nil, 0, err
	}

	return problems, total, nil
}

// Update updates an existing problem
func (r *ProblemRepository) Update(problem *models.Problem) error {
	return r.db.Save(problem).Error
}

// Delete deletes a problem by ID
func (r *ProblemRepository) Delete(id string) error {
	return r.db.Delete(&models.Problem{}, "id = ?", id).Error
}

// ExistsByURL checks if a problem exists with the given URL
func (r *ProblemRepository) ExistsByURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Problem{}).Where("problem_url = ?", url).Count(&count).Error
	return count > 0, err
}

// ExistsByPlatformID checks if a problem exists with the given platform and platform ID
func (r *ProblemRepository) ExistsByPlatformID(platform, platformID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Problem{}).
		Where("platform = ? AND platform_problem_id = ?", platform, platformID).
		Count(&count).Error
	return count > 0, err
}

// GetDB returns the underlying database connection
func (r *ProblemRepository) GetDB() *gorm.DB {
	return r.db
}
