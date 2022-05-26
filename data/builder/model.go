package builder

type Product struct {
	Id string `gorm:"primaryKey"`
	Name string
	Description string
	productImages []ProductImage
}

type ProductImage struct {
	id string `gorm:"primaryKey"`
	productId string
	filename string
}