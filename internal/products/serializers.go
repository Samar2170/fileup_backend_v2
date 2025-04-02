package products

// import (
// 	"impx/internal/db"
// 	"time"

// 	"github.com/google/uuid"
// )

// type ListMarketplaceProduct struct {
// 	Id           uuid.UUID
// 	Name         string
// 	Description  string
// 	Category     string
// 	Mrp          float64
// 	Discount     float64
// 	SellingPrice float64
// 	Color        string
// 	Ratings      int64
// 	NumReviews   int64
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time

// 	ImageUrl string
// }

// func GetMarketplaceProductList() []ListMarketplaceProduct {
// 	var products []ListMarketplaceProduct
// 	db.EcommDB.
// 		Select("product_marketplaceproduct.*", "product_marketplaceproductimage.*").
// 		Table("product_marketplaceproduct").Joins("LEFT JOIN product_marketplaceproductimage ON product_marketplaceproductimage.product_id = product_marketplaceproduct.id").Find(&products)
// 	return products
// }

// type InventoryDetail struct {
// 	Id       uuid.UUID
// 	SizeN    string
// 	Quantity int64
// }
// type ImageDetail struct {
// 	ImageUrl string
// }

// type MarketplaceProductDetail struct {
// 	Id           uuid.UUID
// 	Name         string
// 	Description  string
// 	Category     string
// 	Mrp          float64
// 	Discount     float64
// 	SellingPrice float64
// 	Color        string
// 	Ratings      int64
// 	NumReviews   int64
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time

// 	Images           []ImageDetail
// 	InventoryDetails []InventoryDetail
// }

// func GetMarketplaceProductDetail(id string) MarketplaceProductDetail {
// 	var product MarketplaceProductDetail
// 	db.EcommDB.
// 		// Select("product_marketplaceproduct.*", "product_marketplaceproductimage.*").
// 		// Table("product_marketplaceproduct").Joins("LEFT JOIN product_marketplaceproductimage ON product_marketplaceproductimage.product_id = product_marketplaceproduct.id").
// 		Preload("product_marketplaceproductimage").Preload("product_inventoryproduct").
// 		Where("product_marketplaceproduct.id = ?", id).Find(&product)
// 	return product
// }
