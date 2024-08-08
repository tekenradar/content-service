package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LPPParticipant struct {
	ID                  primitive.ObjectID          `bson:"_id,omitempty" json:"id,omitempty"`
	PID                 string                      `bson:"pid,omitempty" json:"pid,omitempty"`
	InvitationSentAt    time.Time                   `bson:"invitationSentAt" json:"invitationSentAt"`
	ReminderSentAt      time.Time                   `bson:"reminderSentAt" json:"reminderSentAt"`
	Submissions         map[string]time.Time        `bson:"submissions" json:"submissions"`
	ContactInfos        *LPPParticipantContactInfos `bson:"contactInfos" json:"contactInfos,omitempty"`
	Cohort              string                      `bson:"cohort" json:"cohort"`
	StudyData           map[string]string           `bson:"studyData" json:"studyData"`
	TempParticipantInfo *TempParticipantInfo        `bson:"tempParticipantInfo" json:"tempParticipantInfo"`
}

type TempParticipantInfo struct {
	ID        string `bson:"id" json:"id"`
	EnteredAt int64  `bson:"enteredAt" json:"enteredAt"`
}

type LPPParticipantContactInfos struct {
	Email string `bson:"email" json:"email"`
	Name  string `bson:"name" json:"name"`
}
