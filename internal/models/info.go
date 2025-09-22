package models

// ClinicInfo представляет DTO с контактной информацией о клинике.
type ClinicInfo struct {
	Name         string      `json:"name" example:"Клиника 'Здоровье'"`
	Contacts     []Contact   `json:"contacts"`
	Addresses    []Address   `json:"addresses"`
	WorkingHours []WorkHours `json:"workingHours"`
}

// Contact представляет контактные данные (телефон, email).
type Contact struct {
	Type  string `json:"type" example:"phone"`
	Value string `json:"value" example:"+7 (495) 222-33-44"`
}

// Address представляет адрес клиники.
type Address struct {
	ID      int    `json:"id" example:"1"`
	Address string `json:"address" example:"ул. Ленина, д. 100"`
	IsMain  bool   `json:"isMain" example:"true"`
}

// WorkHours представляет часы работы.
type WorkHours struct {
	Days  string `json:"days" example:"пн-пт"`
	Hours string `json:"hours" example:"08:00 - 20:00"`
}

// LegalDocument представляет DTO для юридического документа.
type LegalDocument struct {
	ID         uint64 `gorm:"primarykey" json:"-"`
	Type       string `json:"type" example:"privacy_policy"`
	Title      string `json:"title" example:"Политика конфиденциальности"`
	URL        string `json:"url" example:"/legal/privacy.pdf"`
	Version    string `json:"version" example:"2.0"`
	UpdateDate string `json:"updateDate" example:"2023-09-15"`
}
