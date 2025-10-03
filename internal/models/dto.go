package models

// Recommendation представляет DTO для текста рекомендации.
type Recommendation struct {
	Text string `json:"text"`
}

// AvailableDatesResponse представляет DTO для ответа со свободными датами.
type AvailableDatesResponse struct {
	SpecialistID   uint64   `json:"specialistId"`
	Month          string   `json:"month"`
	AvailableDates []string `json:"availableDates"`
}

// AvailableSlotsResponse представляет DTO для ответа со свободными слотами на ОДНУ дату.
type AvailableSlotsResponse struct {
	SpecialistID   uint64   `json:"specialistId"`
	Date           string   `json:"date"`
	AvailableSlots []string `json:"availableSlots"`
}

// SlotsForDay - вспомогательная структура для ответа по диапазону.
type SlotsForDay struct {
	Date           string   `json:"date"`
	AvailableSlots []string `json:"availableSlots"`
}

// AvailableRangeSlotsResponse представляет DTO для ответа со слотами в диапазоне дат.
type AvailableRangeSlotsResponse struct {
	SpecialistID uint64        `json:"specialistId"`
	ServiceID    uint64        `json:"serviceId"`
	SlotsByDay   []SlotsForDay `json:"slotsByDay"`
}

// PaginatedReviewsResponse представляет DTO для возврата пагинированного списка отзывов.
type PaginatedReviewsResponse struct {
	Page  int      `json:"page" example:"1"`
	Total int64    `json:"total" example:"150"`
	Items []Review `json:"items"`
}
