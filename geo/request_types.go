package geo

type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Location struct {
	LatLng LatLng `json:"latLng"`
}

type Place struct {
	Location Location `json:"location"`
}

type RequestBody struct {
	Origin      Place  `json:"origin"`
	Destination Place  `json:"destination"`
	Intermediates []Place  `json:"intermediates"`
	TravelMode  string `json:"travelMode"`
}


type Route struct {
	DistanceMeters int    `json:"distanceMeters"`
	Duration       string `json:"duration"`
	Polyline       struct {
		EncodedPolyline string `json:"encodedPolyline"`
	} `json:"polyline"`
}

type Response struct {
	Routes []Route `json:"routes"`
}