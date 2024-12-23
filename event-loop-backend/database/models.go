package database

import "time"

// This struct is another table to hold team data
type Team struct {
	// TODO: update by using a UUID
	ID    string `gorm:"primaryKey"`
	Team  string
	Theme string
}

const (
	AccessRoleAdmin     = iota
	AccessRoleOrganizer = iota
	AccessRoleVolunteer = iota
)

// This struct represents the verified organizers and
// volunteers, holding the UUID (sub part of payload of
// google OAuth), VerifiedEmail gmails and user roles
// allowed roles:
//   - admain
//   - organizers
//   - volunteers
type DBAuthoriesedUsers struct {
	SUB           string `gorm:"primaryKey"`
	VerifiedEmail string `gorm:"VerifiedEmail"`
	UserRole      string `gorm:"UserRole"`

	// TOOD:
	// event name they belong to
}

// This struct holds all of the unified participant data stored in the DB
type DBParticipant struct {
	// ID is the generated JWT auth token
	// TODO: Change to a UUID instead of JWT
	ID          string      `gorm:"primaryKey"`
	Participant Participant `gorm:"embedded"`
	Checkpoints Checkpoints `gorm:"embedded"`
}

// This struct stores all fields parsed from the csv
type Participant struct {
	// TODO: Seperate to a separate
	// team table
	Team  string `json:"team"`
	Theme string `json:"theme"`

	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     int64  `json:"phone"`
	College   string `json:"college"`
	Branch    string `json:"branch"`
	PesHostel string `json:"pesHostel"`
}

// This struct stores all event checkpoints
type Checkpoints struct {
	// Other Participant Parameters
	Entry_time time.Time `json:"entry_time"`
	Checkin    bool      `json:"checkin"`
	Checkout   bool      `json:"checkout"`
	Exit_time  time.Time `json:"exit_time"`

	// TODO: Check with organizers regarding
	// the checkpoints of events
	Snacks    bool `json:"snacks"`
	Dinner    bool `json:"dinner"`
	Breakfast bool `json:"breakfast"`
}
