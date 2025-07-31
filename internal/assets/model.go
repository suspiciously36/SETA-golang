package assets

type Folder struct {
	ID      string        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name    string        `gorm:"not null"`
	OwnerID string        `gorm:"not null"`
	Notes   []Note        `gorm:"foreignKey:FolderID"`
	Shares  []FolderShare `gorm:"foreignKey:FolderID"`
}

type Note struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title    string
	Body     string
	FolderID string      `gorm:"index"`
	OwnerID  string      `gorm:"not null"`
	Shares   []NoteShare `gorm:"foreignKey:NoteID"`
}

type FolderShare struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FolderID string `gorm:"index;not null"`
	UserID   string `gorm:"not null"`
	Access   string `gorm:"not null"` // "read" or "write"
}

type NoteShare struct {
	ID     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	NoteID string `gorm:"index;not null"`
	UserID string `gorm:"not null"`
	Access string `gorm:"not null"` // "read" or "write"
}
