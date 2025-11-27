package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	IsVerified   bool           `gorm:"default:false" json:"is_verified"`
	Role         string         `gorm:"default:'user'" json:"role"`
	OIDCSubject  string         `gorm:"uniqueIndex" json:"-"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Profile       *Profile       `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Posts         []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments      []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	Reactions     []Reaction     `gorm:"foreignKey:UserID" json:"reactions,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
	Skills        []UserSkill    `gorm:"foreignKey:UserID" json:"skills,omitempty"`
}

type Profile struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DisplayName string     `json:"display_name"`
	Bio         string     `gorm:"type:text" json:"bio"`
	AvatarURL   string     `json:"avatar_url"`
	CoverURL    string     `json:"cover_url"`
	Department  string     `json:"department"`
	Position    string     `json:"position"`
	Location    string     `json:"location"`
	Phone       string     `json:"phone"`
	Birthday    *time.Time `json:"birthday"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Post struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	MediaURLs  string         `gorm:"type:jsonb" json:"media_urls"`
	Visibility string         `gorm:"default:'public'" json:"visibility"`
	GroupID    *uint          `gorm:"index" json:"group_id,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments  []Comment  `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Reactions []Reaction `gorm:"foreignKey:PostID" json:"reactions,omitempty"`
	Group     *Group     `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PostID    uint           `gorm:"not null;index" json:"post_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Post      *Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent    *Comment   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies   []Comment  `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Reactions []Reaction `gorm:"foreignKey:CommentID" json:"reactions,omitempty"`
}

type Reaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	PostID    *uint     `gorm:"index" json:"post_id,omitempty"`
	CommentID *uint     `gorm:"index" json:"comment_id,omitempty"`
	Type      string    `gorm:"not null" json:"type"`
	CreatedAt time.Time `json:"created_at"`

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post    *Post    `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Comment *Comment `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
}

type Connection struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	TargetID  uint      `gorm:"not null;index" json:"target_id"`
	Status    string    `gorm:"not null" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User   *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Target *User `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	ActorID   *uint     `gorm:"index" json:"actor_id,omitempty"`
	Type      string    `gorm:"not null" json:"type"`
	Title     string    `json:"title"`
	Message   string    `gorm:"type:text" json:"message"`
	Link      string    `json:"link"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	Data      string    `gorm:"type:jsonb" json:"data"`
	CreatedAt time.Time `json:"created_at"`

	User  *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Actor *User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

type Message struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SenderID   uint           `gorm:"not null;index" json:"sender_id"`
	ReceiverID uint           `gorm:"not null;index" json:"receiver_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	ReadAt     *time.Time     `json:"read_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Sender   *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver *User `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
}

type Group struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	AvatarURL   string         `json:"avatar_url"`
	CoverURL    string         `json:"cover_url"`
	Visibility  string         `gorm:"default:'public'" json:"visibility"`
	CreatorID   uint           `gorm:"not null;index" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Creator *User         `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Members []GroupMember `gorm:"foreignKey:GroupID" json:"members,omitempty"`
	Posts   []Post        `gorm:"foreignKey:GroupID" json:"posts,omitempty"`
}

type GroupMember struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	GroupID  uint      `gorm:"not null;index" json:"group_id"`
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	Role     string    `gorm:"default:'member'" json:"role"`
	JoinedAt time.Time `json:"joined_at"`

	Group *Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Skill struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Category    string         `json:"category"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	UserSkills []UserSkill `gorm:"foreignKey:SkillID" json:"user_skills,omitempty"`
}

type UserSkill struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	SkillID    uint      `gorm:"not null;index" json:"skill_id"`
	Level      string    `json:"level"`
	YearsOfExp int       `json:"years_of_experience"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User         *User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Skill        *Skill        `gorm:"foreignKey:SkillID" json:"skill,omitempty"`
	Endorsements []Endorsement `gorm:"foreignKey:UserSkillID" json:"endorsements,omitempty"`
}

type Endorsement struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserSkillID uint      `gorm:"not null;index" json:"user_skill_id"`
	EndorserID  uint      `gorm:"not null;index" json:"endorser_id"`
	Comment     string    `gorm:"type:text" json:"comment"`
	CreatedAt   time.Time `json:"created_at"`

	UserSkill *UserSkill `gorm:"foreignKey:UserSkillID" json:"user_skill,omitempty"`
	Endorser  *User      `gorm:"foreignKey:EndorserID" json:"endorser,omitempty"`
}

type File struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	Name       string         `gorm:"not null" json:"name"`
	Size       int64          `json:"size"`
	MimeType   string         `json:"mime_type"`
	StorageKey string         `gorm:"not null" json:"storage_key"`
	URL        string         `json:"url"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
