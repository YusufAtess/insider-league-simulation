package domain

type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Strength int    `json:"-"`
}

type Match struct {
	ID         int  `json:"id"`
	HomeTeamID int  `json:"home_team_id"`
	AwayTeamID int  `json:"away_team_id"`
	Week       int  `json:"week"`
	HomeScore  int  `json:"home_score"`
	AwayScore  int  `json:"away_score"`
	IsPlayed   bool `json:"is_played"`
}

type Standing struct {
	TeamName  string  `json:"team_name"`
	Played    int     `json:"p"`
	Won       int     `json:"w"`
	Drawn     int     `json:"d"`
	Lost      int     `json:"l"`
	GD        int     `json:"gd"`
	Points    int     `json:"pts"`
	ChampProb float64 `json:"championship_probability"`
}

// Repository Interface
type Repository interface {
	GetTeams() ([]Team, error)
	GetMatches() ([]Match, error)
	UpdateMatch(matchID, homeScore, awayScore int) error
	GetMatchesByWeek(week int) ([]Match, error)
}

// Service Interface
type Simulation interface {
	GetStandings() ([]Standing, error)
	PlayNextWeek() error
	PlayAll() error
	EditMatch(matchID, homeScore, awayScore int) error
	GetMatchesByWeek(week int) ([]Match, error)
}
