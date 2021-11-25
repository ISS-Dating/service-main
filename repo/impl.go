package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ISS-Dating/service-main/model"
	"github.com/lib/pq"
)

var (
	userFields = []string{
		"id", "photo_url", "name", "surname", "username", "password", "email", "gender", "city",
		"country", "age", "description", "looking_for", "status", "education", "mood", "banned",
		"role", "registration_date", "links",
	}
	questionFields = []string{
		"id", "work", "food", "travel", "bio", "main", "user_id",
	}
	statsFields = []string{
		"id", "banned_before", "users_met", "messages_sent", "average_message_length", "links_in_messages", "user_id",
	}
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(user model.User) (model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return model.User{}, err
	}

	if user.Role == "" {
		user.Role = model.RoleUser
	}
	if user.Links == nil {
		user.Links = make([]string, 0)
	}
	user.RegistrationDate = time.Now()

	err = tx.QueryRow("INSERT INTO \"user\" ("+strings.Join(userFields[1:], ", ")+") VALUES ("+
		generatePlaceholders(len(userFields[1:]))+") RETURNING id",
		getReadUserFields(user)[1:]...).Scan(&user.ID)

	if err != nil {
		tx.Rollback()
		return model.User{}, fmt.Errorf("%s, %w", err.Error(), ErrorUsernameDuplication)
	}

	user.Questionary.UserID = user.ID
	err = tx.QueryRow("INSERT INTO questionary("+strings.Join(questionFields[1:], ", ")+") VALUES ("+
		generatePlaceholders(len(questionFields[1:]))+") RETURNING id",
		getReadQuestionFields(user.Questionary)[1:]...).Scan(&user.Questionary.ID)
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	user.Stats.UserID = user.ID
	err = tx.QueryRow("INSERT INTO stats("+strings.Join(statsFields[1:], ", ")+") VALUES ("+
		generatePlaceholders(len(statsFields[1:]))+") RETURNING id",
		getReadStatsFields(user.Stats)[1:]...).Scan(&user.Stats.ID)
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	tx.Commit()
	return user, nil
}

func (r *Repository) ReadUserByLogin(username, password string) (model.User, error) {
	var user model.User

	tx, err := r.db.Begin()
	if err != nil {
		return model.User{}, err
	}

	row := r.db.QueryRow("SELECT "+strings.Join(userFields, ", ")+" FROM \"user\" WHERE username=$1 AND password=$2",
		username, password)
	if err := row.Scan(getModifyUserFields(&user)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrorUserNotExist
		}
		return model.User{}, err
	}

	user.Stats, err = readStats(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	user.Questionary, err = readQuestions(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	tx.Commit()
	return user, nil
}

func readQuestions(tx *sql.Tx, id int64) (model.Questionary, error) {
	var q model.Questionary
	res := tx.QueryRow("SELECT "+strings.Join(questionFields, ", ")+" FROM questionary WHERE user_id=$1",
		id)
	err := res.Scan(getModifyQuestionFields(&q)...)
	if err != nil {
		return q, err
	}

	return q, nil
}

func readStats(tx *sql.Tx, id int64) (model.Stats, error) {
	var s model.Stats
	res := tx.QueryRow("SELECT "+strings.Join(statsFields, ", ")+" FROM stats WHERE user_id=$1",
		id)
	err := res.Scan(getModifyStatsFields(&s)...)
	if err != nil {
		return s, err
	}

	return s, nil
}

// func updateStats(tx *sql.Tx, id int64, stats *model.Stats) error {
// 	toEdit, last := psqlJoin(statsFields, 1)
// 	_, err := tx.Exec("UPDATE stat SET "+toEdit+" WHERE user_id=$"+strconv.Itoa(last),
// 		stats.BannedBefore, stats.UsersMet, stats.MessagesSent, stats.AverageMessageLen, stats.LinksInMessages, id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func updateQuestions(tx *sql.Tx, id uint64, questions *model.Questionary) error {
// 	toEdit, last := psqlJoin(questionFields, 1)
// 	_, err := tx.Exec("UPDATE stat SET "+toEdit+" WHERE user_id=$"+strconv.Itoa(last),
// 		questions.Work, questions.Food, questions.Travel, questions.Biography, questions.Main, id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func updateUser(tx *sql.Tx, id uint64, user *model.User) error {
// 	toEdit, last := psqlJoin(userFields, 1)
// 	_, err := tx.Exec("UPDATE user SET "+toEdit+" WHERE user_id=$"+strconv.Itoa(last),
// 		user.ID, user.PhotoURL, user.Name, user.Surname, user.Username,
// 		user.Password, user.Email, user.Gender, user.City, user.Country, user.Age,
// 		user.Description, user.LookingFor, user.Status, user.Education, user.Education,
// 		user.Mood, user.Banned, user.Role, id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *Repository) ReadUserByUsername(username string) (*model.User, error) {
	data := &model.User{}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	res := r.db.QueryRow("SELECT "+strings.Join(userFields, ", ")+"FROM \"user\" WHERE username=$1",
		username)
	err = res.Scan(&data.ID, &data.PhotoURL, &data.Name, &data.Surname, &data.Username,
		&data.Password, &data.Email, &data.Gender, &data.City, &data.Country, &data.Age,
		&data.Description, &data.LookingFor, &data.Status, &data.Education, &data.Education,
		&data.Mood, &data.Banned, &data.Role)
	if err != nil {
		tx.Rollback()
		return nil, ErrorUserNotExist
	}

	data.Stats, err = readStats(tx, data.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	data.Questionary, err = readQuestions(tx, data.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return data, nil
}

// func (r *Repository) UpdateUser(user *model.User) (*model.User, error) {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := updateUser(tx, user.ID, user); err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	if err := updateStats(tx, user.ID, &user.Stats); err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	if err := updateQuestions(tx, user.ID, &user.Questionary); err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	tx.Commit()
// 	return user, nil
// }

func (r *Repository) CreateAcquaintance(userA, userB string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO acquaintance(user_a, user_b) VALUES($1, $2)", userA, userB)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repository) GetAcquaintanceByUsername(username string) ([]model.Acquaintance, error) {
	var acc []model.Acquaintance

	tx, err := r.db.Begin()
	if err != nil {
		return acc, nil
	}

	rows, err := tx.Query("SELECT user_a, user_b FROM acquaintance WHERE user_a=$1", username)
	if err != nil {
		tx.Rollback()
		return acc, err
	}

	defer rows.Close()
	for rows.Next() {
		var a model.Acquaintance
		if err := rows.Scan(&a.UserAUsername, &a.UserBUsername); err != nil {
			tx.Rollback()
			return acc, err
		}
		acc = append(acc, a)
	}

	tx.Commit()
	return acc, nil
}

func generateEqualsPlaceholder(arr []string, start int) (string, int) {
	b := strings.Builder{}
	for _, s := range arr {
		b.WriteString(s)
		b.WriteString("=$")
		b.WriteString(strconv.Itoa(start))
		b.WriteString(", ")
		start++
	}
	res := b.String()
	return strings.TrimSuffix(res, ", "), start
}

func generatePlaceholders(n int) string {
	var placeholders []string
	for i := 0; i < n; i++ {
		placeholders = append(placeholders, "$"+strconv.Itoa(i+1))
	}

	return strings.Join(placeholders, ", ")
}

func getReadUserFields(user model.User) []interface{} {
	var fields []interface{}
	fields = append(fields,
		user.ID,
		user.PhotoURL,
		user.Name,
		user.Surname,
		user.Username,
		user.Password,
		user.Email,
		user.Gender,
		user.City,
		user.Country,
		user.Age,
		user.Description,
		user.LookingFor,
		user.Status,
		user.Education,
		user.Mood,
		user.Banned,
		user.Role,
		user.RegistrationDate,
		pq.Array(user.Links),
	)

	return fields
}

func getModifyUserFields(user *model.User) []interface{} {
	var fields []interface{}
	fields = append(fields,
		&user.ID,
		&user.PhotoURL,
		&user.Name,
		&user.Surname,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Gender,
		&user.City,
		&user.Country,
		&user.Age,
		&user.Description,
		&user.LookingFor,
		&user.Status,
		&user.Education,
		&user.Mood,
		&user.Banned,
		&user.Role,
		&user.RegistrationDate,
		pq.Array(&user.Links),
	)

	return fields
}

func getReadStatsFields(stats model.Stats) []interface{} {
	var fields []interface{}
	fields = append(fields,
		stats.ID,
		stats.BannedBefore,
		stats.UsersMet,
		stats.MessagesSent,
		stats.AverageMessageLen,
		stats.LinksInMessages,
		stats.UserID,
	)

	return fields
}

func getModifyStatsFields(stats *model.Stats) []interface{} {
	var fields []interface{}
	fields = append(fields,
		&stats.ID,
		&stats.BannedBefore,
		&stats.UsersMet,
		&stats.MessagesSent,
		&stats.AverageMessageLen,
		&stats.LinksInMessages,
		&stats.UserID,
	)

	return fields
}

func getReadQuestionFields(question model.Questionary) []interface{} {
	var fields []interface{}
	fields = append(fields,
		question.ID,
		question.Work,
		question.Food,
		question.Travel,
		question.Biography,
		question.Main,
		question.UserID,
	)

	return fields
}

func getModifyQuestionFields(question *model.Questionary) []interface{} {
	var fields []interface{}
	fields = append(fields,
		&question.ID,
		&question.Work,
		&question.Food,
		&question.Travel,
		&question.Biography,
		&question.Main,
		&question.UserID,
	)

	return fields
}
