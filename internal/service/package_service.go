package service

import (
	model "kondangin-backend/internal/model"
	"kondangin-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type PackageService interface {
	GetAll(c *gin.Context) ([]model.Package, error)
	GetByID(c *gin.Context, id uint) (*model.Package, error)
	GetByCode(c *gin.Context, code string) (*model.Package, error)
	Create(c *gin.Context, pkg *model.Package) error
	Update(c *gin.Context, pkg *model.Package) error
	Delete(c *gin.Context, id uint) error
}

type packageService struct {
	repo repository.PackageRepository
}

func NewPackageService(repo repository.PackageRepository) PackageService {
	return &packageService{repo}
}

func (s *packageService) GetAll(c *gin.Context) ([]model.Package, error) {
	// Bisa tambahkan log atau user audit dari context di sini
	return s.repo.FindAll()
}

func (s *packageService) GetByID(c *gin.Context, id uint) (*model.Package, error) {
	return s.repo.FindByID(id)
}

func (s *packageService) GetByCode(c *gin.Context, code string) (*model.Package, error) {
	return s.repo.FindByCode(code)
}

func (s *packageService) Create(c *gin.Context, pkg *model.Package) error {
	// contoh ambil user dari context: user := c.MustGet("user").(model.User)
	return s.repo.Create(pkg)
}

func (s *packageService) Update(c *gin.Context, pkg *model.Package) error {
	return s.repo.Update(pkg)
}

func (s *packageService) Delete(c *gin.Context, id uint) error {
	return s.repo.Delete(id)
}
