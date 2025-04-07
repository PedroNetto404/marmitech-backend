package dtos

type RestaurantDto struct {
	Id string `json:"id"`
	TradeName string	 `json:"trade_name"`
	LegalName string `json:"legal_name"`
	Cnpj string `json:"cnpj"`
}