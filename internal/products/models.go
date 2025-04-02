package products

// import (
// 	"impx/internal/db"
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type MarketplaceProduct struct {
// 	*gorm.Model
// 	Id           uuid.UUID
// 	Name         string  `gorm:"column:name"`
// 	Description  string  `gorm:"column:description"`
// 	Category     string  `gorm:"column:category"`
// 	Mrp          float64 `gorm:"column:mrp"`
// 	Discount     float64 `gorm:"column:discount"`
// 	SellingPrice float64 `gorm:"column:selling_price"`
// 	Color        string  `gorm:"column:color"`
// 	Ratings      int64   `gorm:"column:ratings"`
// 	NumReviews   int64   `gorm:"column:num_reviews"`

// 	CreatedAt time.Time `gorm:"column:created_at"`
// 	UpdatedAt time.Time `gorm:"column:updated_at"`

// 	ArticleCode  string `gorm:"column:article_code"`
// 	InternalRank int64  `gorm:"column:internal_rank"`
// 	Material     string `gorm:"column:material"`
// 	Tagline      string `gorm:"column:tagline"`
// 	// Tags         db.JSONB `gorm:"column:tags;type:jsonb"`
// }

// type MarketplaceProductImage struct {
// 	*gorm.Model
// 	ID        uuid.UUID
// 	ImageUrl  string    `gorm:"column:image_url"`
// 	Tag       string    `gorm:"column:tag"`
// 	Primary   bool      `gorm:"column:primary"`
// 	ProductId uuid.UUID `gorm:"column:product_id"`
// 	CreatedAt time.Time `gorm:"column:created_at"`
// 	UpdatedAt time.Time `gorm:"column:updated_at"`
// }

// type InventoryProduct struct {
// 	*gorm.Model
// 	Id                   uuid.UUID
// 	SKU                  string    `gorm:"column:sku"`
// 	Name                 string    `gorm:"column:name"`
// 	SizeN                string    `gorm:"column:size_n"`
// 	Freesize             bool      `gorm:"column:freesize"`
// 	Attributes           db.JSONB  `gorm:"column:attributes;type:jsonb"`
// 	Quantity             int64     `gorm:"column:quantity"`
// 	CreatedAt            time.Time `gorm:"column:created_at"`
// 	UpdatedAt            time.Time `gorm:"column:updated_at"`
// 	MarketplaceProductID uuid.UUID `gorm:"column:marketplace_product_id"`
// }
