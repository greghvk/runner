# runner

This app allows you to create a walking/running path with a distance of selected length. The path will start and end in your current location.

To use this app, you need google a Google Maps API key (with routes API enabled). When you obtain the key, replace `API_KEY` string in the maps link in `index.html` file with the key, and set `ROUTES_KEY` environment variable as the key (`export ROUTES_KEY=key`).

Then, start the app with `go run main.go`, go to `localhost:8080` and start generating your paths! (you need to allow location for this app to work) 
