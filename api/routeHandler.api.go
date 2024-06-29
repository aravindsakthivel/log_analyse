package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log_analyse/DB"
	"log_analyse/controller"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SUserDetailPR struct {
	Name      string `json:"name"`
	UserEmail string `json:"userEmail"`
}

type SRoutes struct {
	dbCtrl *controller.SDBCtrl
}

func (sr *SRoutes) setCtrl(dbCtrl *controller.SDBCtrl) {
	sr.dbCtrl = dbCtrl
}

func (sr *SRoutes) main(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func (sr *SRoutes) health(w http.ResponseWriter, r *http.Request) {
	dbProp := DB.SDB{}
	var health bool = dbProp.Health()
	if !health {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database is not healthy"))
		return
	}
	w.Write([]byte("Database is healthy"))
}

func (sr *SRoutes) createUserDummy(w http.ResponseWriter, r *http.Request) {
	userCtrl, err := sr.dbCtrl.Users()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var filesUploaded []primitive.ObjectID = []primitive.ObjectID{}
	user := controller.SUserInfo{
		Name:          "John Doe",
		UserEmail:     "john@mail.com",
		FilesUploaded: filesUploaded,
	}
	err = userCtrl.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error())) // Convert error message to byte slice
		return
	}
	w.WriteHeader(http.StatusOK)
	var response []byte = []byte(fmt.Sprintf("Inserted User %s", user.UserEmail))
	w.Write(response)
}

func (sr SRoutes) createUser(w http.ResponseWriter, r *http.Request) {
	userCtrl, err := sr.dbCtrl.Users()
	if err != nil {
		log.Println("Error in getting the Users DB control ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user SUserDetailPR

	errBD := json.NewDecoder(r.Body).Decode(&user)

	if errBD != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(errBD.Error())) // Convert error message to byte slice
		return
	}

	log.Println("Body ", user)

	var filesUploaded []primitive.ObjectID = []primitive.ObjectID{}
	buildUser := controller.SUserInfo{
		Name:          user.Name,
		UserEmail:     user.UserEmail,
		FilesUploaded: filesUploaded,
	}

	log.Println("buildUser ", buildUser)

	err = userCtrl.CreateUser(buildUser)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error())) // Convert error message to byte slice
		return
	}

	w.WriteHeader(http.StatusOK)
	var response []byte = []byte(fmt.Sprintf("Inserted User %s", buildUser.UserEmail))
	w.Write(response)
}

func (sr *SRoutes) uploadFile(w http.ResponseWriter, r *http.Request) {

	userCtrl, errDB := sr.dbCtrl.Users()

	fileCtrl, fileErr := sr.dbCtrl.Files()

	if errDB != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errDB.Error()))
		return
	}

	if fileErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fileErr.Error()))
		return
	}

	// Parse the multipart form in the request
	err := r.ParseMultipartForm(10 << 20) // limit your maxMultipartMemory
	if err != nil {
		log.Println("Error in parsing the form data ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the file from form data
	file, header, err := r.FormFile("file") // retrieve the file from form data
	if err != nil {
		log.Println("Error in retrieving the file ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// print current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting the current working directory ", err)
		http.Error(w, "Error getting the current working directory", http.StatusInternalServerError)
		return
	}

	log.Println("Current working directory: ", cwd)

	os.Mkdir("./uploads", os.ModePerm)

	splitUd := strings.Split(uuid.New().String(), "-")

	filesPrefixID := strings.Join(splitUd, "")

	log.Println("Id ", filesPrefixID)

	filePath := "./uploads/" + filesPrefixID + "-" + header.Filename

	// Read the file into a byte slice
	destFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating the file ", err)
		http.Error(w, "Error creating the file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		log.Println("Error saving the file ", err)
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	var query = r.URL.Query()
	mail := query.Get("mail")

	if mail == "" {
		log.Println("Error getting mail id ", err)
		http.Error(w, "Error getting mail id ", http.StatusInternalServerError)
		return
	}

	fileId, setFileErr := fileCtrl.SetFiles(filePath, header.Filename, mail) // TODO : test this function

	if setFileErr != nil {
		log.Println("Error setting file ", err)
		http.Error(w, "Error setting file ", http.StatusInternalServerError)
		return
	}

	log.Print("Uploaded file ID ", fileId)

	errPF := userCtrl.PushFileName(mail, fileId)

	if errPF != nil {
		log.Println("Error saving the file path ", errPF)
		http.Error(w, "Error saving the file path ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func (sr *SRoutes) findUnCompressedFiles(w http.ResponseWriter, r *http.Request) {
	fileCtrl, fileDBErr := sr.dbCtrl.Files()

	if fileDBErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fileDBErr.Error()))
		return
	}
	filesExists, files, err := fileCtrl.GetUnCompressedFile()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if !filesExists {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No un compressed files found"))
		return
	}

	filesJSON, err := json.Marshal(files)
	if err != nil {
		log.Println("Error marshalling files: ", err)
		http.Error(w, "Error processing files", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	var response []byte = []byte(filesJSON)
	w.Write(response)
}
