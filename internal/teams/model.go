package teams

type Team struct {
	ID       string        `gorm:"primaryKey" json:"teamId"`
	Name     string        `json:"teamName"`
	Managers []TeamManager `json:"managers" gorm:"foreignKey:TeamID"`
	Members  []TeamMember  `json:"members" gorm:"foreignKey:TeamID"`
}

type TeamManager struct {
	ID     string `json:"managerId"`
	Name   string `json:"managerName"`
	TeamID string
}

type TeamMember struct {
	ID     string `json:"memberId"`
	Name   string `json:"memberName"`
	TeamID string
}
