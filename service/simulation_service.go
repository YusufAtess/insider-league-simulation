package service

import (
	"insider-league/domain"
	"math"
	"math/rand"
	"sort"
)

type SimService struct {
	repo domain.Repository
}

func NewSimService(repo domain.Repository) *SimService {
	return &SimService{repo: repo}
}

func generatePoisson(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	for p > L {
		k++
		p *= rand.Float64()
	}
	return k - 1
}

func (s *SimService) simulateScore(homeStrength, awayStrength int) (int, int) {
	// Ev sahibi avantajı = x1.2
	homeLambda := (float64(homeStrength) / float64(awayStrength)) * 1.5
	awayLambda := (float64(awayStrength) / float64(homeStrength)) * 1.1

	return generatePoisson(homeLambda), generatePoisson(awayLambda)
}

func (s *SimService) calculateStandings(teams []domain.Team, matches []domain.Match) []domain.Standing {
	standingsMap := make(map[int]*domain.Standing)
	for _, t := range teams {
		standingsMap[t.ID] = &domain.Standing{TeamName: t.Name}
	}

	for _, m := range matches {
		if !m.IsPlayed {
			continue
		}
		home, away := standingsMap[m.HomeTeamID], standingsMap[m.AwayTeamID]

		home.Played++
		away.Played++
		home.GD += (m.HomeScore - m.AwayScore)
		away.GD += (m.AwayScore - m.HomeScore)

		if m.HomeScore > m.AwayScore {
			home.Won++
			home.Points += 3
			away.Lost++
		} else if m.HomeScore < m.AwayScore {
			away.Won++
			away.Points += 3
			home.Lost++
		} else {
			home.Drawn++
			away.Drawn++
			home.Points += 1
			away.Points += 1
		}
	}

	var table []domain.Standing
	for _, st := range standingsMap {
		table = append(table, *st)
	}

	sort.Slice(table, func(i, j int) bool {
		if table[i].Points == table[j].Points {
			return table[i].GD > table[j].GD
		}
		return table[i].Points > table[j].Points
	})
	return table
}

func (s *SimService) runMonteCarlo(teams []domain.Team, playedMatches, unplayedMatches []domain.Match) map[string]float64 {
	championCounts := make(map[string]int)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		simulatedMatches := make([]domain.Match, len(playedMatches))
		copy(simulatedMatches, playedMatches)

		for _, um := range unplayedMatches {
			var hStr, aStr int
			for _, t := range teams {
				if t.ID == um.HomeTeamID {
					hStr = t.Strength
				}
				if t.ID == um.AwayTeamID {
					aStr = t.Strength
				}
			}
			um.HomeScore, um.AwayScore = s.simulateScore(hStr, aStr)
			um.IsPlayed = true
			simulatedMatches = append(simulatedMatches, um)
		}

		finalTable := s.calculateStandings(teams, simulatedMatches)
		championCounts[finalTable[0].TeamName]++
	}

	probs := make(map[string]float64)
	for name, count := range championCounts {
		probs[name] = (float64(count) / float64(iterations)) * 100
	}
	return probs
}

func (s *SimService) GetStandings() ([]domain.Standing, error) {
	teams, _ := s.repo.GetTeams()
	matches, _ := s.repo.GetMatches()

	table := s.calculateStandings(teams, matches)

	var playedWeeks int
	var unplayedMatches []domain.Match
	for _, m := range matches {
		if m.IsPlayed {
			if m.Week > playedWeeks {
				playedWeeks = m.Week
			}
		} else {
			unplayedMatches = append(unplayedMatches, m)
		}
	}

	if playedWeeks >= 4 && len(unplayedMatches) > 0 {
		probs := s.runMonteCarlo(teams, matches, unplayedMatches)
		for i := range table {
			table[i].ChampProb = probs[table[i].TeamName]
		}
	}

	return table, nil
}

func (s *SimService) PlayNextWeek() error {
	matches, _ := s.repo.GetMatches()
	teams, _ := s.repo.GetTeams()

	nextWeek := -1
	for _, m := range matches {
		if !m.IsPlayed {
			nextWeek = m.Week
			break
		}
	}
	if nextWeek == -1 {
		return nil
	}

	for _, m := range matches {
		if m.Week == nextWeek {
			hStr, aStr := 0, 0
			for _, t := range teams {
				if t.ID == m.HomeTeamID {
					hStr = t.Strength
				}
				if t.ID == m.AwayTeamID {
					aStr = t.Strength
				}
			}
			hScore, aScore := s.simulateScore(hStr, aStr)
			s.repo.UpdateMatch(m.ID, hScore, aScore)
		}
	}
	return nil
}

func (s *SimService) PlayAll() error {
	for i := 0; i < 6; i++ {
		s.PlayNextWeek()
	}
	return nil
}

func (s *SimService) EditMatch(matchID, homeScore, awayScore int) error {
	return s.repo.UpdateMatch(matchID, homeScore, awayScore)
}

func (s *SimService) GetMatchesByWeek(week int) ([]domain.Match, error) {
	return s.repo.GetMatchesByWeek(week)
}
