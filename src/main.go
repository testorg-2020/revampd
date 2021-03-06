package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func getParameterValues(r *http.Request) (string, string, string) {
	years := r.URL.Query().Get("operatingYear")
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")
	return years, limitString, offsetString
}

func unitRequestParser(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, int, int, int, error) {
	years, limitString, offsetString := getParameterValues(r)
	if len(years) < 1 {
		log.Debug("A valid operating year is required")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Operating year is required"))
		return w, 0, 0, 0, errors.New("No operating year")
	}

	// Check to see if the operating year is a valid int
	year, err := strconv.Atoi(years)
	if err != nil || year < 0 {
		log.Debug("Can't convert year to int.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("A valid operating year is required"))
		return w, 0, 0, 0, err
	}
	//using default if there no limit or offset
	if len(limitString) < 1 {
		log.Debug("setting limit to 100")
		limitString = "100"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil || limit < 0 {
		log.Debug("Can't convert limit to int.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid limit. Is it a positive integer?"))
		return w, 0, 0, 0, err
	}

	if len(offsetString) < 1 {
		log.Debug("setting offset to 0")
		offsetString = "0"
	}
	offset, err := strconv.Atoi(offsetString)
	if err != nil || offset < 0 {
		log.Debug("Can't convert offset to int. Is it a positive integer?")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid offset. Is it a positive integer?"))
		return w, 0, 0, 0, err
	}
	return w, year, limit, offset, nil
}

func findUnitsByOperatingYear(w http.ResponseWriter, r *http.Request, dbService *DatabaseService) {
	// parses request and checks for errors
	w, year, limit, offset, parseErr := unitRequestParser(w, r)
	if parseErr != nil {
		return
	}
	//Gets values for response
	units, err := dbService.paginatedUnitsByOperatingYear(year, limit, offset)
	total, err2 := dbService.getTotalNumberOfRows(year)
	// Check the results
	if err == nil && err2 == nil {
		if len(units) == 0 {
			log.Debug("No units found for year ", year)
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("No units were found"))
			return
		} else {
			var payload *Payload
			payload = new(Payload)
			payload.Units = units
			payload.MetaData.Retrieved = strconv.Itoa(len(units))
			payload.MetaData.Offset = strconv.Itoa(offset)
			payload.MetaData.Total = strconv.Itoa(total)
			jPayload, err := json.Marshal(payload)
			if err == nil {
				w.Write(jPayload)
			}
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Enable Cross-Origin Resource Sharing (CORS)
func enableCors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	})
}

func main() {
	dbService, err := CreateDatabaseService()

	if err != nil {
		log.Fatal("Could create database service: " + err.Error())
	}

	http.HandleFunc("/units/findByOperatingYear", func(w http.ResponseWriter, r *http.Request) {
		findUnitsByOperatingYear(w, r, dbService)
	})

	fs := http.FileServer(http.Dir("./api-spec"))
	http.Handle("/swagger/", http.StripPrefix("/swagger", fs))

	log.Debug("Starting revAMPD API backend, listening on port " + os.Getenv("PORT"))
	//Enable CORS at the server level
	http.ListenAndServe(":"+os.Getenv("PORT"), enableCors(http.DefaultServeMux))
}
