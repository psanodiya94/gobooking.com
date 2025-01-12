package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/driver"
	"github.com/psanodiya94/gobooking.com/internal/forms"
	"github.com/psanodiya94/gobooking.com/internal/helpers"
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	checkIn := r.Form.Get("check_in")
	checkOut := r.Form.Get("check_out")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	checkinDate, err := time.Parse(layout, checkIn)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't parse check out date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	checkoutDate, err := time.Parse(layout, checkOut)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "can't parse check in date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(checkinDate, checkoutDate)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "no room available!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var res models.Reservation

	room, err := repo.DB.GetRoomById(roomId)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't get room id from database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := repo.DB.GetRoomById(res.RoomId)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't find room!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	layout := "2006-01-02"

	sd := r.Form.Get("check_in")
	ed := r.Form.Get("check_out")

	checkIn, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse check in date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	checkOut, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse check out date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse room id date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// send notifications - first to guest
	htmlMessage := fmt.Sprintf(
		`<strong>Reservation confirmation</strong><br>
		Dear %s, <br>
		This is to confirm your reservation from %s to %s.`,
		reservation.FirstName,
		reservation.CheckIn.Format("2006-01-02"),
		reservation.CheckOut.Format("2006-01-02"),
	)

	repo.App.MailChan <- models.MailData{
		To:       reservation.Email,
		From:     "gobookings@mailhog.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.email.html",
	}

	// send notifications - property owner
	ownerMessage := fmt.Sprintf(
		`<strong>Reservation Notification</strong><br>
		A reservation has been made for %s from %s to %s.`,
		reservation.FirstName,
		reservation.CheckIn.Format("2006-01-02"),
		reservation.CheckOut.Format("2006-01-02"),
	)

	repo.App.MailChan <- models.MailData{
		To:       "gobookings@mailhog.com",
		From:     "gobookings@mailhog.com",
		Subject:  "Reservation Notification",
		Content:  ownerMessage,
		Template: "basic.email.html",
	}

	repo.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays the reservation summary
func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

// GetShowLogin displays the login page
func (repo *Repository) GetShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin is a login handler for post login
func (repo *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		_ = render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := repo.DB.Authenticate(email, password)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(r.Context(), "user_id", id)
	repo.App.Session.Put(r.Context(), "flash", "Login successful")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GetLogout logs the user out
func (repo *Repository) GetLogout(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.Destroy(r.Context())
	_ = repo.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// GetAdminDashboard displays the admin dashboard
func (repo *Repository) GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// GetAdminAllReservations displays the admin all reservations
func (repo *Repository) GetAdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := repo.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	_ = render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// GetAdminNewReservations displays the admin new reservations
func (repo *Repository) GetAdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := repo.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	_ = render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// GetAdminShowReservation displays the admin show reservation
func (repo *Repository) GetAdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap := make(map[string]string)
	stringMap["src"] = src
	stringMap["year"] = year
	stringMap["month"] = month

	// get reservation by id
	res, err := repo.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = res

	_ = render.Template(w, r, "admin-show-reservation.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

// PostAdminShowReservation is the admin show reservation handler for post request
func (repo *Repository) PostAdminShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	// get reservation by id
	res, err := repo.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = repo.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.Form.Get("year")
	month := r.Form.Get("month")

	repo.App.Session.Put(r.Context(), "flash", "Changes saved")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

// GetAdminProcessReservation is the admin process reservation handler for get request
func (repo *Repository) GetAdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := chi.URLParam(r, "src")

	err = repo.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	repo.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// GetAdminDeleteReservation is the admin delete reservation handler for get request
func (repo *Repository) GetAdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := chi.URLParam(r, "src")

	err = repo.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	repo.App.Session.Put(r.Context(), "flash", "Reservation Deleted")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// GetAdminReservationsCalendar displays the admin reservations calendar
func (repo *Repository) GetAdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// assume there is no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := repo.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get all reservations
	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		// get all the restrictions for the current room
		restrictions, err := repo.DB.GetRestrictionsForRoomByDate(x.Id, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationId > 0 {
				// it's a reservation
				for d := y.CheckIn; d.After(y.CheckOut) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationId
				}
			} else {
				// it's a block
				blockMap[y.CheckIn.Format("2006-01-2")] = y.Id
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.Id)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.Id)] = blockMap

		repo.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.Id), blockMap)
	}

	_ = render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		IntMap:    intMap,
		Data:      data,
	})
}

// PostAdminReservationsCalendar displays the admin show reservation calendar
func (repo *Repository) PostAdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, err := strconv.Atoi(r.Form.Get("y"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month, err := strconv.Atoi(r.Form.Get("m"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := repo.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {
		// Get the block map from the session
		blockMap := repo.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.Id)).(map[string]int)

		// Loop through the block map
		for key, value := range blockMap {
			// ok will be false if the value is not in the map
			if val, ok := blockMap[key]; ok {
				// Only work on the values > 0, and that are not in the form post
				// The rest are the blocks we have
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.Id, key)) {
						// delete the restriction by id
						err := repo.DB.DeleteBlockById(value)
						if err != nil {
							helpers.ServerError(w, err)
							return
						}
					}
				}
			}
		}
	}

	// handle new blocks
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomId, err := strconv.Atoi(exploded[2])
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

			// get the date from the form
			startDate, err := time.Parse("2006-01-2", exploded[3])
			if err != nil {
				helpers.ServerError(w, err)
				return
			}

			// insert a new block
			err = repo.DB.InsertBlockForRoom(roomId, startDate)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
		}
	}

	repo.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
