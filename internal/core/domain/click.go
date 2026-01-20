package domain

import "time"

// Click represents a single click/visit on a shortened link.
type Click struct {
	ID         int64     `json:"id"`
	LinkID     int64     `json:"link_id"`
	Timestamp  time.Time `json:"timestamp"`
	Referrer   string    `json:"referrer,omitempty"`
	DeviceHash string    `json:"device_hash,omitempty"`
	UserAgent  string    `json:"-"`
	IPHash     string    `json:"-"`
}

// ClickStats contains aggregated click statistics for a link.
type ClickStats struct {
	TotalClicks   int64            `json:"total_clicks"`
	ClicksByDay   []DayStats       `json:"clicks_by_day,omitempty"`
	ClicksByMonth []MonthStats     `json:"clicks_by_month,omitempty"`
	TopReferrers  []ReferrerStats  `json:"top_referrers,omitempty"`
	DeviceStats   []DeviceStats    `json:"device_stats,omitempty"`
}

// DayStats contains click counts for a specific day.
type DayStats struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

// MonthStats contains click counts for a specific month.
type MonthStats struct {
	Month  string `json:"month"`
	Clicks int64  `json:"clicks"`
}

// ReferrerStats contains click counts from a specific referrer.
type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Clicks   int64  `json:"clicks"`
}

// DeviceStats contains click counts by device type.
type DeviceStats struct {
	DeviceType string `json:"device_type"`
	Clicks     int64  `json:"clicks"`
}

// StatsPeriod defines the time range for statistics queries.
type StatsPeriod string

const (
	StatsPeriodDay   StatsPeriod = "day"
	StatsPeriodWeek  StatsPeriod = "week"
	StatsPeriodMonth StatsPeriod = "month"
	StatsPeriodYear  StatsPeriod = "year"
	StatsPeriodAll   StatsPeriod = "all"
)

// StatsFilter contains filter options for retrieving statistics.
type StatsFilter struct {
	Period    StatsPeriod `json:"period,omitempty"`
	StartDate *time.Time  `json:"start_date,omitempty"`
	EndDate   *time.Time  `json:"end_date,omitempty"`
}
