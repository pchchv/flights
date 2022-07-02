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
```
/difference — Differences between the results of the two previous queries
```
