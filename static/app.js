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

loc = navigator.geolocation.getCurrentPosition(showPosition);

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
  lat = position.coords.latitude;
  lng = position.coords.longitude;
  document.getElementById("generatePath").textContent = "Generate!";

  userLocation = { lat, lng };
  map.setCenter(userLocation);
  pointToCircle(userLocation);
}

const doGeneratePath = async (distance) => {
  if (!lat || !lng) {
    console.log("unset!");
    return;
  }
  const url = "/route";
  payload = {
    lat: lat,
    lng: lng,
    distance: distance,
  };
  const queryParams = new URLSearchParams(payload).toString();
  const fullURL = `${url}?${queryParams}`;
  const response = await fetch(fullURL);
  const el = await response.json();
  console.log(el);
  decodedPath = google.maps.geometry.encoding.decodePath(el.polyLine);
  document.getElementById("outputBox").textContent =
    "Actual path distance (m): " + el.distance;
  document.getElementById("outputBox").removeAttribute("hidden");
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
}

const initMap = () => {
  map = new google.maps.Map(document.querySelector("#map"), {
    center: WARSAW,
    zoom: 12,
    clickableIcons: false,
    fullscreenControl: false,
    mapTypeControl: false,
    rotateControl: true,
    scaleControl: false,
    streetViewControl: true,
    zoomControl: true,
  });
};

const fetchStores = async (center) => {
  const url = `/data/dropoffs`;
  const response = await fetch(url);
  return response.json();
};

const image =
  "https://developers.google.com/maps/documentation/javascript/examples/full/images/beachflag.png";
const pointToCircle = (pt) => {
  const circle = new google.maps.Marker({
    position: pt,
    map,
    icon: null,
  });

  return circle;
};
