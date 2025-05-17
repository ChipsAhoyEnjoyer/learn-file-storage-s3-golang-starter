package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (cfg *apiConfig) handlerResetVideos(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	assets, err := os.ReadDir(cfg.assetsRoot)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset assets folder", err)
		return
	}
	for _, f := range assets {
		err = os.Remove(filepath.Join(cfg.assetsRoot, f.Name()))
		if err != nil {
			log.Printf("Error resetting assets folder: \n%v\n", err)
		}
	}

	if err = cfg.db.ResetVideos(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset videos", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Videos table reset to initial state"))
}
