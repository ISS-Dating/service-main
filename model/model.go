package model

import "time"

const (
	RoleUser          = "user"
	RoleModerator     = "moderator"
	RoleAdministrator = "administrator"
)

type User struct {
	ID               int64       `json:"id"`
	PhotoURL         string      `json:"photo_url"`
	Name             string      `json:"name"`
	Surname          string      `json:"surname"`
	Username         string      `json:"username"`
	Password         string      `json:"password"`
	Email            string      `json:"email"`
	Gender           string      `json:"gender"`
	City             string      `json:"city"`
	Country          string      `json:"country"`
	Age              uint        `json:"age"`
	Description      string      `json:"description"`
	LookingFor       string      `json:"looking_for"`
	Status           string      `json:"status"`
	Education        string      `json:"education"`
	Mood             string      `json:"mood"`
	Banned           bool        `json:"banned"`
	Role             string      `json:"role"`
	Stats            Stats       `json:"stats"`
	Questionary      Questionary `json:"questions"`
	RegistrationDate time.Time   `json:"registration_date"`
	Links            []string    `json:"links"`
}

type Stats struct {
	ID                int64 `json:"id"`
	UserID            int64 `json:"user_id"`
	BannedBefore      bool  `json:"banned_before"`
	UsersMet          uint  `json:"users_met"`
	MessagesSent      uint  `json:"messages_sent"`
	AverageMessageLen uint  `json:"average_message_length"`
	LinksInMessages   uint  `json:"links_in_messages"`
}

type Questionary struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Work      string `json:"work_question"`
	Food      string `json:"food_question"`
	Travel    string `json:"travel_question"`
	Biography string `json:"biography_question"`
	Main      string `json:"main_question"`
}

type Acquaintance struct {
	ID            int64 `json:"id"`
	UserAUsername int64 `json:"user_a"`
	UserBUsername int64 `json:"user_b"`
}
