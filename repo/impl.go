package repo

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/ISS-Dating/service-main/model"
)

var (
	userFields = []string{
		"id", "photo_url", "name", "surname", "username", "password", "email", "gender", "city",
		"country", "age", "description", "looking_for", "status", "education", "mood", "banned",
		"role",
	}
	questionFields = []string{
		"work", "food", "travel", "bio", "main",
	}
	statsFields = []string{
		"banned_before", "users_met", "messages_sent", "average_message_length", "links_in_messages",
	}
)

type Repository struct {
	db *sql.DB
}

func NewRepoImpl(db *sql.DB) *Repository {
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

	res, err := tx.Exec("INSERT INTO user("+strings.Join(userFields[1:], ", s")+") VALUES ("+
		generatePlaceholders(len(userFields[1:]))+")",
		getReadOnlyFields(user)[1:]...)

	if err != nil {
		tx.Rollback()
		return model.User{}, ErrorUsernameDuplication
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	_, err = tx.Exec("INSERT INTO question(user_id) VALUES ($1)",
		id)
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	_, err = tx.Exec("INSERT INTO stat(user_id) VALUES ($1)",
		id)
	if err != nil {
		tx.Rollback()
		return model.User{}, err
	}

	row := tx.QueryRow("SELECT "+strings.Join(userFields, ", ")+" FROM user WHERE username=$1 AND password=$2",
		user.Username, user.Password)
	if err := row.Scan(getModifyFields(&user)...); err != nil {
		return model.User{}, err
	}

	tx.Commit()
	return user, nil
}

func (r *Repository) ReadUserByLogin(username, password string) (*model.User, error) {
	data := &model.User{}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	res := r.db.QueryRow("SELECT "+strings.Join(userFields, ", ")+" FROM user WHERE username=$1 AND password=$2",
		username, password)
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

func readQuestions(tx *sql.Tx, id uint64) (model.Questionary, error) {
	var q model.Questionary
	res := tx.QueryRow("SELECT "+strings.Join(questionFields, ", ")+" FROM question WHERE user_id=$1",
		id)
	err := res.Scan(&q.Work, &q.Food, &q.Travel,
		&q.Biography, &q.Main)
	if err != nil {
		return q, err
	}

	return q, nil
}

func readStats(tx *sql.Tx, id uint64) (model.Stats, error) {
	var s model.Stats
	res := tx.QueryRow("SELECT "+strings.Join(statsFields, ", ")+" FROM stat WHERE user_id=$1",
		id)
	err := res.Scan(&s.BannedBefore, &s.UsersMet, &s.MessagesSent,
		&s.AverageMessageLen, &s.LinksInMessages)
	if err != nil {
		return s, err
	}

	return s, nil
}

func updateStats(tx *sql.Tx, id uint64, stats *model.Stats) error {
	toEdit, last := psqlJoin(statsFields, 1)
	_, err := tx.Exec("UPDATE stat SET "+toEdit+" WHERE user_id=$"+strconv.Itoa(last),
		stats.BannedBefore, stats.UsersMet, stats.MessagesSent, stats.AverageMessageLen, stats.LinksInMessages, id)
	if err != nil {
		return err
	}

	return nil
}

func updateQuestions(tx *sql.Tx, id uint64, questions *model.Questionary) error {
	toEdit, last := psqlJoin(questionFields, 1)
	_, err := tx.Exec("UPDATE stat SET "+toEdit+" WHERE user_id=$"+strconv.Itoa(last),
		questions.Work, questions.Food, questions.Travel, questions.Biography, questions.Main, id)
	if err != nil {
		return err
	}

	return nil
}

func updateUser(tx *sql.Tx, id uint64, user *model.User) error {
	toEdit, last := psqlJoin(userFields, 1)
	_, err := tx.Exec("UPDATE user SET "+toEdit+" WHERE user_id=$"+strconv.Itoa(last),
		user.ID, user.PhotoURL, user.Name, user.Surname, user.Username,
		user.Password, user.Email, user.Gender, user.City, user.Country, user.Age,
		user.Description, user.LookingFor, user.Status, user.Education, user.Education,
		user.Mood, user.Banned, user.Role, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ReadUserByUsername(username string) (*model.User, error) {
	data := &model.User{}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	res := r.db.QueryRow("SELECT "+strings.Join(userFields, ", ")+"FROM user WHERE username=$1",
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

func (r *Repository) UpdateUser(user *model.User) (*model.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	if err := updateUser(tx, user.ID, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := updateStats(tx, user.ID, &user.Stats); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := updateQuestions(tx, user.ID, &user.Questionary); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return user, nil
}

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

func psqlJoin(arr []string, start int) (string, int) {
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

func getReadOnlyFields(user model.User) []interface{} {
	var fields []interface{}
	fields = append(fields,
		user.ID,
		user.PhotoURL,
		user.Name,
		user.Surname,
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
	)

	return fields
}

func getModifyFields(user *model.User) []interface{} {
	var fields []interface{}
	fields = append(fields,
		&user.ID,
		&user.PhotoURL,
		&user.Name,
		&user.Surname,
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
	)

	return fields
}
