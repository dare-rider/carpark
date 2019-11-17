package requests

type NearestCarparksRequest struct {
	Latitude  float64 `schema:"latitude" validate:"required"`
	Longitude float64 `schema:"longitude" validate:"required"`
	Page      int     `schema:"page" validate:"omitempty,min=1"`
	PerPage   int     `schema:"per_page" validate:"omitempty,min=1"`
}
