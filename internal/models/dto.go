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

// DoctorReviewsResponse представляет DTO для ответа с отзывами врача и его рейтингом.
type DoctorReviewsResponse struct {
	Items        []Review `json:"items"`
	DoctorRating float64  `json:"doctorRating" example:"4.8"`
}

// CreateReviewRequest представляет DTO для создания отзыва.
type CreateReviewRequest struct {
	DoctorID uint64 `json:"doctor_id" binding:"required" example:"55"`
	Rating   uint16 `json:"rating" binding:"required,min=1,max=5" example:"5"`
	Comment  string `json:"comment" example:"Отличный врач!"`
}

// UpdateReviewRequest представляет DTO для обновления отзыва.
type UpdateReviewRequest struct {
	Rating  uint16 `json:"rating" binding:"required,min=1,max=5" example:"4"`
	Comment string `json:"comment" example:"Обновленный комментарий"`
}

// CreateReviewResponse представляет DTO для ответа при создании отзыва.
type CreateReviewResponse struct {
	ID      uint64 `json:"id" example:"1"`
	Message string `json:"message" example:"Отзыв создан"`
}

// UpdateReviewResponse представляет DTO для ответа при обновлении отзыва.
type UpdateReviewResponse struct {
	Message string `json:"message" example:"Отзыв обновлен"`
}

// DeleteReviewResponse представляет DTO для ответа при удалении отзыва.
type DeleteReviewResponse struct {
	Message string `json:"message" example:"Отзыв удален"`
}
