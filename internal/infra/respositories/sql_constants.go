package respositories

const (
	AddressFields = `
		a.id as address_id,
		a.alias as address_alias,
		a.street as address_street,
		a.number as address_number,
		a.complement as address_complement,
		a.neighborhood as address_neighborhood,
		a.city as address_city,
		a.state as address_state,
		a.country as address_country,
		a.zip_code as address_zip_code,
		a.lat as address_lat,
		a.lng as address_lng`
)
