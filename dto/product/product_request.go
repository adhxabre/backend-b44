package productdto

type ProductRequest struct {
	Name       string `json:"name" form:"name"`
	Desc       string `json:"desc" form:"desc"`
	Price      int    `json:"price" form:"price"`
	Image      string `json:"image" form:"image"`
	Qty        int    `json:"qty" form:"qty"`
	UserID     int    `json:"user_id"`
	CategoryID int    `json:"category_id" form:"category_id"`
}
