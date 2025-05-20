package repository

import (
	model "kondangin-backend/internal/model"

	"gorm.io/gorm"
)

type PackageRepository interface {
	FindAll() ([]model.Package, error)
	FindByID(id uint) (*model.Package, error)
	FindByCode(code string) (*model.Package, error)
	Create(pkg *model.Package) error
	Update(pkg *model.Package) error
	Delete(id uint) error
}

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) PackageRepository {
	return &packageRepository{db}
}

func (r *packageRepository) FindAll() ([]model.Package, error) {
	var packages []model.Package
	err := r.db.Find(&packages).Error
	return packages, err
}

func (r *packageRepository) FindByID(id uint) (*model.Package, error) {
	var pkg model.Package
	err := r.db.First(&pkg, id).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) FindByCode(code string) (*model.Package, error) {
	var pkg model.Package
	err := r.db.Where("code = ?", code).First(&pkg).Error
	return &pkg, err
}

func (r *packageRepository) Create(pkg *model.Package) error {
	return r.db.Create(pkg).Error
}

func (r *packageRepository) Update(pkg *model.Package) error {
	return r.db.Save(pkg).Error
}

func (r *packageRepository) Delete(id uint) error {
	return r.db.Delete(&model.Package{}, id).Error
}
