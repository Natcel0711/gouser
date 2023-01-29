package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Natcel0711/gouser/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id,name, username,email,password,created_at,updated_at FROM public.users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				JSONError(w, err.Error(), 400)
			}
		}(rows)
		var users []models.User
		for rows.Next() {
			var user models.User
			err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Created_at, &user.Updated_at)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(j)
	}
}

func GetUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		row := db.QueryRow("SELECT id, name, username, email, password, created_at, updated_at FROM public.users WHERE id=$1", id)
		var user models.User
		switch err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Created_at, &user.Updated_at); err {
		case sql.ErrNoRows:
			JSONError(w, "No rows returned", http.StatusFound)
		case nil:
			j, err := json.Marshal(user)
			if err != nil {
				JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(j)
		default:
			JSONError(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetUserBySession(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionid := chi.URLParam(r, "sessionid")
		if sessionid == "" {
			JSONError(w, "session id not available", 500)
			return
		}
		row := db.QueryRow("SELECT u.id, u.name, u.username, u.email, u.password, u.created_at, u.updated_at FROM public.users u INNER JOIN public.sessions s ON u.id = s.userid WHERE s.sessionid=$1", sessionid)
		var user models.User
		switch err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Created_at, &user.Updated_at); err {
		case sql.ErrNoRows:
			JSONError(w, "no rows", 500)
		case nil:
			j, err := json.Marshal(user)
			if err != nil {
				JSONError(w, "error turning to json", 500)
				return
			}
			_, err = w.Write(j)
			if err != nil {
				JSONError(w, "error writing to response", 500)
				return
			}
		default:
			JSONError(w, "error while mapping", 500)
		}
	}
}

func GetUserByEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		row := db.QueryRow("SELECT id, name, username, email, password, created_at, updated_at FROM public.users WHERE email=$1", email)
		fmt.Println(row)
		var user models.User
		switch err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password, &user.Created_at, &user.Updated_at); err {
		case sql.ErrNoRows:
			JSONError(w, "No users with that email", 404)
		case nil:
			j, err := json.Marshal(user)
			if err != nil {
				JSONError(w, "error turning to json", 500)
				return
			}
			_, err = w.Write(j)
			if err != nil {
				JSONError(w, "error writing to response", 500)
				return
			}
		default:
			JSONError(w, "error while mapping user", 500)
		}
	}
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var id int
		err = db.QueryRow(`INSERT INTO public.users (name, username, email, password, created_at, updated_at) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`, user.Name, user.Username, user.Email, user.Password, time.Now(), time.Now()).Scan(&id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte(fmt.Sprintf("{\"ID\":%d, \"Success\": %t}", id, true)))
		if err != nil {
			JSONError(w, "error writing to response", http.StatusBadRequest)
			return
		}
	}
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		row := db.QueryRow("SELECT id, name, username, email, password FROM public.users WHERE id=$1", user.Id)
		var foundUser models.User
		switch err := row.Scan(&foundUser.Id, &foundUser.Name, &foundUser.Username, &foundUser.Email, &foundUser.Password); err {
		case sql.ErrNoRows:
			JSONError(w, "no user exists", http.StatusBadRequest)
		case nil:
			row := db.QueryRow("UPDATE public.users SET name = $1, username = $2, password = $3, email=$4 WHERE id = $4", user.Name, user.Username, user.Password, user.Email, user.Id)
			if row.Err() != nil {
				JSONError(w, err.Error(), http.StatusBadRequest)
			}
			JSONError(w, err.Error(), http.StatusBadRequest)
		default:
			JSONError(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func CreateSessionID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var id int
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			JSONError(w, "error while encoding json", 500)
		}
		sqlStatement := `DELETE FROM public.sessions WHERE userid = $1;`
		_, err = db.Exec(sqlStatement, user.Id)
		if err != nil {
			JSONError(w, err.Error(), 500)
			return
		}
		session := uuid.New()
		err = db.QueryRow(`INSERT INTO public.sessions (sessionid, userid) VALUES($1,$2) returning id`, session, user.Id).Scan(&id)
		if err != nil {
			JSONError(w, "error inserting session", 500)
			return
		}
		_, err = w.Write([]byte(fmt.Sprintf("{\"ID\":\"%s\", \"Success\": %t}", session.String(), true)))
	}
}

func JSONError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"error\": true, \"message\":\"%s\"}", message)))
}
