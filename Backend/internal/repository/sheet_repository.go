package repository

import (
	"dojo/internal/models"

	"gorm.io/gorm"
)

type SheetRepository struct {
	db *gorm.DB
}

func NewSheetRepository(db *gorm.DB) *SheetRepository {
	return &SheetRepository{db: db}
}

// Create creates a new problem sheet
func (r *SheetRepository) Create(sheet *models.ProblemSheet) error {
	return r.db.Create(sheet).Error
}

// FindByID retrieves a sheet by ID with its problems
func (r *SheetRepository) FindByID(id string) (*models.ProblemSheet, error) {
	var sheet models.ProblemSheet
	err := r.db.Preload("SheetProblems.Problem").Where("id = ?", id).First(&sheet).Error
	return &sheet, err
}

// FindByUserID retrieves all sheets for a user
func (r *SheetRepository) FindByUserID(userID string) ([]models.ProblemSheet, error) {
	var sheets []models.ProblemSheet
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&sheets).Error
	return sheets, err
}

// FindPublicSheets retrieves all public sheets
func (r *SheetRepository) FindPublicSheets() ([]models.ProblemSheet, error) {
	var sheets []models.ProblemSheet
	err := r.db.Where("is_public = ?", true).Order("created_at DESC").Find(&sheets).Error
	return sheets, err
}

// Update updates a sheet
func (r *SheetRepository) Update(sheet *models.ProblemSheet) error {
	return r.db.Save(sheet).Error
}

// Delete deletes a sheet
func (r *SheetRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.ProblemSheet{}).Error
}

// AddProblem adds a problem to a sheet
func (r *SheetRepository) AddProblem(sheetProblem *models.SheetProblem) error {
	return r.db.Create(sheetProblem).Error
}

// RemoveProblem removes a problem from a sheet
func (r *SheetRepository) RemoveProblem(sheetID, problemID string) error {
	return r.db.Where("sheet_id = ? AND problem_id = ?", sheetID, problemID).Delete(&models.SheetProblem{}).Error
}

// UpdateSheetProblem updates a problem within a sheet (solved status, notes)
func (r *SheetRepository) UpdateSheetProblem(sheetProblem *models.SheetProblem) error {
	return r.db.Save(sheetProblem).Error
}

// FindSheetProblem finds a specific problem in a sheet
func (r *SheetRepository) FindSheetProblem(sheetID, problemID string) (*models.SheetProblem, error) {
	var sheetProblem models.SheetProblem
	err := r.db.Preload("Problem").Where("sheet_id = ? AND problem_id = ?", sheetID, problemID).First(&sheetProblem).Error
	return &sheetProblem, err
}

// ExistsProblemInSheet checks if a problem already exists in a sheet
func (r *SheetRepository) ExistsProblemInSheet(sheetID, problemID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.SheetProblem{}).Where("sheet_id = ? AND problem_id = ?", sheetID, problemID).Count(&count).Error
	return count > 0, err
}

// GetNextPosition gets the next position for a problem in a sheet
func (r *SheetRepository) GetNextPosition(sheetID string) (int, error) {
	var maxPosition int
	err := r.db.Model(&models.SheetProblem{}).Where("sheet_id = ?", sheetID).Select("COALESCE(MAX(position), 0)").Scan(&maxPosition).Error
	return maxPosition + 1, err
}
