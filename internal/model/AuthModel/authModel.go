package authmodel

import "time"

type User struct {
	ID                string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email             string `gorm:"size:255;unique;not null"`
	EmailVerified     bool   `gorm:"default:false"`
	DisplayName       string `gorm:"size:100"`
	ProfilePictureURL string
	AccountStatus     string `gorm:"type:varchar(20);default:'active'"`
	LastLoginAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time

	AuthMethods []AuthenticationMethod `gorm:"foreignKey:UserID"`
	Sessions    []Session              `gorm:"foreignKey:UserID"`
	UserToken   []UserToken            `gorm:"foreignKey:UserID"`
}

type AuthenticationMethod struct {
	ID                    string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID                string `gorm:"type:uuid;not null;index"` // ✅ ADD INDEX
	ProviderType          string `gorm:"type:varchar(20);index"`   // ✅ ADD INDEX
	ProviderUserID        string `gorm:"size:255"`
	PasswordHash          string
	PasswordLastChangedAt *time.Time
	OauthRefreshToken     string
	OauthTokenExpiry      *time.Time
	IsPrimary             bool `gorm:"default:false"`
	CreatedAt             time.Time
	UpdatedAt             time.Time

	User User `gorm:"constraint:OnDelete:CASCADE"`
}

type Session struct {
	SessionID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           string `gorm:"type:uuid;index"`
	RefreshTokenHash string `gorm:"size:255;not null;index"` // ✅ ADD INDEX
	DeviceInfo       string
	IPAddress        string
	ExpiresAt        time.Time
	IsRevoked        bool `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	User User `gorm:"constraint:OnDelete:CASCADE"`
}

type UserToken struct {
	TokenID   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string `gorm:"type:uuid;index;not null"`
	TokenHash string `gorm:"size:255;not null"`
	TokenType string `gorm:"type:varchar(50);not null"` // ← add not null
	ExpiresAt time.Time
	IsUsed    bool `gorm:"default:false"`
	CreatedAt time.Time

	User User `gorm:"constraint:OnDelete:CASCADE"`
}
