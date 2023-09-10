package geo

import (
	"fmt"
	"math"
	"net/url"
	"testing"
)

// Query params parser test
func TestQueryParamFailingSplit(t *testing.T) {
	tests := map[string]struct {
		distanceStr string
		latStr string
		lngStr string
	} {
		"wrongDist": {distanceStr: "x", latStr: "0", lngStr: "0"},
		"wrongLat": {distanceStr: "0", latStr: "x", lngStr: "0"},
		"wrongLng": {distanceStr: "0", latStr: "0", lngStr: "x"},
		"wrongAll": {distanceStr: "x", latStr: "x", lngStr: "x"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T){
			vals := url.Values{}
			vals.Set("distance", tc.distanceStr)
			vals.Set("lat", tc.latStr)
			vals.Set("lng", tc.lngStr)
			params, err := ParseRouteQueryParams(vals)
			if err == nil {
				t.Errorf("should fail, not success: %f", params)
			}
		})
	}
}

func TestQueryParamsOkay(t *testing.T) {
	vals := url.Values{}
	vals.Set("distance", "0")
	vals.Set("lat", "0")
	vals.Set("lng", "0")
	_, err := ParseRouteQueryParams(vals)
	if err != nil {
		t.Errorf("should succeed, got error: %s", err.Error())
	}
}

func TestRejectEmptyParams(t *testing.T) {
	var vals url.Values
	// vals.G
	params, err := ParseRouteQueryParams(vals)
	
	if params != (QueryParams{}) {
		t.Error()
	}

	if err == nil {
		t.Error("err should be empty")
	}
}



// Distance checker test
func TestDistanceTooLargeFails(t *testing.T) {
	err := CheckParams(QueryParams{
		distance: 101,
	})
	if err == nil {
		t.Errorf("should reject big distance")
	}
}

func TestDistanceNegativeFails(t *testing.T) {
	err := CheckParams(QueryParams{
		distance: -1,
	})
	if err == nil {
		t.Errorf("should reject negative distance")
	}
}
func TestDistanceOk(t *testing.T) {
	err := CheckParams(QueryParams{
		distance: 50,
	})
	if err != nil {
		t.Errorf("should have no error, got %s", err.Error())
	}
}

const eps = 1e-9

func TestGeneratePointOnCircle(t  *testing.T) {
	warsawPoint := Point{52.2297, 21.0122}
	distKm := 5.0
	res := RandomPointOnCircleWithRadius(warsawPoint, distKm)
	calculatedDistKm := haversine(warsawPoint, res)
	if math.Abs(calculatedDistKm - distKm) > eps {
		t.Errorf("expected distance %f, got %f", distKm, calculatedDistKm)
	}
}

func TestPointsSumDistanceOkay(t  *testing.T) {
	warsawPoint := Point{52.2297, 21.0122}
	distKm := 6.0
	p2, p3 := GenerateTrianglePoints(warsawPoint, distKm)
	dist1 := haversine(warsawPoint, p2)
	fmt.Println("dist1: ", dist1)
	dist2 := haversine(p2, p3)
	fmt.Println("dist2: ", dist2)
	dist3 := haversine(p3, warsawPoint)
	fmt.Println("dist3: ", dist3)
	totalDistance := dist1 + dist2 + dist3
	if math.Abs(totalDistance - distKm) > eps {
		t.Errorf("expected distance %f, got %f", distKm, totalDistance)
	}
}

// Convert degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Haversine function to compute distance between two latitude-longitude points
func haversine(p1, p2 Point) float64 {
	deltaLat := toRadians(p2.Lat - p1.Lat)
	deltaLng := toRadians(p2.Lng - p1.Lng)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(toRadians(p1.Lat))*math.Cos(toRadians(p2.Lat))*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}
