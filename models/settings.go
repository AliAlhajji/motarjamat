package models

type SiteSettings struct {
	IsRunning      bool   `db:"is_running" form:"isRunning"`
	Announcement   string `db:"announcement" form:"announcement"`
	Title          string `db:"title" form:"title"`
	ResultsPerPage int    `db:"results_per_page" form:"resultsPerPage"`
}
