# Card Api in Go Lang

> Simple RESTful API to create deck, open deck, and draw cards

##  Install


``` bash
# Install mux router
go get -u github.com/gorilla/mux
```

``` bash
go build
./go-card_rest
```

## Endpoints

### Get All Deck List
``` bash
GET http://localhost:8000/decks
```
### Create New Deck
``` bash
POST http://localhost:8000/create_deck
Body : {
    "shuffled":false,
    "cards":"KH,QH,4S"
}
```

### Open Deck
``` bash
GET http://localhost:8000/open_deck/{uuid}
```

### Draw Card
``` bash
GET http://localhost:8000/draw_card/{uuid}?count=2

PARAMS : count = 1

```

