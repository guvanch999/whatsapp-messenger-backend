package models

import (
	"github.com/google/uuid"
	"github.com/medium-messenger/messenger-backend/utils/enums"
	supa "github.com/nedpals/supabase-go"
	"time"
)

type User struct {
	ID                 uuid.UUID `json:"id"`
	Aud                string    `json:"aud"`
	Role               string    `json:"role"`
	Email              string    `json:"email"`
	InvitedAt          time.Time `json:"invited_at"`
	ConfirmedAt        time.Time `json:"confirmed_at"`
	ConfirmationSentAt time.Time `json:"confirmation_sent_at"`
	UserMetadata       map[string]interface {
	} `json:"user_metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *User) FromSupabaseUser(user *supa.User) {
	d.ID = uuid.MustParse(user.ID)
	d.Aud = user.Aud
	d.Email = user.Email
	d.InvitedAt = user.InvitedAt
	d.ConfirmedAt = user.ConfirmedAt
	d.ConfirmationSentAt = user.ConfirmationSentAt
	d.UserMetadata = user.UserMetadata
	d.CreatedAt = user.CreatedAt
	d.UpdatedAt = user.UpdatedAt
}

type UserDetail struct {
	User
	Role enums.Role `json:"role"`
}

func (d *UserDetail) FromSupabaseUser(user *supa.User, info *UserInfo) {
	d.ID = uuid.MustParse(user.ID)
	d.Aud = user.Aud
	d.Email = user.Email
	d.InvitedAt = user.InvitedAt
	d.ConfirmedAt = user.ConfirmedAt
	d.ConfirmationSentAt = user.ConfirmationSentAt
	d.UserMetadata = user.UserMetadata
	d.CreatedAt = user.CreatedAt
	d.UpdatedAt = user.UpdatedAt
	d.Role = info.Role
}

func (d *UserDetail) FromUser(user *User, info *UserInfo) {
	d.ID = user.ID
	d.Aud = user.Aud
	d.Email = user.Email
	d.InvitedAt = user.InvitedAt
	d.ConfirmedAt = user.ConfirmedAt
	d.ConfirmationSentAt = user.ConfirmationSentAt
	d.UserMetadata = user.UserMetadata
	d.CreatedAt = user.CreatedAt
	d.UpdatedAt = user.UpdatedAt
	d.Role = info.Role
}

type UserInfo struct {
	Id        uuid.UUID  `json:"id,omitempty" gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
	Role      enums.Role `json:"role"`
	UserGuid  uuid.UUID  `json:"user_guid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
