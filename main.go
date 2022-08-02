package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v3/plaid"
)

var (
	PLAID_CLIENT_ID                      = ""
	PLAID_SECRET                         = ""
	PLAID_ENV                            = ""
	PLAID_PRODUCTS                       = ""
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	client              *plaid.APIClient = nil
)

func init() {
	// Load env vars from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading environment variables from .env file %w", err)
	}

	// Set constants from env
	PLAID_CLIENT_ID = os.Getenv("PLAID_CLIENT_ID")
	PLAID_SECRET = os.Getenv("PLAID_SECRET")
	if PLAID_CLIENT_ID == "" || PLAID_SECRET == "" {
		log.Fatal("Error: PLAID_SECRET or PLAID_CLIENT_ID is not set. Did you copy .env.example to .env and fill it out?")
	}
	PLAID_ENV = os.Getenv("PLAID_ENV")
	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")

	// Set defaults
	if PLAID_PRODUCTS == "" {
		PLAID_PRODUCTS = "transactions"
	}
	if PLAID_COUNTRY_CODES == "" {
		PLAID_COUNTRY_CODES = "US"
	}
	if PLAID_ENV == "" {
		PLAID_ENV = "sandbox"
	}

	// Create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", PLAID_SECRET)
	configuration.UseEnvironment(plaid.Sandbox) // plaid.Development, plaid.Production
	client = plaid.NewAPIClient(configuration)
}

func main() {
	fmt.Println("Hello World!")

	ctx := context.Background()

	// Get a Sandbox Public Token
	sandboxInstitution := "ins_109508" // First Platypus Bank
	testProducts := []plaid.Products{plaid.Products("auth")}
	sandboxPublicTokenResp, _, _ := client.PlaidApi.SandboxPublicTokenCreate(ctx).SandboxPublicTokenCreateRequest(
		*plaid.NewSandboxPublicTokenCreateRequest(
			sandboxInstitution,
			testProducts,
		),
	).Execute()
	publicToken := sandboxPublicTokenResp.GetPublicToken()
	fmt.Println(publicToken)

	// Exchange the publicToken for an accessToken
	exchangePublicTokenResp, _, _ := client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()
	accessToken := exchangePublicTokenResp.GetAccessToken()
	fmt.Println(accessToken)

	// Get accounts
	accountsGetResp, _, _ := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()
	accounts := accountsGetResp.GetAccounts()
	accountID := accounts[0].GetAccountId()
	fmt.Println(accounts)
	fmt.Println(accountID)

	for _, a := range accounts {
		fmt.Println("\n\nAccountName: ", a.Name, "\nAccountID: ", a.AccountId, "\nBalances:", a.Balances)
	}
}
