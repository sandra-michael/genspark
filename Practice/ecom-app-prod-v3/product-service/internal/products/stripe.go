package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"product-service/pkg/logkey"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
)

func (c *Conf) CreatePricingStripe(ctx context.Context, productId string, val uint64, prodName string) error {
	// Step 1: Retrieve the Stripe secret key from the environment variables
	sKey := os.Getenv("STRIPE_TEST_KEY")
	if sKey == "" {
		// If the key is not set, return an error
		return fmt.Errorf("STRIPE_TEST_KEY not set")
	}

	// Step 2: Assign the Stripe API key to the Stripe library's internal configuration
	stripe.Key = sKey

	// Step 3: Begin a database transaction using the `withTx` method (assumed to be defined elsewhere)
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		// Step 4: Define a SQL query to check if the user already has a Stripe customer ID in the database
		sqlQuery := `
				SELECT stripe_product_id	 
				FROM product_pricing_stripe
				WHERE product_id = $1
				`

		// Step 5: Declare a variable to hold the Stripe customer ID we fetch from the database
		var stripeCustomerId string

		// Step 6: Execute the query to get the Stripe customer ID for the given user ID
		err := tx.QueryRowContext(ctx, sqlQuery, productId).Scan(&stripeCustomerId)
		if err != nil {
			// Step 7: Handle the case where no rows are found (i.e., the user doesn't have a Stripe customer ID yet)
			if !errors.Is(err, sql.ErrNoRows) {
				// If the error is not `sql.ErrNoRows`, that means something went wrong with the query execution; return an error
				return fmt.Errorf("failed to fetch Stripe customer ID: %w", err)
			}

			params := &stripe.PriceParams{
				Currency:   stripe.String(string(stripe.CurrencyINR)),
				UnitAmount: stripe.Int64(int64(val)),

				ProductData: &stripe.PriceProductDataParams{Name: stripe.String(prodName)},
			}

			// Step 9: Call the Stripe API to create a new customer using the parameters
			priceResult, err := price.New(params)

			if err != nil {
				// Log the error and return it if the creation of the Stripe customer fails
				slog.Error("failed to create Stripe pricnig", slog.Any(logkey.ERROR, err))
				return fmt.Errorf("failed to create Stripe pricnig: %w", err)
			}

			// Step 10: Define the SQL query to insert the new Stripe customer into the `users_stripe` table
			query := `
		INSERT INTO product_pricing_stripe (product_id, stripe_product_id, price_id, price,created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
			// Step 11: Get the current timestamp for record creation and updates
			createdAt := time.Now().UTC()
			updatedAt := createdAt

			fmt.Println(query, productId, priceResult.Product, priceResult.ID, createdAt, updatedAt)

			fmt.Println(priceResult)

			// Step 12: Execute the query to insert the new Stripe customer into the database
			res, err := tx.ExecContext(ctx, query, productId, priceResult.Product.ID, priceResult.ID, val, createdAt, updatedAt)
			if err != nil {
				// Log the error and return it if the database insertion fails
				slog.Error("failed to insert Stripe pricing ID", slog.Any(logkey.ERROR, err))
				return fmt.Errorf("failed to insert Stripe customer ID: %w", err)
			}

			// Step 13: Check if the insertion affected any rows (it should affect exactly one row if successful)
			if num, err := res.RowsAffected(); num == 0 || err != nil {
				// If no rows were affected or another error occurred, return an error
				return fmt.Errorf("failed to insert Stripe customer ID: %w", err)
			}

			// Step 14: Return `nil` if the Stripe customer is successfully added to the database
			return nil
		}

		// Step 15: CustomerId already exist on stripe, no need to add the customer
		return nil
	})

	// Step 16: Handle any errors from the transaction function
	if err != nil {
		return err
	}

	// Step 17: If everything succeeds, return `nil`
	return nil
}

func (c *Conf) GetStripeProductDetails(ctx context.Context, productId string) (ProductOrder, error) {
	var prodOrder ProductOrder
	//var stock int
	//var price_id string

	// SQL query to retrieve the Stripe customer ID for the given user ID
	query := `
	select pr.stock as stock, pps.price_id as price_id
	from products pr
	inner join product_pricing_stripe pps on pr.id = pps.product_id
	where pr.id = $1
	`
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, query, productId).Scan(&prodOrder.Stock, &prodOrder.PriceId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("no stripe price id  found for product %s: %w", productId, err)
			}
			return fmt.Errorf("failed to fetch stripe price id: %w", err)
		}
		return nil
	})
	if err != nil {
		return ProductOrder{}, err
	}

	return prodOrder, nil

}
