package dishtype

type DishType string

const (
	//carnes
	MEAT DishType = "meat"
	// acompanhamento
	ACCOMPANIMENT DishType = "accompaniment"
	//guarnição
	SIDE_DISH DishType = "side_dish"
	//sobremesa
	DESSERT DishType = "dessert"
	//bebida
	DRINK DishType = "drink"
	//salada	
	SALAD DishType = "salad"
	//outros
	OTHER DishType = "other"
)