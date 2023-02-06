package api

type TokenDetails struct {
	AccessToken   string
	AccessID      string
	RefreshToken  string
	RefreshID     string
	AccessExpiry  int64
	RefreshExpiry int64
}