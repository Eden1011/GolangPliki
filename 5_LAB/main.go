package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	URL_STOPS         = "https://ckan.multimediagdansk.pl/dataset/c24aa637-3619-4dc2-a171-a23eec8f2172/resource/4c4025f0-01bf-41f7-a39f-d156d201b82b/download/stops.json"
	URL_BUS_DEPARTURE = "https://ckan2.multimediagdansk.pl/departures?stopId="
)

type STOP struct {
	StopId             int     `json:"stopId"`
	StopCode           string  `json:"stopCode"`
	StopName           string  `json:"stopName"`
	StopShortName      string  `json:"stopShortName"`
	StopDesc           string  `json:"stopDesc"`
	SubName            string  `json:"subName"`
	Date               string  `json:"date"`
	ZoneId             int     `json:"zoneId"`
	ZoneName           string  `json:"zoneName"`
	Virtual            int     `json:"virtual"`
	Nonpassenger       int     `json:"nonpassenger"`
	Depot              int     `json:"depot"`
	TicketZoneBorder   int     `json:"ticketZoneBorder"`
	OnDemand           int     `json:"onDemand"`
	ActivationDate     string  `json:"activationDate"`
	StopLat            float64 `json:"stopLat"`
	StopLon            float64 `json:"stopLon"`
	Type               string  `json:"type"`
	StopUrl            string  `json:"stopUrl"`
	LocationType       any     `json:"locationType"`
	ParentStation      any     `json:"parentStation"`
	StopTimezone       string  `json:"stopTimezone"`
	WheelchairBoarding any     `json:"wheelchairBoarding"`
}

type DATEDATA struct {
	LastUpdate string `json:"lastUpdate"`
	Stops      []STOP `json:"stops"`
}

type StopsResponse map[string]DATEDATA

type BUS struct {
	ID                     string `json:"id"`
	DelayInSeconds         int    `json:"delayInSeconds"`
	EstimatedTime          string `json:"estimatedTime"`
	Headsign               string `json:"headsign"`
	RouteID                int    `json:"routeId"`
	RouteShortName         string `json:"routeShortName"`
	ScheduledTripStartTime string `json:"scheduledTripStartTime"`
	TripID                 int    `json:"tripId"`
	Status                 string `json:"status"`
	TheoreticalTime        string `json:"theoreticalTime"`
	Timestamp              string `json:"timestamp"`
	Trip                   int    `json:"trip"`
	VehicleCode            int    `json:"vehicleCode"`
	VehicleID              int    `json:"vehicleId"`
	VehicleService         string `json:"vehicleService"`
}

type BusResponse struct {
	LastUpdate string `json:"lastUpdate"`
	Departures []BUS  `json:"departures"`
}

type ComparisonResult struct {
	Stop1Name       string
	Stop1ID         int
	Stop1Departures int
	Stop2Name       string
	Stop2ID         int
	Stop2Departures int
	AverageDelay1   float64
	AverageDelay2   float64
}

func getDate() string {
	now := time.Now()
	formatted := now.Format("2006-01-02")
	return formatted
}

func getStops() (*[]STOP, error) {
	resp, err := http.Get(URL_STOPS)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	var stopsResponse StopsResponse
	err = json.Unmarshal(body, &stopsResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}
	dateKey := getDate()
	dateData, exists := stopsResponse[dateKey]
	if !exists || len(stopsResponse) == 0 {
		for key := range stopsResponse {
			dateData = stopsResponse[key]
			break
		}
		if len(stopsResponse) == 0 {
			return &[]STOP{}, nil
		}
	}
	return &dateData.Stops, nil
}

func getStop(stopName string, stopArray *[]STOP) (STOP, error) {
	if stopArray == nil || len(*stopArray) == 0 {
		return STOP{}, fmt.Errorf("stop array is nil or empty")
	}
	for _, stop := range *stopArray {
		if stop.StopName == stopName {
			return stop, nil
		}
	}
	return STOP{}, fmt.Errorf("no stop found with name: %s", stopName)
}

func getBus(stopName string, stopArray *[]STOP) (*BusResponse, error) {
	stop, err := getStop(stopName, stopArray)
	if err != nil {
		return nil, fmt.Errorf("error finding stop: %w", err)
	}

	stopId := stop.StopId

	url := URL_BUS_DEPARTURE + strconv.Itoa(stopId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var busResponse BusResponse
	err = json.Unmarshal(body, &busResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &busResponse, nil
}

func getBusCompare(stopName1, stopName2 string, stopArray *[]STOP) (*ComparisonResult, error) {
	var wg sync.WaitGroup
	var busResponse1, busResponse2 *BusResponse
	var err1, err2 error
	var stop1, stop2 STOP

	stop1, err1 = getStop(stopName1, stopArray)
	if err1 != nil {
		return nil, fmt.Errorf("error finding first stop: %w", err1)
	}

	stop2, err2 = getStop(stopName2, stopArray)
	if err2 != nil {
		return nil, fmt.Errorf("error finding second stop: %w", err2)
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		busResponse1, err1 = getBus(stopName1, stopArray)
	}()

	go func() {
		defer wg.Done()
		busResponse2, err2 = getBus(stopName2, stopArray)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, fmt.Errorf("error getting bus departures for %s: %w", stopName1, err1)
	}
	if err2 != nil {
		return nil, fmt.Errorf("error getting bus departures for %s: %w", stopName2, err2)
	}

	var totalDelay1, totalDelay2 int
	for _, bus := range busResponse1.Departures {
		totalDelay1 += bus.DelayInSeconds
	}
	for _, bus := range busResponse2.Departures {
		totalDelay2 += bus.DelayInSeconds
	}

	var avgDelay1, avgDelay2 float64
	if len(busResponse1.Departures) > 0 {
		avgDelay1 = float64(totalDelay1) / float64(len(busResponse1.Departures))
	}
	if len(busResponse2.Departures) > 0 {
		avgDelay2 = float64(totalDelay2) / float64(len(busResponse2.Departures))
	}

	result := &ComparisonResult{
		Stop1Name:       stopName1,
		Stop1ID:         stop1.StopId,
		Stop1Departures: len(busResponse1.Departures),
		Stop2Name:       stopName2,
		Stop2ID:         stop2.StopId,
		Stop2Departures: len(busResponse2.Departures),
		AverageDelay1:   avgDelay1,
		AverageDelay2:   avgDelay2,
	}

	return result, nil
}

func main() {
	stops, err := getStops()
	if err != nil {
		fmt.Println("Error getting stops:", err)
		return
	}

	fmt.Printf("Retrieved %d stops\n", len(*stops))

	var validStops []STOP
	for _, stop := range *stops {
		if stop.StopName != "" {
			validStops = append(validStops, stop)
			if len(validStops) >= 5 {
				break
			}
		}
	}

	fmt.Println("\nValid stops with non-null names:")
	for i, stop := range validStops {
		fmt.Printf("%d. Stop %d: %s (%s) - %s\n",
			i+1, stop.StopId, stop.StopName, stop.StopCode, stop.Type)
	}

	if len(validStops) >= 2 {
		fmt.Println("\nComparing bus departures for two stops...")

		stop1Name := validStops[0].StopName
		stop2Name := validStops[1].StopName

		var wg sync.WaitGroup
		wg.Add(1)

		resultChan := make(chan *ComparisonResult, 1)
		errChan := make(chan error, 1)

		go func() {
			defer wg.Done()

			result, err := getBusCompare(stop1Name, stop2Name, stops)
			if err != nil {
				errChan <- err
				return
			}

			resultChan <- result
		}()

		stopName := "Kameliowa"
		busResponse, err := getBus(stopName, stops)
		if err != nil {
			fmt.Println("\nError getting bus departures for", stopName, ":", err)
		} else {
			stop, _ := getStop(stopName, stops)

			fmt.Printf("\nBus departures for %s (ID: %d) (Last updated: %s):\n",
				stopName, stop.StopId, busResponse.LastUpdate)

			if len(busResponse.Departures) == 0 {
				fmt.Println("No departures found.")
			} else {
				for i, bus := range busResponse.Departures {
					fmt.Printf("%d. Route %s to %s - Estimated time: %s (Delay: %d seconds)\n",
						i+1, bus.RouteShortName, bus.Headsign, bus.EstimatedTime, bus.DelayInSeconds)
				}
			}
		}

		wg.Wait()

		select {
		case result := <-resultChan:
			fmt.Printf("\n--- Bus Comparison Results ---\n")
			fmt.Printf("Stop 1: %s (ID: %d)\n", result.Stop1Name, result.Stop1ID)
			fmt.Printf("  - Number of departures: %d\n", result.Stop1Departures)
			fmt.Printf("  - Average delay: %.2f seconds\n", result.AverageDelay1)
			fmt.Printf("\nStop 2: %s (ID: %d)\n", result.Stop2Name, result.Stop2ID)
			fmt.Printf("  - Number of departures: %d\n", result.Stop2Departures)
			fmt.Printf("  - Average delay: %.2f seconds\n", result.AverageDelay2)

			if result.Stop1Departures > result.Stop2Departures {
				fmt.Printf("\nStop %s has more departures than %s\n", result.Stop1Name, result.Stop2Name)
			} else if result.Stop1Departures < result.Stop2Departures {
				fmt.Printf("\nStop %s has more departures than %s\n", result.Stop2Name, result.Stop1Name)
			} else {
				fmt.Printf("\nBoth stops have the same number of departures\n")
			}

			if result.AverageDelay1 < result.AverageDelay2 {
				fmt.Printf("Buses at %s are more punctual (less delay) than at %s\n", result.Stop1Name, result.Stop2Name)
			} else if result.AverageDelay1 > result.AverageDelay2 {
				fmt.Printf("Buses at %s are more punctual (less delay) than at %s\n", result.Stop2Name, result.Stop1Name)
			} else {
				fmt.Printf("Buses at both stops have similar punctuality\n")
			}

		case err := <-errChan:
			fmt.Printf("\nError comparing stops: %s\n", err)

		default:
			fmt.Println("\nNo comparison results available")
		}
	}
}

