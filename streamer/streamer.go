package streamer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Charts data representation.
type datapoint struct {
	PanelID int                    `json:"panelid"` // Grafana panel id.
	RefID   string                 `json:"refid"`   // Grafana panel query ref id.
	Values  map[string]interface{} `json:"values"`  // Values associated with the row text. Timestamp is hardcoded every time.
}

type QueryConfig struct {
	Writer http.ResponseWriter
	Series []string
	Ticker *time.Ticker
}

var configMap = make(map[int]map[string]QueryConfig)
var Origin = ""

func sendData(pID int, rID string, fillDataRow func(string, map[string]interface{})) {
	dataRow := make(map[string]interface{})
	dataRow["timestamp"] = time.Now().UnixNano() / 1000000

	for _, row := range configMap[pID][rID].Series {
		fillDataRow(row, dataRow)
	}

	currentPoint := &datapoint{
		RefID:   rID,
		PanelID: pID,
		Values:  dataRow,
	}
	j, _ := json.Marshal(currentPoint)
	configMap[pID][rID].Writer.Header().Set("Access-Control-Allow-Origin", Origin)
	fmt.Fprintf(configMap[pID][rID].Writer, "%s\n", j)
	configMap[pID][rID].Writer.(http.Flusher).Flush() // Trigger "chunked" encoding and s
}

func Configure(pID int, rID string, rowsParam string, startTime int, endTime int, dataPointCount int, w http.ResponseWriter) {
	// Sets the sampling time based on the allowed grafana data points
	samplingTimeMs := (endTime - startTime) / dataPointCount
	if _, ok := configMap[pID]; !ok {
		configMap[pID] = make(map[string]QueryConfig)
	}

	// Creates the config map
	configMap[pID][rID] = QueryConfig{
		Series: strings.Split(rowsParam, ","),
		Ticker: time.NewTicker(time.Duration(samplingTimeMs) * time.Millisecond),
		Writer: w,
	}
}

// StreamData starts data streaming using the respective ticker stored in configMap.
// This function shall be called from a REST request handler function
// Based on the parsed parameters of the request the configMap is extended.
// Structure of config map:
//			Panel -> queries/panel -> series/queries
// At the moment each query has its own ticker channel setup based on the sampling time (calculated based on the grafana data points config)
// Each series is represented by a string that can be configured in grafana through the dataText field.
// Once the ticker is set, the respective query results will be streamed back (periodical http flush) to the grafana server.
func StreamData(pID int, rID string, fillDataRow func(string, map[string]interface{})) {
	log.Printf("Start streaming for Panel %d and query %s", pID, rID)
	defer configMap[pID][rID].Ticker.Stop()
	for ; true; <-configMap[pID][rID].Ticker.C {
		sendData(pID, rID, fillDataRow)
	}
}
