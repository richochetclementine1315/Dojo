package repository

import (
	"dojo/internal/models"
	"time"

	"gorm.io/gorm"
)

type ContestRepository struct {
	db *gorm.DB
}

func NewContestRepository(db *gorm.DB) *ContestRepository {
	return &ContestRepository{db: db}
}

// FindAll retrieves all contests withfilters
func (r *ContestRepository) FindAll(filters map[string]interface{}, page, limit int) ([]models.Contest, int64, error) {
	var contests []models.Contest
	var total int64

	query := r.db.Model(&models.Contest{})

	// Apply filters
	platform, ok := filters["platform"]
	if ok {
		query = query.Where("platform=?", platform)
	}
	upcoming, ok := filters["upcoming"]
	if ok && upcoming.(bool) {
		query = query.Where("start_time>?", time.Now())
	}
	ongoing, ok := filters["ongoing"]
	if ok && ongoing.(bool) {
		query = query.Where("start_time<=? AND end_time>=?", time.Now(), time.Now())
	}
	// Total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	//Pagination jobzz again :)
	offset := (page - 1) * limit
	err = query.Order("start_time  ASC").Offset(offset).Limit(limit).Find(&contests).Error
	if err != nil {
		return nil, 0, err
	}
	return contests, total, nil
}

// FindByID retrieves the contests by ID
func (r *ContestRepository) FindByID(id string) (*models.Contest, error) {
	var contest models.Contest
	err := r.db.Where("id=?", id).First(&contest).Error
	if err != nil {
		return nil, err
	}
	return &contest, nil
}

// CreareRemainder createes a new contest remainder
func (r *ContestRepository) CreateReminder(reminder *models.ContestReminder) error {
	return r.db.Create(reminder).Error
}

// FindReminderByID retrieves the remainder by ID
func (r *ContestRepository) FindReminderByID(id string) (*models.ContestReminder, error) {
	var reminder models.ContestReminder
	err := r.db.Preload("Contest").Where("id=?", id).First(&reminder).Error
	if err != nil {
		return nil, err
	}
	return &reminder, nil
}

// Deletereminder deletes a contest reminder
func (r *ContestRepository) DeleteReminder(id string) error {
	return r.db.Where("id=?", id).Delete(&models.ContestReminder{}).Error
}

// FindReminderByuserID retrieves the reminder by userID

func (r *ContestRepository) FindReminderByuserID(userID string) ([]models.ContestReminder, error) {
	var reminders []models.ContestReminder
	err := r.db.Where("user_id=?", userID).Preload("Contest").Find(&reminders).Error
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

// Existsreminder checks if a reminder exists for a user and contest
func (r *ContestRepository) ExistsReminder(userID, contestID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ContestReminder{}).Where("user_id=? AND contest_id=?", userID, contestID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
