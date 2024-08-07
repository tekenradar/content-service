package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LPPParticipant struct {
	ID               primitive.ObjectID          `bson:"_id,omitempty" json:"id,omitempty"`
	PID              string                      `bson:"pid,omitempty" json:"pid,omitempty"`
	InvitationSentAt time.Time                   `bson:"invitationSentAt" json:"invitationSentAt"`
	Submissions      map[string]time.Time        `bson:"submissions" json:"submissions"`
	ContactInfos     *LPPParticipantContactInfos `bson:"contactInfos" json:"contactInfos,omitempty"`
	StudyData        map[string]string           `bson:"studyData" json:"studyData"`
}

type LPPParticipantContactInfos struct {
	Email string `bson:"email" json:"email"`
	Name  string `bson:"name" json:"name"`
}
