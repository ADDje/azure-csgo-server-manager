package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,string"`
}

// JSON write json response. Doesn't need to be public but otherwise namespaces collide
func JSON(w http.ResponseWriter, resp JSONResponse) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %s", err)
	}
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
	} else {
		resp.Success = true
	}

	JSON(w, resp)
}

// GetDefaultServerConfig Returns JSON response of default server config
func GetDefaultServerConfig(w http.ResponseWriter, r *http.Request) {

	resp := JSONResponse{
		Success: false,
		Data:    GetDefaultSettings(),
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	resp.Success = true

	JSON(w, resp)
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
	JSON(w, resp)
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
	JSON(w, resp)
}

func GetServerConfigTextByName(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var err error
	vars := mux.Vars(r)

	resp.Data, err = GetServerConfigTextFromFile(vars["configName"])
	if err != nil {
		log.Printf("Error getting server config %s", err)
		return
	}

	resp.Success = true
	JSON(w, resp)
}

func UpdateServerConfigText(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	vars := mux.Vars(r)
	fileName := CONFIG_DIRECTORY + vars["configName"]

	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error reading file %s: %s", fileName, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error in starting updating config handler body: %s", err)
		return
	}

	var valid bool
	valid, err = CheckConfigValid(body)
	if !valid {
		log.Printf("Error parsing file: %s", err)
		resp.Data = fmt.Sprintf("%s", err)
		JSON(w, resp)
		return
	}

	err = ioutil.WriteFile(fileName, body, 0770)

	if err != nil {
		log.Printf("Error writing file %s: %s", fileName, err)
		return
	}

	resp.Data = fmt.Sprintf("Config: %s, edited successfully", vars["configName"])
	resp.Success = true

	JSON(w, resp)
}

func CreateServerConfig(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	vars := mux.Vars(r)
	fileName := CONFIG_DIRECTORY + vars["configName"]

	err := ioutil.WriteFile(fileName, nil, 0770)

	if err != nil {
		log.Printf("Error creating file %s: %s", fileName, err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func DeleteServerConfig(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	vars := mux.Vars(r)
	fileName := CONFIG_DIRECTORY + vars["configName"]

	err := os.Remove(fileName)

	if err != nil {
		log.Printf("Error deleting file %s: %s", fileName, err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func GetDeploymentTemplates(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	data, err := GetTemplates()

	if err == nil {
		resp.Success = true
		resp.Data = data
	}

	JSON(w, resp)
}

func UpdateTemplateParameters(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	fileName := TEMPLATE_DIRECTORY + vars["templateName"] + ".parameters.json"

	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error reading file %s: %s", fileName, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error in starting updating template parameters handler body: %s", err)
		return
	}

	var valid bool
	valid, err = CheckTemplateValid(body)

	if !valid {
		log.Printf("Error parsing file: %s", err)
		resp.Data = fmt.Sprintf("%s", err)
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error parsing template parameters: %s", err)
		}
		return
	}

	err = ioutil.WriteFile(fileName, body, 0770)

	if err != nil {
		log.Printf("Error writing file %s: %s", fileName, err)
		return
	}

	resp.Data = fmt.Sprintf("Config: %s, edited successfully", vars["configName"])
	resp.Success = true

	JSON(w, resp)
}

func UpdateTemplateText(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	fileName := TEMPLATE_DIRECTORY + vars["templateName"] + ".json"

	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error reading file %s: %s", fileName, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error in starting updating template handler body: %s", err)
		return
	}

	var valid bool
	valid, err = CheckTemplateValid(body)

	if !valid {
		log.Printf("Error parsing file: %s", err)
		resp.Data = fmt.Sprintf("%s", err)
		JSON(w, resp)
		return
	}

	err = ioutil.WriteFile(fileName, body, 0770)

	if err != nil {
		log.Printf("Error writing file %s: %s", fileName, err)
		return
	}

	resp.Data = fmt.Sprintf("Config: %s, edited successfully", vars["templateName"])
	resp.Success = true

	JSON(w, resp)
}

func CreateDeploymentTemplate(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	vars := mux.Vars(r)
	templateFileName := TEMPLATE_DIRECTORY + vars["templateName"] + ".json"
	parameterFileName := TEMPLATE_DIRECTORY + vars["templateName"] + ".parameters.json"

	err := ioutil.WriteFile(templateFileName, []byte("{}"), 0770)

	if err != nil {
		log.Printf("Error creating template file %s: %s", templateFileName, err)
		return
	}

	err = ioutil.WriteFile(parameterFileName, []byte("{}"), 0770)

	if err != nil {
		log.Printf("Error creating parameter file %s: %s", parameterFileName, err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func DeleteDeploymentTemplate(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	vars := mux.Vars(r)

	templateFileName := TEMPLATE_DIRECTORY + vars["templateName"] + ".json"
	parameterFileName := TEMPLATE_DIRECTORY + vars["templateName"] + ".parameters.json"

	err := os.Remove(templateFileName)

	if err != nil {
		log.Printf("Error deleting template %s: %s", templateFileName, err)
		return
	}

	err = os.Remove(parameterFileName)

	if err != nil {
		log.Printf("Error deleting parameters %s: %s", parameterFileName, err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var user User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in starting csgo server handler body: %s", err)
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
		JSON(w, resp)
		return
	}

	log.Printf("User: %s, logged in successfully", user.Username)
	resp.Data = fmt.Sprintf("User: %s, logged in successfully", user.Username)
	resp.Success = true

	JSON(w, resp)
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
	JSON(w, resp)
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
		JSON(w, resp)
		return
	}

	resp.Success = true
	resp.Data = user

	JSON(w, resp)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	users, err := Auth.listUsers()
	if err != nil {
		log.Printf("Error in ListUsers handler: %s", err)
		resp.Data = fmt.Sprint("Error listing users")
		resp.Success = false
	} else {
		resp.Success = true
		resp.Data = users
	}

	JSON(w, resp)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	user := User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in reading add user POST: %s", err)
		resp.Data = fmt.Sprintf("Error in adding user: %s", err)
		resp.Success = false
		JSON(w, resp)
		return
	}

	log.Printf("Adding user: %v", string(body))

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("Error unmarshaling user add JSON: %s", err)
		resp.Data = fmt.Sprintf("Error in adding user: %s", err)
		resp.Success = false
		JSON(w, resp)
		return
	}

	err = Auth.addUser(user.Username, user.Password, user.Email, user.Role)
	if err != nil {
		log.Printf("Error in adding user: %s", err)
		resp.Data = fmt.Sprintf("Error in adding user: %s", err)
		resp.Success = false
		JSON(w, resp)
		return
	}

	resp.Success = true
	resp.Data = fmt.Sprintf("User: %s successfully added.", user.Username)

	JSON(w, resp)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	user := User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in reading remove user POST: %s", err)
		resp.Data = fmt.Sprintf("Error in removing user: %s", err)
		resp.Success = false
		JSON(w, resp)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("Error unmarshaling user remove JSON: %s", err)
		resp.Data = fmt.Sprintf("Error in removing user: %s", err)
		resp.Success = false
		JSON(w, resp)
		return
	}

	err = Auth.removeUser(user.Username)
	if err != nil {
		log.Printf("Error in remove user handler: %s", err)
		resp.Data = fmt.Sprintf("Error in removing user: %s", err)
		resp.Success = false
		JSON(w, resp)
		return
	}

	resp.Success = true
	resp.Data = fmt.Sprintf("User: %s successfully removed.", user.Username)

	JSON(w, resp)
}

// Return JSON response of server-settings.json file
func GetServerSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	//resp.Data = FactorioServ.Settings
	resp.Success = true

	JSON(w, resp)
}

func UpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in reading server settings POST: %s", err)
		resp.Data = fmt.Sprintf("Error in updating settings: %s", err)
		resp.Success = false
		JSON(w, resp)
		return
	}
	log.Printf("Received settings JSON: %s", body)

}
