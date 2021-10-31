package model

const (
	RoleUser          = "user"
	RoleModerator     = "moderator"
	RoleAdministrator = "administrator"
)

type User struct {
	ID          uint64      `json:"id"`
	PhotoURL    string      `json:"photo_url"`
	Name        string      `json:"name"`
	Surname     string      `json:"surname"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Email       string      `json:"email"`
	Gender      string      `json:"gender"`
	City        string      `json:"city"`
	Country     string      `json:"country"`
	Age         uint        `json:"age"`
	Description string      `json:"description"`
	LookingFor  string      `json:"looking_for"`
	Status      string      `json:"status"`
	Education   string      `json:"education"`
	Mood        string      `json:"mood"`
	Banned      bool        `json:"banned"`
	Role        string      `json:"role"`
	Stats       Stats       `json:"stats"`
	Questionary Questionary `json:"questions"`
}

type Stats struct {
	BannedBefore      bool `json:"banned_before"`
	UsersMet          uint `json:"users_met"`
	MessagesSent      uint `json:"messages_sent"`
	AverageMessageLen uint `json:"average_message_length"`
	LinksInMessages   uint `json:"links_in_messages"`
}

type Questionary struct {
	Work      string `json:"work_question"`
	Food      string `json:"food_question"`
	Travel    string `json:"travel_question"`
	Biography string `json:"biography_question"`
	Main      string `json:"main_question"`
}

type Acquaintance struct {
	UserAUsername uint64 `json:"user_a"`
	UserBUsername uint64 `json:"user_b"`
}
