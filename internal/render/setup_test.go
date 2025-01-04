package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/psanodiya94/gobooking.com/internal/config"
	"github.com/psanodiya94/gobooking.com/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am i going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type writerResponse struct{}

func (w *writerResponse) Header() http.Header {
	return http.Header{}
}

func (w *writerResponse) WriteHeader(statusCode int) {}

func (w *writerResponse) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
