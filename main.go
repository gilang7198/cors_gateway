package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
)

type Response struct {
	GeocodedWaypoints []GeocodedWaypoint `json:"geocoded_waypoints"`
	Routes            []Route            `json:"routes"`            
	Status            string             `json:"status"`            
}

type GeocodedWaypoint struct {
	GeocoderStatus string   `json:"geocoder_status"`
	PlaceID        string   `json:"place_id"`       
	Types          []string `json:"types"`          
}

type Route struct {
	Bounds           Bounds        `json:"bounds"`           
	Copyrights       string        `json:"copyrights"`       
	Legs             []Leg         `json:"legs"`             
	OverviewPolyline Polyline      `json:"overview_polyline"`
	Summary          string        `json:"summary"`          
	Warnings         []interface{} `json:"warnings"`         
	WaypointOrder    []interface{} `json:"waypoint_order"`   
}

type Bounds struct {
	Northeast Northeast `json:"northeast"`
	Southwest Northeast `json:"southwest"`
}

type Northeast struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Leg struct {
	Distance          Distance      `json:"distance"`           
	Duration          Distance      `json:"duration"`           
	EndAddress        string        `json:"end_address"`        
	EndLocation       Northeast     `json:"end_location"`       
	StartAddress      string        `json:"start_address"`      
	StartLocation     Northeast     `json:"start_location"`     
	Steps             []Step        `json:"steps"`              
	TrafficSpeedEntry []interface{} `json:"traffic_speed_entry"`
	ViaWaypoint       []interface{} `json:"via_waypoint"`       
}

type Distance struct {
	Text  string `json:"text"` 
	Value int64  `json:"value"`
}

type Step struct {
	Distance         Distance   `json:"distance"`          
	Duration         Distance   `json:"duration"`          
	EndLocation      Northeast  `json:"end_location"`      
	HTMLInstructions string     `json:"html_instructions"` 
	Polyline         Polyline   `json:"polyline"`          
	StartLocation    Northeast  `json:"start_location"`    
	TravelMode       TravelMode `json:"travel_mode"`       
	Maneuver         *string    `json:"maneuver,omitempty"`
}

type Polyline struct {
	Points string `json:"points"`
}

type TravelMode string
const (
	Driving TravelMode = "DRIVING"
)

func homePage(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    key := vars["url"]

    response, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?" + key)
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var responseObject Response
    json.Unmarshal(responseData, &responseObject)
	json.NewEncoder(w).Encode(responseObject)
    // fmt.Fprintf(w, "https://maps.googleapis.com/maps/api/directions/json?" + key)
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/{url}", homePage)
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}