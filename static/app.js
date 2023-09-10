let distanceMatrixService;
let map;
let originMarker;
let infowindow;
let userLocation;
let poly;
let circles = [];
let waypoints = [];

const distanceInput = document.getElementById("distance");
const generatePathBtn = document.getElementById("generatePath");

generatePathBtn.addEventListener("click", readUserDistanceAndTryGenerating);
distanceInput.addEventListener("keypress", (event) => {
  if (event.key != "Enter") {
    return;
  }
  readUserDistanceAndTryGenerating();
});

function readUserDistanceAndTryGenerating() {
  const pathDistance = distanceInput.value;
  if (!pathDistance) {
    distanceInput.classList.add("error");
  } else {
    distanceInput.classList.remove("error");
  }

  doGeneratePath(pathDistance);
}

const WARSAW = { lat: 52.2297, lng: 21.0122 };

// callback to show position
// loc = navigator.geolocation.getCurrentPosition(showPosition);

function refillWaypoints(resp) {
  clearWaypoints();
  for (let i = 0; i < resp.Points.length; ++i) {
    waypoints.push(pointToCircle(resp.Points[i]));
  }
}

function clearWaypoints() {
  for (let i = 0; i < waypoints.length; ++i) {
    waypoints[i].setMap(null);
  }
  waypoints = [];
}

function showPosition(position) {
  // const lat = position.coords.latitude;
  // const lng = position.coords.longitude;
  lat = position.lat;
  lng = position.lng;

  userLocation = { lat, lng };
  map.setCenter(userLocation);
  pointToCircle(userLocation);
  // var element = document.getElementById("map");
  // element.removeAttribute("hidden");
}

const doGeneratePath = async (distance) => {
  const url = "/route";
  payload = {
    lat: WARSAW.lat,
    lng: WARSAW.lng,
    distance: distance,
  };
  const queryParams = new URLSearchParams(payload).toString();
  const fullURL = `${url}?${queryParams}`;
  const response = await fetch(fullURL);
  const el = await response.json();
  console.log(el);
  decodedPath = google.maps.geometry.encoding.decodePath(el.polyLine);
  // decodedLvls = google.maps.geometry.encoding.decodeLevels("BBBB");
  // decodedLevels = goog;
  if (poly) {
    poly.setMap(null);
  }
  refillWaypoints(el);
  poly = new google.maps.Polyline({
    path: decodedPath,
    geodesic: true,
    strokeColor: "#FF0000",
    strokeOpacity: 1.0,
    strokeWeight: 5,
  });
  poly.setMap(map);
};

async function initialize() {
  initMap();
  showPosition(WARSAW);
  doGeneratePath(5);
}

const initMap = () => {
  // TODO: Start Distance Matrix service

  // The map, centered on Austin, TX
  map = new google.maps.Map(document.querySelector("#map"), {
    center: WARSAW,
    // zoom: 14,
    zoom: 12,
    // mapId: 'YOUR_MAP_ID_HERE',
    clickableIcons: false,
    fullscreenControl: false,
    mapTypeControl: false,
    rotateControl: true,
    scaleControl: false,
    streetViewControl: true,
    zoomControl: true,
  });
};

const fetchAndRenderStores = async (center) => {
  // Fetch the stores from the data source
  stores = (await fetchStores(center)).features;

  // Create circular markers based on the stores
  circles = stores.map((store) => storeToCircle(store, map));
};

const fetchStores = async (center) => {
  const url = `/data/dropoffs`;
  const response = await fetch(url);
  return response.json();
};

const image =
  "https://developers.google.com/maps/documentation/javascript/examples/full/images/beachflag.png";
const pointToCircle = (pt) => {
  // const [lng, lat] = store.geometry.coordinates;
  const circle = new google.maps.Marker({
    position: pt,
    map,
    icon: image,
  });

  return circle;
};
