package models

type Product struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Image           string   `json:"image"`
	Price           float64  `json:"price"`
	DiscountedPrice float64  `json:"discountedPrice"`
	Rating          float64  `json:"rating"`
	Tags            []string `json:"tags"`
	IsNew           bool     `json:"isNew"`
	IsBestseller    bool     `json:"isBestseller"`
}
