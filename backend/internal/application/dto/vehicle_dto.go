package dto

// Vehicle DTOs
type CreateVehicleRequest struct {
	ClientName  string `json:"client_name" validate:"required"`
	Phone       string `json:"phone"`
	Brand       string `json:"brand" validate:"required"`
	Model       string `json:"model" validate:"required"`
	Year        int    `json:"year"`
	Color       string `json:"color"`
	PlateNumber string `json:"plate_number"`
}

type UpdateVehicleRequest struct {
	ClientName  string `json:"client_name"`
	Phone       string `json:"phone"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Year        int    `json:"year"`
	Color       string `json:"color"`
	PlateNumber string `json:"plate_number"`
}
