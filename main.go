package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Funds struct {
	Id               int             `json:"id"`
	Name             string          `json:"name"`
	Ticker           string          `json:"ticker"`
	Amount           int64           `json:"amount"`
	PricePerItem     decimal.Decimal `json:"priceperitem"`
	PurchasePrice    decimal.Decimal `json:"purchaseprice"`
	PriceCurrent     decimal.Decimal `json:"pricecurrent"`
	PercentChanges   decimal.Decimal `json:"percentchanges"`
	YearlyInvestment decimal.Decimal `json:"yearlyinvestment"`
	ClearMoney       decimal.Decimal `json:"clearmoney"`
	DatePurchase     time.Time       `json:"datepurchase"`
	DateLastUpdate   time.Time       `json:"datelastupdate"`
	Type             string          `json:"type"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var sampleSecretKey = []byte("SecretYouShouldHide")

var user = User{
	Username: "1",
	Password: "1",
}

func main() {
	fmt.Println("My Rest Api")

	r := mux.NewRouter()
	r.Handle("/funds/usd/shares", checkAuth(getUSDFundsShares)).Methods("GET")
	r.HandleFunc("/login", login).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", r))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	// fmt.Println("user: ", u)
	checkLogin(u)
}

func checkLogin(u User) string {
	if user.Username != u.Username || user.Password != u.Password {
		fmt.Println("NOT CORRECT")
		err := "error"
		return err
	}
	validToken, err := GenerateJWT()
	fmt.Println(validToken)
	if err != nil {
		fmt.Println(err)
	}
	return validToken
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = "username"
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func getUSDFundsShares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	ArrShares := myCurrentFunds("share")
	json.NewEncoder(w).Encode(ArrShares)
}

func myCurrentFunds(fundType string) []Funds {
	var amountShared []Funds
	db, err := sql.Open("postgres", "postgres://nariman@127.0.0.1/fin?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * from fundsusd WHERE type=$1", fundType)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		f := Funds{}
		err = rows.Scan(
			&f.Id,
			&f.Name,
			&f.Ticker,
			&f.Amount,
			&f.PricePerItem,
			&f.PurchasePrice,
			&f.PriceCurrent,
			&f.PercentChanges,
			&f.YearlyInvestment,
			&f.ClearMoney,
			&f.DatePurchase,
			&f.DateLastUpdate,
			&f.Type)
		if err != nil {
			log.Fatal(err)
		}
		amountShared = append(amountShared, f)
	}

	return amountShared
}

func checkAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, nil
				}
				return sampleSecretKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "not authorized")
		}
	})
}
