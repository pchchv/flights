# Running the application
```
docker-compose up --build
```
# Running the application without Docker
```
go run .
```
## Running tests (app must be running)
```
go test
```
## HTTP Methods
```
/ping — Checking the server connection
```
```
/flights — All flight options from DXB to BKK
```
```
/variants — The most expensive/cheapest, fastest/longest and best flight variants
options: 
    duration — Fastest/longest options
    optimal — Best option
    price — Most expensive/cheapest options
    
example: http://localhost:8080/variants?optinons=duration,optimal,price
```
