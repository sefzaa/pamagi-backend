package route

import (
	"pamagi/api/controller"
	"pamagi/repository"
	"pamagi/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewCategoryRouter(db *gorm.DB, protectedRouter *gin.RouterGroup) {
	categoryRepository := repository.NewCategoryRepository(db)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepository)
	categoryController := &controller.CategoryController{
		CategoryUsecase: categoryUsecase,
	}

	protectedRouter.POST("/categories", categoryController.CreateCategory)
	protectedRouter.GET("/categories", categoryController.GetCategories)
	protectedRouter.PUT("/categories/:id", categoryController.UpdateCategory)
	protectedRouter.DELETE("/categories/:id", categoryController.DeleteCategory)
}