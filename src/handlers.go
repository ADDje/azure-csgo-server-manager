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
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
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

	resp.Data, err = GetVms(config)
	if err != nil {
		resp.Data = fmt.Sprintf("Error in GetAllServers handler: %s", err)
	} else {
		resp.Success = true
	}

	JSON(w, resp)
}

func DeployServers(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := JSONResponse{
		Success: false,
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Invalid Body: %s", err)
		return
	}

	argumentsJSON := make(map[string]interface{})
	err = json.Unmarshal(body, &argumentsJSON)
	if err != nil {
		log.Printf("Invalid arguments json: %s", err)
		return
	}

	numberOfServers, _ := argumentsJSON["numberOfServers"].(float64)
	vmName, _ := argumentsJSON["vmName"].(string)
	adminUserName, _ := argumentsJSON["adminUserName"].(string)
	adminPassword, _ := argumentsJSON["adminPassword"].(string)
	configFile, _ := argumentsJSON["configFile"].(string)
	templateFile, _ := argumentsJSON["templateFile"].(string)

	DeployXTemplates(int(numberOfServers), config, vmName, adminUserName,
		adminPassword, configFile, templateFile)

	resp.Success = true
	JSON(w, resp)
}

// StartServer starts an existing VM that has been stopped (deallocated)
func StartServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	log.Printf("Starting Server: %s", vars["vmName"])
	err := StartVM(config, vars["vmName"])

	if err != nil {
		resp.Data = err
	} else {
		resp.Success = true
	}

	JSON(w, resp)
}

// StopServer deallocates a VM, as just "stopping" it retains resources
func StopServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	log.Printf("Stopping (deallocating) Server: %s", vars["vmName"])
	err := DeallocateVM(config, vars["vmName"])

	if err != nil {
		resp.Data = err
	} else {
		resp.Success = true
	}

	JSON(w, resp)
}

// DeleteServer deletes a VM and associated components (network, IP)
func DeleteServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	log.Printf("Beginning deletion of VM: %s", vars["vmName"])
	err := FullDeleteVM(config, vars["vmName"])

	if err != nil {
		resp.Data = err
	} else {
		resp.Success = true
		log.Printf("VM Successfully deleted: %s", vars["vmName"])
	}

	JSON(w, resp)
}

// ReplayServer exports a server's replays (demos)
func ReplayServer(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Invalid Body: %s", err)
		return
	}
	bodyJSON := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		log.Printf("Invalid body JSON: %s", err)
		return
	}

	log.Printf("Exporting replays for VM: %s", vars["vmName"])
	err = ExportReplays(config, vars["vmName"], bodyJSON["username"].(string), bodyJSON["password"].(string))

	if err != nil {
		resp.Data = err
	} else {
		resp.Success = true
		log.Printf("Replays Successfully exported: %s", vars["vmName"])
	}

	JSON(w, resp)
}

// GetDefaultServerConfig Returns JSON response of default server config
func GetDefaultServerConfig(w http.ResponseWriter, r *http.Request) {

	resp := JSONResponse{
		Success: false,
		Data:    GetDefaultSettings(),
	}

	resp.Success = true

	JSON(w, resp)
}

// GetServerConfigs Returns list of server.conf files from cloud or local storage
func GetServerConfigs(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	var err error
	if config.UseCloudStorage {
		resp.Data, err = GetServerConfigsFromAzure()
	} else {
		resp.Data, err = GetServerConfigsFromFile()
	}
	if err != nil {
		log.Printf("Error getting server configs")
		return
	}

	resp.Success = true
	JSON(w, resp)
}

// GetServerConfigByName Returns a server.conf file from cloud or local storage
func GetServerConfigByName(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	var err error
	vars := mux.Vars(r)

	if config.UseCloudStorage {
		resp.Data, err = GetStorageFile(config, CONFIG_FILE_STORE, vars["configName"])
	} else {
		resp.Data, err = GetServerConfigFromFile(vars["configName"])
	}
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

	var err error
	vars := mux.Vars(r)

	if config.UseCloudStorage {
		azureFile, err2 := GetStorageFile(config, CONFIG_FILE_STORE, vars["configName"])
		err = err2
		myBytes, err := ReadConfigIntoBytes(azureFile.Body)
		if err == nil {
			resp.Data = string(myBytes)
		}
	} else {
		resp.Data, err = GetServerConfigTextFromFile(vars["configName"])
	}
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

	vars := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error in updating server config handler body: %s", err)
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

	if config.UseCloudStorage {
		err = UpdateStorageFile(config, CONFIG_FILE_STORE, vars["configName"], body)
	} else {
		err = ioutil.WriteFile(CONFIG_DIRECTORY+vars["configName"], body, 0770)
	}

	if err != nil {
		log.Printf("Error writing file %s: %s", vars["configName"], err)
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

	vars := mux.Vars(r)

	var err error
	if config.UseCloudStorage {
		err = CreateStorageFile(config, CONFIG_FILE_STORE, vars["configName"], nil)
	} else {
		err = ioutil.WriteFile(CONFIG_DIRECTORY+vars["configName"], nil, 0770)
	}

	if err != nil {
		log.Printf("Error creating file %s", err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func DeleteServerConfig(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	var err error
	if config.UseCloudStorage {
		err = DeleteStorageFile(config, CONFIG_FILE_STORE, vars["configName"])
	} else {
		err = os.Remove(CONFIG_DIRECTORY + vars["configName"])
	}

	if err != nil {
		log.Printf("Error deleting file: %s", err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func GetDeploymentTemplates(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	var err error
	if config.UseCloudStorage {
		resp.Data, err = GetTemplatesFromAzure()
	} else {
		resp.Data, err = GetTemplatesFromFile()
	}

	if err == nil {
		resp.Success = true
	}

	JSON(w, resp)
}

func UpdateTemplateParameters(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	fileName := vars["templateName"] + ".parameters.json"

	var err error
	existingParameters := TemplateParameterFile{}
	if config.UseCloudStorage {
		azureFile, err := GetStorageFile(config, TEMPLATE_FILE_STORE, fileName)
		if err != nil {
			log.Printf("Could not read azure parameters file: %s", err)
			return
		}
		err = json.NewDecoder(azureFile.Body).Decode(&existingParameters)
		log.Printf(azureFile.Properties.ContentType)
	} else {
		file, err := ioutil.ReadFile(TEMPLATE_DIRECTORY + fileName)
		if err != nil {
			log.Printf("Could not read parameters file: %s", err)
			return
		}
		err = json.Unmarshal(file, &existingParameters)
	}
	if err != nil {
		log.Printf("Error reading file %s: %s", fileName, err)
		return
	}

	log.Printf("Existing Schema: %s", existingParameters.Schema)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error in updating template parameters handler body: %s", err)
		return
	}

	paramsFile, err := CheckParametersValid(body)
	if err != nil {
		return
	}

	// Replace existing parameters with their new ones
	existingParameters.Parameters = paramsFile.Parameters
	saveBody, err := json.Marshal(existingParameters)
	if err != nil {
		log.Printf("Could not convert parameters body: %s", err)
		return
	}

	if config.UseCloudStorage {
		err = UpdateStorageFile(config, TEMPLATE_FILE_STORE, fileName, saveBody)
	} else {
		err = ioutil.WriteFile(TEMPLATE_DIRECTORY+fileName, saveBody, 0770)
	}

	if err != nil {
		log.Printf("Error writing file %s: %s", fileName, err)
		return
	}

	resp.Data = fmt.Sprintf("Template: %s, edited successfully", fileName)
	resp.Success = true

	log.Printf("Updated template parameters: %s", saveBody)

	JSON(w, resp)
}

func UpdateTemplateText(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	fileName := vars["templateName"] + ".json"

	var err error
	if config.UseCloudStorage {
		_, err = GetStorageFile(config, TEMPLATE_FILE_STORE, fileName)
	} else {
		_, err = ioutil.ReadFile(TEMPLATE_DIRECTORY + fileName)
	}
	if err != nil {
		log.Printf("Error reading file %s: %s", fileName, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error in updating template handler body: %s", err)
		return
	}

	_, err = CheckTemplateValid(body)
	if err != nil {
		return
	}

	if config.UseCloudStorage {
		err = UpdateStorageFile(config, TEMPLATE_FILE_STORE, fileName, body)
	} else {
		err = ioutil.WriteFile(TEMPLATE_DIRECTORY+fileName, body, 0770)
	}

	if err != nil {
		log.Printf("Error writing file %s: %s", fileName, err)
		return
	}

	resp.Data = fmt.Sprintf("Template: %s, edited successfully", vars["templateName"])
	resp.Success = true

	JSON(w, resp)
}

func CreateDeploymentTemplate(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)
	templateFileName := vars["templateName"] + ".json"
	parameterFileName := vars["templateName"] + ".parameters.json"

	var err error
	if config.UseCloudStorage {
		err = CreateStorageFile(config, TEMPLATE_FILE_STORE, templateFileName, []byte("{}"))
	} else {
		err = ioutil.WriteFile(TEMPLATE_DIRECTORY+templateFileName, []byte("{}"), 0770)
	}

	if err != nil {
		log.Printf("Error creating template file %s: %s", templateFileName, err)
		return
	}

	if config.UseCloudStorage {
		err = CreateStorageFile(config, TEMPLATE_FILE_STORE, parameterFileName, GetDefaultParametersFile())
	} else {
		err = ioutil.WriteFile(TEMPLATE_DIRECTORY+parameterFileName, GetDefaultParametersFile(), 0770)
	}

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

	vars := mux.Vars(r)

	templateFileName := vars["templateName"] + ".json"
	parameterFileName := vars["templateName"] + ".parameters.json"

	var err error
	if config.UseCloudStorage {
		err = DeleteStorageFile(config, TEMPLATE_FILE_STORE, templateFileName)
	} else {
		err = os.Remove(TEMPLATE_DIRECTORY + templateFileName)
	}

	if err != nil {
		log.Printf("Error deleting template %s: %s", templateFileName, err)
		return
	}

	if config.UseCloudStorage {
		err = DeleteStorageFile(config, TEMPLATE_FILE_STORE, parameterFileName)
	} else {
		err = os.Remove(TEMPLATE_DIRECTORY + parameterFileName)
	}

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

	var user User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in Login User handler body: %s", err)
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

	if err := Auth.aaa.Logout(w, r); err != nil {
		log.Print("Error logging out current user")
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

// GetSettings Return JSON response of conf.json file
func GetSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: true,
		Data:    config,
	}

	JSON(w, resp)
}

// UpdateSettings updates the conf.json file
func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error in reading server settings POST: %s", err)
		resp.Data = fmt.Sprintf("Error in updating settings: %s", err)
		JSON(w, resp)
		return
	}
	log.Printf("Received settings JSON: %s", body)

	_, err = ioutil.ReadFile(config.ConfFile)
	if err != nil {
		log.Printf("Could not open config file %s: %s", config.ConfFile, err)
		resp.Data = fmt.Sprintf("Error in updating settings: %s", err)
		JSON(w, resp)
		return
	}

	// Decode into config
	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Printf("Could not unmarshal config: %s", err)
		return
	}

	log.Printf("Config: %s", config.ResourceGroup)

	var newJSON []byte
	newJSON, err = json.MarshalIndent(config, "", "    ")
	err = ioutil.WriteFile(config.ConfFile, newJSON, 0770)
	if err != nil {
		log.Printf("Could not write config %s: %s", config.ConfFile, err)
		resp.Data = fmt.Sprintf("Error in updating settings: %s", err)
		JSON(w, resp)
		return
	}

	resp.Success = true
	JSON(w, resp)
}

func GetAllActions(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	actions, err := GetScheduleActions()
	if err == nil {
		resp.Success = true
		resp.Data = actions
	}

	JSON(w, resp)
}

func CreateOrUpdateAction(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	name, ok := vars["actionName"]
	if !ok {
		log.Printf("Missing action name")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Invalid Body: %s", err)
		return
	}
	log.Printf("Body: %s", body)

	action := ScheduleAction{}
	err = json.Unmarshal(body, &action)
	if err != nil {
		log.Printf("Invalid schedule action: %s", err)
		return
	}

	err = AddOrUpdateScheduleAction(name, &action)
	if err != nil {
		log.Printf("Could not add or update action: %s", err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func DeleteAction(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	name, ok := vars["actionName"]
	if !ok {
		log.Printf("Missing action name")
		return
	}

	err := DeleteScheduleAction(name)
	if err != nil {
		log.Printf("Could delete action: %s", err)
		return
	}

	resp.Success = true

	JSON(w, resp)
}

func ExecuteAction(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Success: false,
	}

	vars := mux.Vars(r)

	name, ok := vars["actionName"]
	if !ok {
		log.Printf("Missing action name")
		return
	}

	params := r.URL.Query()

	// Flatten params, only support the first value for each param.
	// Easier to deal with in the scheduler
	newParams := make(map[string]string)
	for k, v := range params {
		newParams[k] = v[0]
	}

	log.Printf("Executing action: %s", name)
	err := ExecuteScheduleAction(name, newParams)
	if err != nil {
		resp.Data = fmt.Sprintf("Error: %s", err)
		log.Printf("Could execute action %s: %s", name, err)
	} else {
		resp.Success = true
	}

	JSON(w, resp)
}
