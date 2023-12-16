package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nayeemnishaat/go-web-app/api/lib"
	"github.com/nayeemnishaat/go-web-app/api/model"
	"github.com/nayeemnishaat/go-web-app/web/template"
)

func (app *Application) TerminalPage(w http.ResponseWriter, r *http.Request) {
	if err := app.RenderTemplate(w, r, "terminal", &template.TemplateData{StringMap: map[string]string{"publishable_key": app.Stripe.Key}}, "stripe-js"); err != nil {
		app.ErrorLog.Println(err)
	}
}

func (app *Application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")
	widgetID, _ := strconv.Atoi(r.Form.Get("product_id"))

	card := lib.Card{
		Secret: app.Config.Stripe.Secret,
		Key:    app.Config.Stripe.Key,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntent)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	// create a new customer
	customerID, err := app.SaveCustomer(firstName, lastName, email)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	// create a new transaction
	amount, _ := strconv.Atoi(paymentAmount)
	txn := model.Transactions{
		Amount:              amount,
		Currency:            paymentCurrency,
		LastFour:            lastFour,
		ExpiryMonth:         int(expiryMonth),
		ExpiryYear:          int(expiryYear),
		BankReturnCode:      pi.LatestCharge.ID,
		PaymentIntent:       paymentIntent,
		PaymentMethod:       paymentMethod,
		TransactionStatusID: 7, // Note: 7 -> Cleared
	}

	txnID, err := app.SaveTransaction(txn)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	// create a new order
	order := model.Order{
		WidgetID:      widgetID,
		TransactionID: txnID,
		CustomerID:    customerID,
		StatusID:      4, // Note: 4 -> Cleared
		Quantity:      1,
		Amount:        amount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, err = app.SaveOrder(order)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	data := make(map[string]any)
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency
	data["last_four"] = lastFour
	data["expiry_month"] = expiryMonth
	data["expiry_year"] = expiryYear
	data["bank_return_code"] = pi.LatestCharge.ID
	data["first_name"] = firstName
	data["last_name"] = lastName

	app.Session.Put(r.Context(), "receipt", data)

	http.Redirect(w, r, "/terminal/receipt", http.StatusSeeOther)
}

func (app *Application) Receipt(w http.ResponseWriter, r *http.Request) {
	data := app.Session.Get(r.Context(), "receipt").(map[string]any)
	// app.Session.Remove(r.Context(), "receipt")

	if err := app.RenderTemplate(w, r, "receipt", &template.TemplateData{Data: data}); err != nil {
		app.ErrorLog.Println(err)
	}
}

// SaveCustomer saves a customer and returns id
func (app *Application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := model.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveTransaction saves a txn and returns id
func (app *Application) SaveTransaction(txn model.Transactions) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// SaveOrder saves a order and returns id
func (app *Application) SaveOrder(order model.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}
	return id, nil
}
