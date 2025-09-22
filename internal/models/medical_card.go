package models

// VisitHistoryItem представляет один элемент в истории посещений.
type VisitHistoryItem struct {
	Appointment
	PatientName string `json:"patientName"`
	DoctorName  string `json:"doctorName"`
	ServiceName string `json:"serviceName"`
}

// PaginatedVisitsResponse DTO для пагинированного списка визитов.
type PaginatedVisitsResponse struct {
	Total int64              `json:"total"`
	Items []VisitHistoryItem `json:"items"`
}

// RecentVisit DTO для информации о последнем визите в сводке.
type RecentVisit struct {
	Date       string `json:"date"`
	DoctorName string `json:"doctorName"`
}

// MedicalCardSummary DTO для сводки по медкарте.
type MedicalCardSummary struct {
	RecentVisit         *RecentVisit `json:"recentVisit"`
	PendingAnalyses     int          `json:"pendingAnalyses"`
	ActivePrescriptions int          `json:"activePrescriptions"`
}
