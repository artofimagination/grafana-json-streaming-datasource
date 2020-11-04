package main

import (
	"crypto/rand"
	"html/template"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"grafana-json-streaming-datasource/streamer"
)

func main() {
	// Grafana sends a request to initiate streaming in the form of http://localhost:8080?panelid=5&refid=A&data-rows=test1,test2
	http.HandleFunc("/", streamHandler)
	http.HandleFunc("/show", showChart)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func showChart(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t := template.Must(template.ParseFiles(wd + "/chart.html"))

	empty := 0
	err = t.ExecuteTemplate(w, "chart.html", empty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	streamer.Origin = r.Header.Get("Origin")
	// Parses the parameters.
	panelParam := r.URL.Query().Get("panelid")
	if panelParam == "" {
		panic("No panel id")
	}
	panelid, _ := strconv.Atoi(panelParam)

	startTimeParam := r.URL.Query().Get("start")
	if startTimeParam == "" {
		panic("No start time")
	}
	startTime, _ := strconv.Atoi(startTimeParam)

	endTimeParam := r.URL.Query().Get("end")
	if endTimeParam == "" {
		panic("No end time")
	}
	endTime, _ := strconv.Atoi(endTimeParam)

	datapointsParam := r.URL.Query().Get("datapoints")
	if datapointsParam == "" {
		panic("No datapoints")
	}
	datapoints, _ := strconv.Atoi(datapointsParam)

	rowsParam := r.URL.Query().Get("data-rows")
	if rowsParam == "" {
		panic("No rows")
	}

	refidParam := r.URL.Query().Get("refid")
	if refidParam == "" {
		panic("No ref id")
	}

	fillDataRows := func(row string, dataRow map[string]interface{}) {
		switch row {
		case "cpu_load":
			random, _ := rand.Int(rand.Reader, big.NewInt(10))
			dataRow[row] = random.Int64()
		case "available_memory":
			random, _ := rand.Int(rand.Reader, big.NewInt(10))
			dataRow[row] = random.Int64() + 10
		case "row_count":
			random, _ := rand.Int(rand.Reader, big.NewInt(10))
			dataRow[row] = random.Int64() + 20
		}
	}

	streamer.Configure(panelid, refidParam, rowsParam, startTime, endTime, datapoints, w)
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	streamer.StreamData(panelid, refidParam, fillDataRows)
}
