package song

import (
	"net/http"
)

// Song HTTP Handlers interface
type Handlers interface {
	GetList(w http.ResponseWriter, r *http.Request)
	GetText(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
}
