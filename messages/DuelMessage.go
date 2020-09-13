package messages

type DuelMessage struct {
	Winner      User `json:"winner"`
	Loser       User `json:"loser"`
	IsChallenge bool `json:"isChallenge"`
	IsGuildDuel bool `json:"isGuildDuel"`
}
