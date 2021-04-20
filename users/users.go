package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/database"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    uint32 `json:"id"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, error := ioutil.ReadAll(r.Body)

	if error != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error"))
		return
	}

	var user user

	if error = json.Unmarshal(body, &user); error != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao converter usuário para o struct"))
		return
	}

	db, erro := database.Connection()

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao conectar ao database"))
		return
	}

	defer db.Close()

	statement, erro := db.Prepare("insert into users (name, email) values (?, ?)")

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao criar o statement"))
		return
	}
	defer statement.Close()

	insert, erro := statement.Exec(user.Name, user.Email)

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error ao executar o Exec"))
		return
	}

	userId, erro := insert.LastInsertId()

	if erro != nil {
		w.WriteHeader(404)
		w.Write([]byte("Erro ao obter o id"))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("Id: %d", userId)))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, erro := database.Connection()

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao realizar conexão com o banco"))
		return
	}

	defer db.Close()

	response, erro := db.Query("select * from users")

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao realizar consulta"))
		return
	}

	defer response.Close()

	var users []user
	for response.Next() {
		var user user

		if erro := response.Scan(&user.Id, &user.Name, &user.Email); erro != nil {
			w.WriteHeader(400)
			w.Write([]byte("Erro a carregar dados SCAN"))
			return
		}
		users = append(users, user)
	}
	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(users); erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error ao converter JSON"))
		return
	}

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro"))
		return
	}

	db, erro := database.Connection()

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}
	defer db.Close()

	response, erro := db.Query("select * from users where id = ?", id)

	if erro != nil {
		w.WriteHeader(404)
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}

	var user user

	if response.Next() {
		if erro := response.Scan(&user.Id, &user.Name, &user.Email); erro != nil {
			w.WriteHeader(400)
			w.Write([]byte("Erro ao carregar dados SCAN"))
			return
		}
	}

	if user.Id == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Usuário não encontrado!"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(user); erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro buscar usuário"))
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro"))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao ler o body"))
		return
	}

	var user user

	if erro := json.Unmarshal(body, &user); erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao converter JSON"))
		return
	}

	db, erro := database.Connection()

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao conectar com banco"))
		return
	}

	statement, erro := db.Prepare("update users set name = ?, email = ? where id = ?")
	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao criar o statement"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(user.Name, user.Email, id); erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao atualizar usuário"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, erro := strconv.ParseUint(params["id"], 10, 32)

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao obter o id"))
		return
	}

	db, erro := database.Connection()

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}

	defer db.Close()

	statement, erro := db.Prepare("delete from users where id = ?")

	if erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro realizar o statement"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(id); erro != nil {
		w.WriteHeader(400)
		w.Write([]byte("Erro ao deletar usuário"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
