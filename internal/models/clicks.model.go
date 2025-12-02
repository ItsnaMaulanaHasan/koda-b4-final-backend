package models

type Click struct {
	ID          int
	ShortLinkID int
	ClickedAt   string
	IPAddress   string
	Referer     string
	UserAgent   string
	Country     string
	City        string
	DeviceType  string
	Browser     string
	OS          string
}
