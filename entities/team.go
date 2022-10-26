package entities

// Team is the representation entity for teams
type Team struct {
	Name     string   `json:"name"`
	LogoFile *string  `json:"logo_file"`
	Slug     *string  `json:"slug"`
	Matches  *[]Match `json:"matches"`
}
