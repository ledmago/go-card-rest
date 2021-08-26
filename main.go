package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)


func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}



// Deck struct (Model)
type Deck struct {
	DeckId string  `json:"deck_id"`
	Shuffled   bool  `json:"shuffled"`
	Remaining int  `json:"remaining"`
	Cards []Card `json:"cards"`
}

type CarList struct {
	Cards []Card  `json:"cards"`

}
// Cards struct

var decks []Deck


var default_cards []Card


type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}



type Message struct {
	Shuffled bool `json:"shuffled"`
	Cards string `json:"cards"`
}
type ErrorMessage struct {
	Message string `json:"message"`
}


// Get all Decks
func getDecks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(decks)
}

func createCard(Shuffled bool, Cards string)[]Card {

	cards_input:=  strings.Split(strings.Replace(Cards," ", "",-1), ",")

	var fullCard bool
	if len(Cards) == 0{
		fullCard = true
	}else {
		fullCard = false
	}

	var cards []Card;

	for d := 0; d < 4; d++ {
		for card := 0; card < 13; card++ {
			var suit string
			var code string
			switch d {
			case 0: {
				suit = "SPADES"
				code = "S"
			}
			case 1:{
					suit = "DIAMONDS"
					code = "D"
				}
			case 2:{
					suit = "CLUBS"
					code = "C"
				}
			case 3:
				{
					suit = "HEARTS"
					code = "H"
				}
			}
			var value string
			switch card {
			case 0:{
				value = "ACE"
				code = "A" + code
			}
			case 10:{
				value ="JACK"
				code = "J" + code
			}
			case 11:{
				value ="QUEEN"
				code = "Q" + code
			}
			case 12:{
				value ="KING"
				code = "K" + code
			}
			default:{
				value = strconv.Itoa(card +1)
				code = strconv.Itoa(card + 1) + code
			}
			}


			if !fullCard{
				if contains(cards_input, code){
					cards = append(cards,Card {value,suit,code} )
				}
			}else {
				cards = append(cards,Card {value,suit,code} )
			}


		}
	}
	if(Shuffled){
		Shuffle(cards)
	}

	return cards

}


func Shuffle(slc []Card) {
	for i := 1; i < len(slc); i++ {
		r := rand.Intn(i + 1)
		if i != r {
			slc[r], slc[i] = slc[i], slc[r]
		}
	}
}


func openDeck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, item := range decks {
		if item.DeckId == params["deck_id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	var errorMessage ErrorMessage = ErrorMessage{"The deck is not found by id"}
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(errorMessage)
}

func remove(slice []Card, s int) []Card {
	return append(slice[:s], slice[s+1:]...)
}


func removeCard(DeckId string, Code string){
	for itemIndex, item := range decks {
		if item.DeckId == DeckId {

			for index, item2 := range item.Cards{
				if item2.Code == Code{
					decks[itemIndex].Cards = remove(decks[itemIndex].Cards,index)
					decks[itemIndex].Remaining = len(decks[itemIndex].Cards)
					return
				}
			}
		}
	}
}

func drawDeck(w http.ResponseWriter, r *http.Request) {



	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, item := range decks {
		if item.DeckId == params["deck_id"] {

			count, _ := strconv.Atoi(r.URL.Query().Get("count"))
			var cardList CarList
			cardList.Cards = []Card{}

			for i := 0; i < count && i < len(item.Cards); i++ {
				card := item.Cards[i]
			 	cardList.Cards = append(cardList.Cards,item.Cards[i])
			 	removeCard(item.DeckId,card.Code)
			}

			json.NewEncoder(w).Encode(cardList)
			return
		}
	}
	var errorMessage ErrorMessage = ErrorMessage{"The deck is not found by id"}
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(errorMessage)
}

func generateId() string{
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
// Add new deck
func createDeck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var deck Deck
	var message Message
	_ = json.NewDecoder(r.Body).Decode(&message)

	deck.DeckId = generateId() // Mock ID - not safe
	deck.Shuffled = message.Shuffled
	var cards []*Card;
	cards = append(cards,&Card {"asd","asd","code"} )

	deck.Cards = createCard(message.Shuffled, message.Cards)
	deck.Remaining = len(deck.Cards)
	decks = append(decks,deck)
	json.NewEncoder(w).Encode(decks)





}

// Main function
func main() {
	// Init router
	r := mux.NewRouter()



	// Route handles & endpoints
	r.HandleFunc("/decks", getDecks).Methods("GET")
	r.HandleFunc("/create_deck", createDeck).Methods("POST")
	r.HandleFunc("/open_deck/{deck_id}", openDeck).Methods("GET")
	r.HandleFunc("/draw_card/{deck_id}", drawDeck).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
