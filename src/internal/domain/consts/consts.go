package consts

const (
	Empty uint = iota
	Cross
	Zero
)

const (
	Player = iota
	Computer
	Player1
	Player2
	Nobody
)

const (
	Active = iota
	Finished
	Draw
	AwaitingSecondPlayer
)

const (
	Rows = 3
	Columns
)

const (
	Negative = iota - 1
	Neutral
	Positive
)

const (
	VersusBot = iota
	VersusPlayer
)

const (
	BotMode    = "BotMode"
	PlayerMode = "PlayerMode"
)
