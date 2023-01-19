package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type CreatedUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetSt := "SELECT id,name,email,password,created_at,updated_at FROM public.users"
		rows, err := db.Query(GetSt)
		if err != nil {
			w.WriteHeader(422)
			w.Write([]byte("error getting users: " + err.Error()))
			return
		}
		defer rows.Close()
		var users []User
		for rows.Next() {
			var user User
			err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Created_at, &user.Updated_at)
			if err != nil {
				w.Write([]byte("Error getting users " + err.Error()))
				return
			}
			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			w.Write([]byte("Error on rows"))
			return
		}
		j, err := json.Marshal(users)
		if err != nil {
			w.Write([]byte("Error encoding users to JSON"))
			return
		}
		w.Write([]byte(j))
	}
}

func GetUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.Write([]byte("ID is not a number"))
			return
		}
		row := db.QueryRow("SELECT id, name, email, password, created_at, updated_at FROM public.users WHERE id=$1", id)
		var user User
		switch err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Created_at, &user.Updated_at); err {
		case sql.ErrNoRows:
			w.Write([]byte("No rows were returned"))
		case nil:
			j, err := json.Marshal(user)
			if err != nil {
				w.Write([]byte("error turning to json"))
				return
			}
			w.Write([]byte(j))
		default:
			w.Write([]byte(fmt.Sprintf("Error while mapping user: %v", err)))
		}
	}
}

func GetUserByEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		row := db.QueryRow("SELECT id, name, email, password, created_at, updated_at FROM public.users WHERE email=$1", email)
		var user User
		switch err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Created_at, &user.Updated_at); err {
		case sql.ErrNoRows:
			w.Write([]byte("No rows were returned"))
		case nil:
			j, err := json.Marshal(user)
			if err != nil {
				w.Write([]byte("error turning to json"))
				return
			}
			w.Write([]byte(j))
		default:
			w.Write([]byte(fmt.Sprintf("Error while mapping user: %v", err)))
		}
	}
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user CreatedUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var id int
		err = db.QueryRow(`INSERT INTO public.users (name, email, password, created_at, updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id`, user.Name, user.Email, user.Password, time.Now(), time.Now()).Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"ID\":%d, \"Success\": %t}", id, true)))
	}
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row := db.QueryRow("SELECT id, name, email, password FROM public.users WHERE id=$1", user.Id)
		var foundUser User
		switch err := row.Scan(&foundUser.Id, &foundUser.Name, &foundUser.Email, &foundUser.Password); err {
		case sql.ErrNoRows:
			w.Write([]byte("No user exists"))
		case nil:
			row := db.QueryRow("UPDATE public.users SET name = $1, password = $2, email=$3 WHERE id = $4", user.Name, user.Password, user.Email, user.Id)
			if row.Err() != nil {
				w.Write([]byte(fmt.Sprintf("Error updating user: %v", err)))
			}
			w.Write([]byte(fmt.Sprintf("{\"Success\": %t}", true)))
		default:
			w.Write([]byte(fmt.Sprintf("Error while mapping user: %v", err)))
		}
	}
}
