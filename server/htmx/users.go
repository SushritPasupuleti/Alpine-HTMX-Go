package htmx

import (
	"net/http"
	// "server/helpers"

	"html/template"

	"server/helpers"
	"server/middleware"
	"server/models"

	"github.com/rs/zerolog/log"

	"bytes"
)

// Returns a HTML template with a list of users
func GetAllUsers(w http.ResponseWriter, r *http.Request, users []*models.User) {
	// log.Info().Msg("Accepts text/html")

	t, _ := template.ParseFiles("templates/users-list.html")

	var buf bytes.Buffer

	err := t.Execute(&buf, users)

	// log.Info().Msgf("buf: %s", buf.String())

	//save to cache
	middleware.SaveToCacheRaw(r, buf.String())

	helpers.WriteHTML(w, http.StatusOK, buf)

	if err != nil {
		log.Error().Err(err).Msg("Error executing template")
		w.Write([]byte("Error executing template"))
		return
	}
}
