package dto

// Service DTOs
type ServiceItemRequest struct {
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,gt=0"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type CreateServiceRequest struct {
	VehicleID   string               `json:"vehicle_id" validate:"required"`
	Description string               `json:"description"`
	Items       []ServiceItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateServiceRequest struct {
	Description string               `json:"description"`
	Items       []ServiceItemRequest `json:"items"`
	Status      string               `json:"status"`
}

type QuickTicketServiceItem struct {
	Descripcion string  `json:"descripcion" validate:"required"`
	Cantidad    int     `json:"cantidad" validate:"required,min=1"`
	Costo       float64 `json:"costo" validate:"required,min=0"`
}

type QuickTicketRequest struct {
	ClienteNombre   string                   `json:"clienteNombre" validate:"required"`
	ClienteTelefono string                   `json:"clienteTelefono"`
	VehiculoDesc    string                   `json:"vehiculoDesc" validate:"required"`
	VehiculoColor   string                   `json:"vehiculoColor"`
	Servicios       []QuickTicketServiceItem `json:"servicios" validate:"required,min=1,dive"`
}
