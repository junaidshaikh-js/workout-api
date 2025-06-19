package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/junaidshaikh-js/workout-api/internal/store"
	"github.com/junaidshaikh-js/workout-api/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)

	if err != nil {
		h.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout id"})
	}

	workout, err := h.workoutStore.GetWorkoutByID(workoutId)

	if err != nil {
		h.logger.Printf("ERROR: getWorkoutById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"workout": workout})
}

func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		h.logger.Printf("ERROR: decodeCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request"})
	}

	createdWorkout, err := h.workoutStore.CreateWorkout(&workout)

	if err != nil {
		h.logger.Printf("ERROR: createWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "failed to create workout"})
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelop{"workout": createdWorkout})
}

func (h *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)

	if err != nil {
		h.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout id"})
	}

	existingWorkout, err := h.workoutStore.GetWorkoutByID(workoutId)

	if err != nil {
		h.logger.Printf("ERROR: getWorkoutById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "failed to fetch workout"})
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		h.logger.Printf("ERROR: decodeError: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request body"})
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	err = h.workoutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		h.logger.Printf("ERROR: updateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"workout": existingWorkout})
}

func (h *WorkoutHandler) HandleDeleteById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)

	if err != nil {
		h.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout id"})
	}

	err = h.workoutStore.DeleteWorkout(workoutId)

	if err == sql.ErrNoRows {
		h.logger.Printf("ERROR: deleteWorkout: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelop{"error": "workout not found"})
	}

	if err != nil {
		h.logger.Printf("ERROR: deleteWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelop{})
}
