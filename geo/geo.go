package geo

import (
	"bytes"
	"encoding/json"

	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

// Earth radius in KM
const earthRadius = 6371.0
// const url = "https://routes.googleapis.com/directions/v2:computeRoutes"

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func RandomPointOnCircleWithRadius(point Point, distance float64) Point {
	// Convert latitude and longitude from degrees to radians
	latRad := point.Lat * (math.Pi / 180.0)
	lngRad := point.Lng * (math.Pi / 180.0)

	// Random angle
	randomAngle := rand.Float64() * 2 * math.Pi

	// Convert distance to angular distance in radians
	angularDistance := distance / earthRadius

	// Calculate new latitude
	newLatRad := math.Asin(math.Sin(latRad)*math.Cos(angularDistance) +
		math.Cos(latRad)*math.Sin(angularDistance)*math.Cos(randomAngle))

	// Calculate new longitude
	newLngRad := lngRad + math.Atan2(math.Sin(randomAngle)*math.Sin(angularDistance)*math.Cos(latRad),
		math.Cos(angularDistance)-math.Sin(latRad)*math.Sin(newLatRad))

	// Convert new latitude and longitude from radians back to degrees
	newLat := newLatRad * (180.0 / math.Pi)
	newLng := newLngRad * (180.0 / math.Pi)

	return Point{newLat, newLng}
}

func rotateAroundPoint(dir Point, angle float64) Point {
	// Rotate point around the origin (0, 0) using a rotation matrix
	newX := dir.Lat*math.Cos(angle) - dir.Lng*math.Sin(angle)
	newY := dir.Lat*math.Sin(angle) + dir.Lng*math.Cos(angle)
	return Point{newX, newY}
}

func GenerateTrianglePoints(a Point, d float64) (Point, Point) {
	b := RandomPointOnCircleWithRadius(a, d/3)

	// Get the direction vector from a to b
	dir := Point{b.Lat - a.Lat, b.Lng - a.Lng}

	// Rotate the direction vector by 60 degrees to get the direction from b to c
	cDir := rotateAroundPoint(dir, math.Pi/3)

	c := Point{b.Lat + cDir.Lat, b.Lng + cDir.Lng}

	return b, c
}

type QueryParams struct {
	point Point
	distance float64
}

type RouteResponse struct {
	PolyLine string `json:"polyLine"`
	Points []Point `json:"points`
}



func ParseRouteQueryParams(params url.Values) (QueryParams, error) {
	distance, err := strconv.ParseFloat(params.Get("distance"), 64); 
	if err != nil {
		return QueryParams{}, fmt.Errorf("cannot parse distance: %s", params.Get("distance"))
	}

	lat, err := strconv.ParseFloat(params.Get("lat"), 64); 
	if err != nil {
		return QueryParams{}, errors.New("cannot parse lat")
	}

	lng, err := strconv.ParseFloat(params.Get("lng"), 64); 
	if err != nil {
		return QueryParams{}, errors.New("cannot parse lng")
	}

	return QueryParams{Point{lat, lng}, distance}, nil
}

func CheckParams(params QueryParams) error {
	if params.distance > 100 || params.distance < 0 {
		return errors.New("distance too large")
	}
	return nil
}

func GetRouteData(params QueryParams) (RouteResponse, error) {
	p1 := params.point
	// p2 := RandomPointOnCircleWithRadius(p1, params.distance)
	p2, p3 := GenerateTrianglePoints(p1, params.distance / 2)
	// fmt.Println("first point: %f, %f, generated point: %f, %f", p1.Lat, p1.Lng, p2.Lat, p2.Lng)
	poly, err := getPolyLine(p1, p2, p3, p1)	
	if err != nil {
		return RouteResponse{}, fmt.Errorf("cannot get route data: %s", err.Error())
	}
	return RouteResponse{poly, []Point{p1, p2, p3}}, nil
}

func PlaceFromPoint(p Point) Place {
	return Place{
		Location: Location{
			LatLng: LatLng{
				Latitude:  p.Lat,
				Longitude: p.Lng,
			},
		},
	}
}

func GetRequestBody(p []Point) ([]byte, error) {
	origin := p[0]
	destination := p[len(p)-1]
	body := RequestBody{
		Origin: PlaceFromPoint(origin),
		Destination: PlaceFromPoint(destination),
		TravelMode: "WALK",
	}
	for i := 1; i < len(p) - 1; i += 1 {
		fmt.Println("adding intermediate: ", PlaceFromPoint(p[i]))
		body.Intermediates = append(body.Intermediates, PlaceFromPoint(p[i]))
	}
	fmt.Println("num intermediates: ", len(body.Intermediates))
	
	return json.Marshal(body)
}

func getPolyLine(points ...Point) (string, error) {
	if len(points) < 2 {
		return "", errors.New("too few points to create a route")
	}

	
	url := "https://routes.googleapis.com/directions/v2:computeRoutes"
	body, err := GetRequestBody(points)
	if err != nil {
		return "", fmt.Errorf("cannot get req body: %s", err.Error())
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("cannot create url: ", err.Error())
		return "", fmt.Errorf("cannot create url: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", "AIzaSyCeCRjWtqsHZGds5QmZzxuro7oVbXRlMTQ")
	req.Header.Set("X-Goog-FieldMask", "routes.duration,routes.distanceMeters,routes.polyline.encodedPolyline")


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot do url request: %s", err.Error())
	}
	defer resp.Body.Close()

	var rspBody Response
	json.NewDecoder(resp.Body).Decode(&rspBody)
	fmt.Printf("Resp: %v", rspBody)
	return rspBody.Routes[0].Polyline.EncodedPolyline, nil
}
