package handlers

import (
	"github.com/geo-provider/app/ctx"
	"github.com/geo-provider/app/render"
	"github.com/geo-provider/app/storage"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

const DefaultRawsCount = 100

func GetSources(w http.ResponseWriter, r *http.Request) {
	storagesList := ctx.Config(r).ListSources()

	render.Respond(w, http.StatusOK, render.Message(storagesList))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	log := ctx.Log(r)

	owner := chi.URLParam(r, "owner")
	if owner == "" {
		render.Respond(w, http.StatusForbidden, render.Message("Owner address is empty"))
		return
	}
	if !ctx.Config(r).IsOwner(owner) {
		render.Respond(w, http.StatusForbidden, render.Message("Owner is not in the list of allowance"))
		return
	}

	source := chi.URLParam(r, "source")
	if source == "" {
		render.Respond(w, http.StatusBadRequest, render.Message("Not Found"))
		return
	}

	var offset, count uint64 = 0, DefaultRawsCount
	offsetRaw, countRaw := r.URL.Query().Get("offset"), r.URL.Query().Get("count")
	if offsetRaw != "" {
		var err error
		offset, err = strconv.ParseUint(offsetRaw, 10, 64)
		if err != nil {
			log.WithField("offset", offsetRaw).WithError(err).Debug("Bad offset was provided")
			render.Respond(w, http.StatusBadRequest, render.Message("Bad offset value"))
			return
		}
	}
	if countRaw != "" {
		var err error
		count, err = strconv.ParseUint(countRaw, 10, 64)
		if err != nil {
			log.WithField("count", countRaw).WithError(err).Debug("Bad count was provided")
			render.Respond(w, http.StatusBadRequest, render.Message("Bad count value"))
			return
		}
	}

	switch source {
	case storage.LocationsStorageKey:
		locations, err := ctx.Locations(r).Select(count, offset)
		if err != nil {
			log.WithError(err).Debug("Failed to select locations")
			render.Respond(w, http.StatusInternalServerError, render.Message(err.Error()))
			return
		}
		render.Respond(w, http.StatusOK, render.Message(locations))
	case storage.DevicesStorageKey:
		devices, err := ctx.Devices(r).Select(count, offset)
		if err == nil {
			log.WithError(err).Debug("Failed to select devices")
			render.Respond(w, http.StatusInternalServerError, render.Message(err.Error()))
			return
		}
		render.Respond(w, http.StatusOK, render.Message(devices))
	default:
		render.Respond(w, http.StatusBadRequest, render.Message("Source Not Found"))
	}
}
