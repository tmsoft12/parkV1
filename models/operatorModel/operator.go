package modeloperator

type Operator struct {
	ID       int64  `json:"id"`
	Park     string `json:"park"`
	LoginAt  string `json:"login_at"`
	LogoutAt string `json:"logout_at"`
	Operator string `json:"operator"`
	Money    int    `json:"money"`
}
