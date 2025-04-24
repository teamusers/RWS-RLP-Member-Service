package requests

type UpdateBurnPin struct {
	Email   string `json:"email" binding:"required"`
	BurnPin int64  `json:"burn_pin" binding:"required"`
}
