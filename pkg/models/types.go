package models

type Image struct {
	ID          int     `json:"id"`
	SpeciesName string  `json:"species_name"`
	GpsLong     float64 `json:"gps_long"`
	GpsLat      float64 `json:"gps_lat"`
	ImagePath   string  `json:"image_path"`
	UserID      int     `json:"user_id"`
}

type ApiKey struct {
	ID        int     `json:"id"`
	Key       string  `json:"key"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
	ExpiresAt string  `json:"expires_at"`
	LastUsed  *string `json:"last_used"`
	Revoked   bool    `json:"revoked"`
}
