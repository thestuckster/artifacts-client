package internal

type ArtifactsTask struct {
	Type           *string
	X, Y           *int
	Resource       *string
	AdditionalInfo *any
	Priority       int
}

type CharacterTask struct {
	Task func(characterName string)
}
