package types

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type State string 

const (
	Acre                State = "Acre"
	Alagoas             State = "Alagoas"
	Amapá               State = "Amapá"
	Amazonas            State = "Amazonas"
	Bahia               State = "Bahia"
	Ceará               State = "Ceará"
	DistritoFederal     State = "Distrito Federal"
	EspíritoSanto       State = "Espírito Santo"
	Goiás               State = "Goiás"
	Maranhão            State = "Maranhão"
	MatoGrosso          State = "Mato Grosso"
	MatoGrossoDoSul     State = "Mato Grosso do Sul"
	MinasGerais         State = "Minas Gerais"
	Pará                State = "Pará"
	Paraíba             State = "Paraíba"
	Paraná              State = "Paraná"
	Pernambuco          State = "Pernambuco"
	Piauí               State = "Piauí"
	RioDeJaneiro        State = "Rio de Janeiro"
	RioGrandeDoNorte    State = "Rio Grande do Norte"
	RioGrandeDoSul      State = "Rio Grande do Sul"
	Rondônia            State = "Rondônia"
	Roraima             State = "Roraima"
	SantaCatarina       State = "Santa Catarina"
	SãoPaulo            State = "São Paulo"
	Sergipe             State = "Sergipe"
	Tocantins           State = "Tocantins"
)

