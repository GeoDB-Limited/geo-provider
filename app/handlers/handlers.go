package handlers

import (
	"github.com/geo-provider/app/ctx"
	"github.com/geo-provider/app/render"
	"github.com/geo-provider/storage"
	"github.com/geo-provider/utils"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

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

	var srcPath string
	if srcPath = ctx.Config(r).Source(source); srcPath == "" {
		log.WithField("source", source).Debug("Provided source not found")
		render.Respond(w, http.StatusBadRequest, render.Message("Source Not Found"))
		return
	}

	offset, count := 0, storage.MaxReadRows

	offsetRaw, countRaw := r.URL.Query().Get("offset"), r.URL.Query().Get("count")

	if offsetRaw != "" {
		var err error
		offset, err = strconv.Atoi(offsetRaw)
		if err != nil {
			log.WithField("offset", offsetRaw).WithError(err).Debug("Bad offset was provided")
			render.Respond(w, http.StatusBadRequest, render.Message("Bad offset value"))
			return
		}
	}

	if countRaw != "" {
		var err error
		count, err = strconv.Atoi(countRaw)
		if err != nil {
			log.WithField("count", countRaw).WithError(err).Debug("Bad count was provided")
			render.Respond(w, http.StatusBadRequest, render.Message("Bad count value"))
			return
		}
	}

	keeper, err := storage.Open(srcPath)
	if err != nil {
		log.WithField("source_path", srcPath).WithError(err).Debug("storage path not found")
		render.Respond(w, http.StatusInternalServerError, render.Message("Couldn't open storage"))
		return
	}

	resp, err := keeper.Read(offset, utils.Min(storage.MaxReadRows, count))
	if err != nil {
		log.WithError(err).Debug("failed to read source")
		render.Respond(w, http.StatusInternalServerError, render.Message("Bad count value"))
		return
	}

	render.Respond(w, http.StatusOK, render.Message(resp))
}
