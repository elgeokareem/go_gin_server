package auth

type LOGIN struct {
	EMAIL    string `json:"email" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

type REGISTER struct {
	EMAIL            string `json:"email" binding:"required"`
	PASSWORD         string `json:"password" binding:"required"`
	CONFIRM_PASSWORD string `json:"confirm_password" binding:"required"`
}
