package score

type TeamScore struct {
	Name  string
	Score int
}

type ScoreManager struct {
	Teams       [2]*TeamScore
	ActiveMatch bool
	MatchTime   int // In seconds
}

func NewScoreManager() *ScoreManager {
	return &ScoreManager{
		Teams: [2]*TeamScore{
			{Name: "Home", Score: 0},
			{Name: "Away", Score: 0},
		},
		ActiveMatch: false,
		MatchTime:   0,
	}
}

func (sm *ScoreManager) AddGoal(teamIndex int) {
	if teamIndex >= 0 && teamIndex < len(sm.Teams) {
		sm.Teams[teamIndex].Score++
	}
}

func (sm *ScoreManager) GetScore(teamIndex int) int {
	if teamIndex >= 0 && teamIndex < len(sm.Teams) {
		return sm.Teams[teamIndex].Score
	}
	return 0
}

func (sm *ScoreManager) GetTeamName(teamIndex int) string {
	if teamIndex >= 0 && teamIndex < len(sm.Teams) {
		return sm.Teams[teamIndex].Name
	}
	return ""
}

func (sm *ScoreManager) SetTeamName(teamIndex int, name string) {
	if teamIndex >= 0 && teamIndex < len(sm.Teams) {
		sm.Teams[teamIndex].Name = name
	}
}

func (sm *ScoreManager) StartMatch(matchDuration int) {
	sm.Teams[0].Score = 0
	sm.Teams[1].Score = 0
	sm.ActiveMatch = true
	sm.MatchTime = matchDuration
}

func (sm *ScoreManager) UpdateMatchTime(deltaSeconds int) bool {
	if !sm.ActiveMatch {
		return false
	}

	sm.MatchTime -= deltaSeconds
	if sm.MatchTime <= 0 {
		sm.MatchTime = 0
		sm.ActiveMatch = false
		return false
	}
	return true
}

func (sm *ScoreManager) GetMatchTime() int {
	return sm.MatchTime
}

func (sm *ScoreManager) IsMatchActive() bool {
	return sm.ActiveMatch
}
