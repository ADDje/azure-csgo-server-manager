package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,string"`
}

// GetAllServers Returns JSON response of all servers found
func GetAllServers(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp.Data, err = getServers(config)
	if err != nil {
		resp.Data = fmt.Sprintf("Error in GetAllServers handler: %s", err)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in list mods: %s", err)
		}
		return
	}

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in get servers: %s", err)
	}
}

// GetDefaultServerConfig Returns JSON response of default server config
func GetDefaultServerConfig(w http.ResponseWriter, r *http.Request) {

	resp := JSONResponse{
		Success: false,
		Data:    GetDefaultSettings(),
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error in get default config: %s", err)
	}
}

func GetServerConfigs(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var err error
	resp.Data, err = GetServerConfigsFromFile()
	if err != nil {
		log.Printf("Error getting server configs")
		return
	}

	resp.Success = true
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding server configs: %s", err)
	}
}

func GetServerConfigByName(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var err error
	vars := mux.Vars(r)

	resp.Data, err = GetServerConfigFromFile(vars["configName"])
	if err != nil {
		log.Printf("Error getting server config %s", err)
		return
	}

	resp.Success = true
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding server config: %s", err)
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for login handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	case "POST":
		var user User
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in starting factorio server handler body: %s", err)
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("Error unmarshaling server settings JSON: %s", err)
			return
		}

		log.Printf("Logging in user: %s", user.Username)

		err = Auth.aaa.Login(w, r, user.Username, user.Password, "/")
		if err != nil {
			log.Printf("Error logging in user: %s, error: %s", user.Username, err)
			resp.Data = fmt.Sprintf("Error logging in user: %s", user.Username)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error listing mods: %s", err)
			}
			return
		}

		log.Printf("User: %s, logged in successfully", user.Username)
		resp.Data = fmt.Sprintf("User: %s, logged in successfully", user.Username)
		resp.Success = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if err := Auth.aaa.Logout(w, r); err != nil {
		log.Printf("Error logging out current user")
		return
	}

	resp.Success = true
	resp.Data = "User logged out successfully."
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error logging out: %s", err)
	}
}

func GetCurrentLogin(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	user, err := Auth.aaa.CurrentUser(w, r)
	if err != nil {
		log.Printf("Error getting current user status: %s", err)
		resp.Data = fmt.Sprintf("Error getting user status: %s", user.Username)
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
		return
	}

	resp.Success = true
	resp.Data = user

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error getting user status: %s", err)
	}

}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	users, err := Auth.listUsers()
	if err != nil {
		log.Printf("Error in ListUsers handler: ", err)
		resp.Data = fmt.Sprint("Error listing users")
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error listing mods: %s", err)
		}
		return
	}

	resp.Success = true
	resp.Data = users

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error getting user status: %s", err)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for add user handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error adding user: %s", err)
		}
	case "POST":
		user := User{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in reading add user POST: %s", err)
			resp.Data = fmt.Sprintf("Error in adding user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		log.Printf("Adding user: %v", string(body))

		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("Error unmarshaling user add JSON: %s", err)
			resp.Data = fmt.Sprintf("Error in adding user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		err = Auth.addUser(user.Username, user.Password, user.Email, user.Role)
		if err != nil {
			log.Printf("Error in adding user: %s", err)
			resp.Data = fmt.Sprintf("Error in adding user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		resp.Success = true
		resp.Data = fmt.Sprintf("User: %s successfully added.", user.Username)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in returning added user response: %s", err)
		}
	}
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for add user handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error adding user: %s", err)
		}
	case "POST":
		user := User{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in reading remove user POST: %s", err)
			resp.Data = fmt.Sprintf("Error in removing user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("Error unmarshaling user remove JSON: %s", err)
			resp.Data = fmt.Sprintf("Error in removing user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error removing user: %s", err)
			}
			return
		}

		err = Auth.removeUser(user.Username)
		if err != nil {
			log.Printf("Error in remove user handler: %s", err)
			resp.Data = fmt.Sprintf("Error in removing user: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error adding user: %s", err)
			}
			return
		}

		resp.Success = true
		resp.Data = fmt.Sprintf("User: %s successfully removed.", user.Username)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error in returning remove user response: %s", err)
		}
	}
}

// Return JSON response of server-settings.json file
func GetServerSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//resp.Data = FactorioServ.Settings
	resp.Success = true

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding server settings JSON reponse: ", err)
	}

	log.Printf("Sent server settings response")
}

func UpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	switch r.Method {
	case "GET":
		log.Printf("GET not supported for add user handler")
		resp.Data = "Unsupported method"
		resp.Success = false
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error adding user: %s", err)
		}
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error in reading server settings POST: %s", err)
			resp.Data = fmt.Sprintf("Error in updating settings: %s", err)
			resp.Success = false
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Error updating settings: %s", err)
			}
			return
		}
		log.Printf("Received settings JSON: %s", body)
	}
}
