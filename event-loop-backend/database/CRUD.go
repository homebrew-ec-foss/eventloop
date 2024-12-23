package database

import (
	"fmt"
	"log"
	"slices"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbGlobal  *gorm.DB
	errGlobal error
)

// Rather proposed custom error types
var (
	ErrDbOpenFailure   = fmt.Errorf("failed to run `OpenDB()`")
	ErrDbMissingRecord = fmt.Errorf("failed to fetch record")
	ErrIncorrectField  = fmt.Errorf("checkpoint missing in db")
)

// Event specific errors
// DB_Participants related errors
var (
	ErrParticipantAbsent = fmt.Errorf("participant never checkedin")
	ErrParticipantLeft   = fmt.Errorf("participant has left the event")
)

// Event authorised user specific
var (
	ErrUnauthorised = fmt.Errorf("incoming request isn't from an authorised login")
	ErrNoAccess     = fmt.Errorf("incoming request was authoriesed but has no access to the endpoint")
)

// Open and return db access struct
func openDB() (*gorm.DB, error) {
	if dbGlobal == nil {
		dbGlobal, errGlobal = gorm.Open(sqlite.Open("event.db"), &gorm.Config{})
		if errGlobal != nil {
			log.Println(errGlobal)
			return nil, errGlobal
		}

		dbGlobal.AutoMigrate(&DBParticipant{}, &DBAuthoriesedUsers{})
	}

	return dbGlobal, nil
}

// Create records for all participants parsed from the csv
func CreateParticipants(dbParticipants []DBParticipant) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	log.Println("Attempting to write to db")
	// Writing to DB
	db.Create(dbParticipants)

	return nil
}

func CreateAuthorisedUsersDB() error {
	db, err := openDB()
	if err != nil {
		return ErrDbOpenFailure
	}

	var dbAuth []DBAuthoriesedUsers

	db.Create(&dbAuth)

	return nil
}

// Function to fetch all checkpoints and
// forward to backend
//
// for dynamic checkpoint loading for site
func FetchCheckpoints() []string {
	var checkpoints []string

	return checkpoints
}

func VerifyLogin(userDetails DBAuthoriesedUsers) (*DBAuthoriesedUsers, error) {
	db, err := openDB()
	if err != nil {
		return nil, ErrDbOpenFailure
	}

	log.Println("request: ", userDetails)

	var dbAuthUser DBAuthoriesedUsers
	_ = db.First(&dbAuthUser, "sub = ?", userDetails.SUB)

	if dbAuthUser.VerifiedEmail == "" {

		log.Println("Missing from records")

		// need admin side approval for
		// organisers and volunteers
		admins := []string{
			"adityahegde.clg@gmail.com",
			"adheshathrey2004@gmail.com",
		}

		organisers := []string{
			"anirudh.sudhir1@gmail.com",
		}

		volunteers := []string{
			"adimhegde@gmail.com",
			"naysha.k0708@gmail.com",
			"devesh6742@gmail.com",
			"omshivshankar21@gmail.com",
			"vickspatil1404@gmail.com",
			"kunalkishoremaverick@gmail.com",
			"kavyaprakashscei@gmail.com",
			"nehanshetty2003@gmail.com",
			"b.himank101@gmail.com",
			"moulikmachaiah724@gmail.com",
			"manum262sagara@gmail.com",
			"disha14072003@gmail.com",
			"ananya975.p@gmail.com",
			"prathamshetty0826@gmail.com",
			"ruthu.hm03@gmail.com",
			"shashanknadigm03@gmail.com",
			"keerthanaumesh161@gmail.com",
			"shreyalizbethrobin@gmail.com",
			"eshwarra5@gmail.com",
			"jiteshnayak2004@gmail.com",
			"sarkarsoham73@gmail.com",
			"santoshrajpurohit89@gmail.com",
			"prawnee99@gmail.com",
			"shubhammookim@gmail.com",
			"kushagraagarwal2003@gmail.com",
			"roshinlinson67281@gmail.com",
		}

		if slices.Contains(admins, userDetails.VerifiedEmail) {
			userDetails.UserRole = "admin"
			log.Println("Hello admin")
		} else if slices.Contains(organisers, userDetails.VerifiedEmail) {
			userDetails.UserRole = "organiser"
			log.Println("Hello organiser")
		} else if slices.Contains(volunteers, userDetails.VerifiedEmail) {
			userDetails.UserRole = "volunteer"
		} else {
			return nil, ErrDbMissingRecord
		}

		db.Create(userDetails)
		return &userDetails, nil
	}

	return &dbAuthUser, nil
}

func SubAuthentication(sub string, userRole string) (*DBAuthoriesedUsers, error) {
	db, err := openDB()
	if err != nil {
		return nil, ErrDbOpenFailure
	}

	var dbAuthUser DBAuthoriesedUsers
	_ = db.First(&dbAuthUser, "sub = ?", sub)

	if dbAuthUser.VerifiedEmail == "" {
		// bro doesnt exist
		return nil, ErrDbMissingRecord
	}

	if dbAuthUser.UserRole != userRole {
		return nil, ErrNoAccess
	} else {
		return &dbAuthUser, nil
	}
}

func JWTFetchParticipant(jwt string) (*DBParticipant, error) {
	db, err := openDB()
	if err != nil {
		return nil, ErrDbOpenFailure
	}

	var dbParticipant DBParticipant
	_ = db.First(&dbParticipant, "id = ?", jwt)

	return &dbParticipant, nil
}

func FetchParticipant(name string, phone string) (*DBParticipant, error) {
	db, err := openDB()
	if err != nil {
		return nil, ErrDbOpenFailure
	}

	var dbParticipant DBParticipant
	_ = db.First(&dbParticipant, "name = ? and phone = ?", name, phone)

	log.Println(dbParticipant)

	if dbParticipant.Participant.Name == "" {
		return nil, ErrDbMissingRecord
	}

	return &dbParticipant, nil
}

// Update DB with the participant entry checkpoint
//
// Return signature
// - Pointer to participant
// - checkin : true if not checked in
// - error
func ParticipantEntry(p_uuid string) (*DBParticipant, bool, error) {
	db, err := openDB()
	if err != nil {
		return nil, false, ErrDbOpenFailure
	}

	var dbParticipant DBParticipant
	_ = db.First(&dbParticipant, "id = ?", p_uuid)
	log.Println(dbParticipant.Participant)

	if dbParticipant.Participant.Name == "" {
		// FIX:
		// Returns a pointer to an empty struct
		return nil, false, ErrDbMissingRecord
	}

	if dbParticipant.Checkpoints.Checkin && dbParticipant.Checkpoints.Checkout {
		return nil, false, ErrParticipantLeft
	}

	if !dbParticipant.Checkpoints.Checkin {
		dbParticipant.Checkpoints.Entry_time = time.Now()
		dbParticipant.Checkpoints.Checkin = true
		db.Save(&dbParticipant)
		return &dbParticipant, true, nil
	}

	return &dbParticipant, false, nil
}

// Update DB with the participant exit checkpoint
//
// Return signature
//   - Pointer to participant
//   - checkin : true if sucessfulyl checked in and
//     false if alreayd checked in
//   - error
func ParticipantExit(p_uuid string) (*DBParticipant, bool, error) {
	db, err := openDB()
	if err != nil {
		return nil, false, ErrDbOpenFailure
	}

	var dbParticipants DBParticipant

	_ = db.First(&dbParticipants, "id = ?", p_uuid)

	// Extra check if there has been some tampered entry
	// in the db
	if dbParticipants.Participant.Name == "" {
		return nil, false, ErrDbMissingRecord
	}

	// Participant MUST have checked in
	// for the checkout proceedure to be valid
	if !dbParticipants.Checkpoints.Checkin {
		return nil, false, ErrParticipantAbsent
	}

	if !dbParticipants.Checkpoints.Checkout {
		dbParticipants.Checkpoints.Checkout = true
		dbParticipants.Checkpoints.Exit_time = time.Now()
		db.Save(&dbParticipants)
		return &dbParticipants, true, nil
	}

	return &dbParticipants, false, nil
}

func ParticipantCheckpoint(p_uuid string, checkpointName string) (*DBParticipant, bool, error) {
	db, err := openDB()
	if err != nil {
		return nil, false, ErrDbOpenFailure
	}

	var dbParticipants DBParticipant

	_ = db.First(&dbParticipants, "id = ?", p_uuid)

	// Check if participant is in the db
	if dbParticipants.Participant.Name == "" {
		return nil, false, ErrDbMissingRecord
	}

	if !dbParticipants.Checkpoints.Checkin {
		return nil, false, ErrParticipantAbsent
	}

	if dbParticipants.Checkpoints.Checkin && dbParticipants.Checkpoints.Checkout {
		return nil, false, ErrParticipantLeft
	}

	// FIX: Refractor
	switch checkpointName {
	case "Breakfast":
		{
			if dbParticipants.Checkpoints.Breakfast {
				break
			}
			dbParticipants.Checkpoints.Breakfast = true
			db.Save(&dbParticipants)
			return &dbParticipants, true, nil
		}
	case "Dinner":
		{
			if dbParticipants.Checkpoints.Dinner {
				break
			}
			dbParticipants.Checkpoints.Dinner = true
			db.Save(&dbParticipants)
			return &dbParticipants, true, nil
		}
	case "Snacks":
		{
			if dbParticipants.Checkpoints.Snacks {
				break
			}
			dbParticipants.Checkpoints.Snacks = true
			db.Save(&dbParticipants)
			return &dbParticipants, true, nil
		}
	default:
		{
			return nil, false, ErrIncorrectField
		}
	}

	// Participant has already opted for the option
	return &dbParticipants, false, nil
}
