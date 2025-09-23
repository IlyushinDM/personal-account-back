package models

// Константы для статусов записи на прием (таблица appointmentstatuses)
const (
	StatusScheduled          uint32 = 1 // Запланировано
	StatusCompleted          uint32 = 2 // Завершено
	StatusCancelledByPatient uint32 = 3 // Отменено пациентом
	StatusCancelledByClinic  uint32 = 4 // Отменено клиникой
)

// Константы для статусов анализов (таблица analysisstatuses)
const (
	AnalysisStatusAssigned   uint32 = 1 // Назначено
	AnalysisStatusInProgress uint32 = 2 // В работе
	AnalysisStatusReady      uint32 = 3 // Готов
)
