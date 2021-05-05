package handlers

import (
	ctx2 "github.com/geo-provider/internal/services/api/ctx"
	render2 "github.com/geo-provider/internal/services/api/render"
	"github.com/geo-provider/internal/storage"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

const DefaultRawsCount = 100

func GetSources(w http.ResponseWriter, r *http.Request) {
	storagesList := ctx2.Config(r).ListSources()

	render2.Respond(w, http.StatusOK, render2.Message(storagesList))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	log := ctx2.Log(r)

	owner := chi.URLParam(r, "owner")
	if owner == "" {
		render2.Respond(w, http.StatusForbidden, render2.Message("Owner address is empty"))
		return
	}
	if !ctx2.Config(r).IsOwner(owner) {
		render2.Respond(w, http.StatusForbidden, render2.Message("Owner is not in the list of allowance"))
		return
	}

	source := chi.URLParam(r, "source")
	if source == "" {
		render2.Respond(w, http.StatusBadRequest, render2.Message("Not Found"))
		return
	}

	var offset, count uint64 = 0, DefaultRawsCount
	offsetRaw, countRaw := r.URL.Query().Get("offset"), r.URL.Query().Get("count")
	if offsetRaw != "" {
		var err error
		offset, err = strconv.ParseUint(offsetRaw, 10, 64)
		if err != nil {
			log.WithField("offset", offsetRaw).WithError(err).Debug("Bad offset was provided")
			render2.Respond(w, http.StatusBadRequest, render2.Message("Bad offset value"))
			return
		}
	}
	if countRaw != "" {
		var err error
		count, err = strconv.ParseUint(countRaw, 10, 64)
		if err != nil {
			log.WithField("count", countRaw).WithError(err).Debug("Bad count was provided")
			render2.Respond(w, http.StatusBadRequest, render2.Message("Bad count value"))
			return
		}
	}

	switch source {
	case storage.LocationsStorageKey:
		locations, err := ctx2.Locations(r).Select(count, offset)
		if err != nil {
			log.WithError(err).Debug("Failed to select locations")
			render2.Respond(w, http.StatusInternalServerError, render2.Message(err.Error()))
			return
		}
		render2.Respond(w, http.StatusOK, render2.Message(locations))
	case storage.DevicesStorageKey:
		devices, err := ctx2.Devices(r).Select(count, offset)
		if err != nil {
			log.WithError(err).Debug("Failed to select devices")
			render2.Respond(w, http.StatusInternalServerError, render2.Message(err.Error()))
			return
		}
		render2.Respond(w, http.StatusOK, render2.Message(devices))
	default:
		render2.Respond(w, http.StatusBadRequest, render2.Message("Source Not Found"))
	}
}
