package handlers

import (
	"encoding/json"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/driver"
	"github.com/psanodiya94/gobooking.com/internal/forms"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"github.com/psanodiya94/gobooking.com/internal/render"
	"github.com/psanodiya94/gobooking.com/internal/repository"
	"github.com/psanodiya94/gobooking.com/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DBRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DataBase) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new testing repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the homepage handler
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Contact is the contact page handler
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Majors is the majors page handler
func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Generals is the generals page handler
func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// GetAvailability checks the availability of rooms for get request
func (repo *Repository) GetAvailability(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability checks the availability of rooms for post request
func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	checkIn := r.Form.Get("check_in")
	checkOut := r.Form.Get("check_out")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	checkinDate, err := time.Parse(layout, checkIn)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't parse check out date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	checkoutDate, err := time.Parse(layout, checkOut)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't parse check in date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(checkinDate, checkoutDate)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "no room available!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		repo.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		CheckIn:  checkinDate,
		CheckOut: checkoutDate,
	}

	repo.App.Session.Put(r.Context(), "reservation", res)

	_ = render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// JsonAvailability checks the availability of rooms for post request with ajax response
func (repo *Repository) JsonAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error parsing start date",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error parsing end date",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error parsing room id",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}

	available, err := repo.DB.SearchAvailabilityForDatesByRoomId(roomId, startDate, endDate)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomId),
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}

// ChooseRoom Display list of rooms
func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	roomId, err := strconv.Atoi(exploded[2])
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomId = roomId

	repo.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL params, build a session variable, and takes user to make-reservation screen
func (repo *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var res models.Reservation

	room, err := repo.DB.GetRoomById(roomId)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't get room id from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomId = roomId
	res.CheckIn = startDate
	res.CheckOut = endDate

	repo.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// GetReservation is the reservation page handler for get request
func (repo *Repository) GetReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := repo.DB.GetRoomById(res.RoomId)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	repo.App.Session.Put(r.Context(), "reservation", res)

	checkIn := res.CheckIn.Format("2006-01-02")
	checkOut := res.CheckOut.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["check_in"] = checkIn
	stringMap["check_out"] = checkOut

	data := make(map[string]interface{})
	data["reservation"] = res

	_ = render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation is the reservation page handler for post request
func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	layout := "2006-01-02"

	sd := r.Form.Get("check_in")
	ed := r.Form.Get("check_out")

	checkIn, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse check in date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	checkOut, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse check out date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse room id date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		CheckIn:   checkIn,
		CheckOut:  checkOut,
		RoomId:    roomId,
	}

	reservation.Room.RoomName = r.Form.Get("room_name")

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		http.Error(w, "Form is not valid", http.StatusSeeOther)

		_ = render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	reservationId, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		CheckIn:       reservation.CheckIn,
		CheckOut:      reservation.CheckOut,
		RoomId:        reservation.RoomId,
		ReservationId: reservationId,
		RestrictionId: 1,
	}

	err = repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't insert room restriction into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repo.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays the reservation summary
func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repo.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	checkIn := reservation.CheckIn.Format("2006-01-02")
	checkOut := reservation.CheckOut.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["check_in"] = checkIn
	stringMap["check_out"] = checkOut

	_ = render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}
