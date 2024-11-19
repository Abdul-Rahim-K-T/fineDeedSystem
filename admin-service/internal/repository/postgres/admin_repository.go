package postgres

import (
	"fineDeedSystem/admin-service/internal/models"
	// proto "fineDeedSystem/proto/employer"

	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(admin models.Admin) error
	FindAdminByID(id uint) (models.Admin, error)
	FindAdminByAdminname(adminname string) (models.Admin, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) CreateAdmin(admin models.Admin) error {
	return r.db.Create(&admin).Error
}

func (r *adminRepository) FindAdminByID(id uint) (models.Admin, error) {
	var admin models.Admin
	err := r.db.First(&admin, id).Error
	return admin, err
}

func (r *adminRepository) FindAdminByAdminname(adminname string) (models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("adminname = ?", adminname).First(&admin).Error
	return admin, err
}

// func (r *adminRepository) GetAllEmployers() ([]*proto.Employer, error) {
// 	var employers []*proto.Employer
// 	err := r.db.Find(&employers).Error
// 	return employers, err
// }
