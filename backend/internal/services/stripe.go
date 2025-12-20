package services

import (
	"fmt"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"

	"github.com/gecogreen/backend/internal/config"
	"github.com/gecogreen/backend/internal/models"
)

type StripeService struct {
	config *config.Config
}

func NewStripeService(cfg *config.Config) *StripeService {
	stripe.Key = cfg.StripeSecretKey
	return &StripeService{config: cfg}
}

// CreateCheckoutSession creates a Stripe Checkout Session for an order
func (s *StripeService) CreateCheckoutSession(order *models.Order, product *models.Product, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	// Calculate amounts in cents (Stripe uses smallest currency unit)
	unitAmountCents := int64(order.UnitPrice * 100)
	shippingCents := int64(order.ShippingCost * 100)

	lineItems := []*stripe.CheckoutSessionLineItemParams{
		{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("eur"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(product.Title),
					Description: stripe.String(truncate(product.Description, 500)),
				},
				UnitAmount: stripe.Int64(unitAmountCents),
			},
			Quantity: stripe.Int64(int64(order.Quantity)),
		},
	}

	// Add shipping as separate line item if applicable
	if shippingCents > 0 {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("eur"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Spedizione"),
				},
				UnitAmount: stripe.Int64(shippingCents),
			},
			Quantity: stripe.Int64(1),
		})
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		LineItems:          lineItems,
		ExpiresAt:          stripe.Int64(order.CreatedAt.Add(30 * 60 * 1000000000).Unix()), // 30 minutes
		Metadata: map[string]string{
			"order_id":   order.ID.String(),
			"buyer_id":   order.BuyerID.String(),
			"seller_id":  order.SellerID.String(),
			"product_id": order.ProductID.String(),
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"order_id": order.ID.String(),
			},
		},
	}

	// Collect shipping address if shipping is required
	if order.DeliveryType == models.DeliverySellerShips {
		params.ShippingAddressCollection = &stripe.CheckoutSessionShippingAddressCollectionParams{
			AllowedCountries: stripe.StringSlice([]string{"IT"}),
		}
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess, nil
}

// VerifyWebhookSignature verifies a Stripe webhook signature
func (s *StripeService) VerifyWebhookSignature(payload []byte, signature string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signature, s.config.StripeWebhookSecret)
}

// Helper function to truncate strings
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
