package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fahriedr/golang-tdl/models"
	"github.com/fahriedr/golang-tdl/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	taskCollection *mongo.Collection
	dbName         string
	port           string
}

// var validate *validator.Validate
var Validate = validator.New()

func NewApiServer(taskCollection *mongo.Collection, dbName string, port string) *APIServer {
	return &APIServer{
		taskCollection: taskCollection,
		dbName:         dbName,
		port:           port,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	router.HandleFunc("/task/create", s.handleCreateTask).Methods(http.MethodPost)
	router.HandleFunc("/task/status/update", s.handleUpdateTaskStatus).Methods(http.MethodPost)
	router.HandleFunc("/task", s.handleGetTasks).Methods(http.MethodGet)
	router.HandleFunc("/task/{id}", s.handleGetTask).Methods(http.MethodGet)
	router.HandleFunc("/task/delete/{id}", s.handleDeleteTask).Methods(http.MethodDelete)

	log.Println("Listening on", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), router)
}

func (s *APIServer) handleCreateTask(w http.ResponseWriter, r *http.Request) {

	// Get context
	var ctx = r.Context()

	var payload models.TaskPayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = Validate.Struct(payload)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// generate unique id
	uniqueId := utils.RandomString(12)

	// initialize doc
	doc := models.Task{
		UniqueId:    uniqueId,
		Title:       payload.Title,
		Description: payload.Description,
		Status:      payload.Status,
		CreatedAt:   time.Now(),
	}

	// insert doc
	res, err := s.taskCollection.InsertOne(ctx, doc)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// type for return data
	var insertedDoc models.Task

	// filter for get one doc
	filter := bson.M{"_id": res.InsertedID}

	// get one doc
	err = s.taskCollection.FindOne(ctx, filter).Decode(&insertedDoc)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"message": "Success", "data": insertedDoc})
}

func (s *APIServer) handleUpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	var payload models.TaskStatusPayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	err = Validate.Struct(payload)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var value models.TaskPayload

	filter := bson.M{"uniqueId": payload.UniqueId}
	update := bson.M{"$set": bson.M{"status": payload.Status}}

	err = s.taskCollection.FindOneAndUpdate(ctx, filter, update).Decode(&value)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = bson.MarshalExtJSON(value, false, false)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Success"})

}

func (s *APIServer) handleGetTasks(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")

	ctx := r.Context()

	filter := bson.M{}

	if query != "" {

		filter = bson.M{
			"title": bson.M{
				"$regex":   query,
				"$options": "i",
			},
		}
	}

	var tasks []models.Task

	cursor, err := s.taskCollection.Find(ctx, filter)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err = cursor.All(ctx, &tasks); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Success", "data": tasks})

}

func (s *APIServer) handleGetTask(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	q := mux.Vars(r)

	var task models.Task

	filter := bson.M{"uniqueId": q["id"]}

	err := s.taskCollection.FindOne(ctx, filter).Decode(&task)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Success", "data": task})

}

func (s *APIServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	q := mux.Vars(r)

	filter := bson.M{"uniqueId": q["id"]}

	var task models.Task

	err := s.taskCollection.FindOneAndDelete(ctx, filter).Decode(&task)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Data successfully deleted"})

}
