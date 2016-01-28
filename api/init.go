package api

import (
	// "github.com/gin-gonic/gin"
	// "net/http"
	"xsbAdmin/models"
)

var (
	user_manager     = models.NewUserManager()
	imei_manager     = models.NewImeiManager()
	company_manager  = models.NewCompanyManager()
	category_manager = models.NewCategoryManager()
	group_manager    = models.NewGroupManager()
	book_manager     = models.NewBookManager()

	default_password = "111"

	max_bound_imei = 3

	default_avatar = "default_avatar.png"

	Avatar_dir         = "./files/avatars"
	Category_cover_dir = "./files/categoryCovers"
	Book_dir           = "./files/books"
)
