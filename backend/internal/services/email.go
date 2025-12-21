package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/gecogreen/backend/internal/config"
	"github.com/gecogreen/backend/internal/models"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

// SendOrderConfirmationToBuyer sends order confirmation to the buyer
func (s *EmailService) SendOrderConfirmationToBuyer(order *models.Order, product *models.Product, buyerEmail, buyerName string) error {
	subject := fmt.Sprintf("Conferma ordine #%s - GecoGreen", order.ID.String()[:8])

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #22c55e, #16a34a); color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 20px; border: 1px solid #e5e7eb; }
        .product { background: white; padding: 15px; border-radius: 8px; margin: 15px 0; }
        .footer { text-align: center; padding: 20px; color: #6b7280; font-size: 12px; }
        .btn { display: inline-block; background: #22c55e; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin-top: 15px; }
        .total { font-size: 24px; color: #22c55e; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Ordine Confermato!</h1>
        </div>
        <div class="content">
            <p>Ciao <strong>%s</strong>,</p>
            <p>Grazie per il tuo acquisto su GecoGreen! Il tuo ordine e' stato confermato con successo.</p>

            <div class="product">
                <h3>%s</h3>
                <p>Quantita': %d</p>
                <p>Prezzo unitario: %.2f EUR</p>
                %s
            </div>

            <p class="total">Totale: %.2f EUR</p>

            %s

            <p style="margin-top: 20px;">
                <a href="%s/orders/%s" class="btn">Visualizza Ordine</a>
            </p>
        </div>
        <div class="footer">
            <p>GecoGreen - La piattaforma antispreco</p>
            <p>Insieme contro lo spreco alimentare</p>
        </div>
    </div>
</body>
</html>
`,
		buyerName,
		product.Title,
		order.Quantity,
		order.UnitPrice,
		s.getShippingInfo(order),
		order.TotalAmount,
		s.getDeliveryInfo(order),
		s.config.FrontendURL,
		order.ID.String(),
	)

	return s.sendEmail(buyerEmail, subject, body)
}

// SendNewOrderToSeller notifies the seller of a new order
func (s *EmailService) SendNewOrderToSeller(order *models.Order, product *models.Product, sellerEmail, sellerName, buyerName string) error {
	subject := fmt.Sprintf("Nuovo ordine ricevuto! #%s - GecoGreen", order.ID.String()[:8])

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #3b82f6, #1d4ed8); color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 20px; border: 1px solid #e5e7eb; }
        .product { background: white; padding: 15px; border-radius: 8px; margin: 15px 0; }
        .footer { text-align: center; padding: 20px; color: #6b7280; font-size: 12px; }
        .btn { display: inline-block; background: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin-top: 15px; }
        .total { font-size: 24px; color: #22c55e; font-weight: bold; }
        .alert { background: #fef3c7; border: 1px solid #f59e0b; padding: 15px; border-radius: 8px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Hai ricevuto un nuovo ordine!</h1>
        </div>
        <div class="content">
            <p>Ciao <strong>%s</strong>,</p>
            <p>Ottime notizie! Hai ricevuto un nuovo ordine su GecoGreen.</p>

            <div class="product">
                <h3>%s</h3>
                <p><strong>Acquirente:</strong> %s</p>
                <p><strong>Quantita':</strong> %d</p>
                <p><strong>Prezzo unitario:</strong> %.2f EUR</p>
            </div>

            <p class="total">Totale ordine: %.2f EUR</p>

            %s

            <div class="alert">
                <strong>Prossimi passi:</strong>
                <p>%s</p>
            </div>

            <p style="margin-top: 20px;">
                <a href="%s/seller/orders/%s" class="btn">Gestisci Ordine</a>
            </p>
        </div>
        <div class="footer">
            <p>GecoGreen - La piattaforma antispreco</p>
        </div>
    </div>
</body>
</html>
`,
		sellerName,
		product.Title,
		buyerName,
		order.Quantity,
		order.UnitPrice,
		order.TotalAmount,
		s.getShippingAddressForSeller(order),
		s.getNextStepsForSeller(order),
		s.config.FrontendURL,
		order.ID.String(),
	)

	return s.sendEmail(sellerEmail, subject, body)
}

func (s *EmailService) getShippingInfo(order *models.Order) string {
	if order.ShippingCost > 0 {
		return fmt.Sprintf("<p>Spedizione: %.2f EUR</p>", order.ShippingCost)
	}
	return ""
}

func (s *EmailService) getDeliveryInfo(order *models.Order) string {
	if order.DeliveryType == models.DeliveryPickup {
		return `<p><strong>Modalita':</strong> Ritiro in sede</p>
		<p>Il venditore ti contattera' per organizzare il ritiro.</p>`
	}
	return fmt.Sprintf(`<p><strong>Spedizione a:</strong></p>
	<p>%s<br>%s %s %s</p>`,
		order.ShippingAddress,
		order.ShippingPostalCode,
		order.ShippingCity,
		order.ShippingProvince,
	)
}

func (s *EmailService) getShippingAddressForSeller(order *models.Order) string {
	if order.DeliveryType == models.DeliveryPickup {
		return ""
	}
	return fmt.Sprintf(`<div class="product">
		<h4>Indirizzo di spedizione:</h4>
		<p>%s<br>%s %s %s<br>%s</p>
	</div>`,
		order.ShippingAddress,
		order.ShippingPostalCode,
		order.ShippingCity,
		order.ShippingProvince,
		order.ShippingCountry,
	)
}

func (s *EmailService) getNextStepsForSeller(order *models.Order) string {
	if order.DeliveryType == models.DeliveryPickup {
		return "Contatta l'acquirente per organizzare il ritiro del prodotto."
	}
	return "Prepara il pacco e spedisci il prodotto all'indirizzo indicato. Ricordati di inserire il numero di tracking!"
}

// sendEmail sends an email via SMTP with STARTTLS support
func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	if s.config.SMTPUser == "" || s.config.SMTPPassword == "" {
		fmt.Printf("SMTP not configured, skipping email to %s\n", to)
		return nil
	}

	from := s.config.SMTPFrom

	// Build email message
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	tlsConfig := &tls.Config{
		ServerName: s.config.SMTPHost,
	}

	var client *smtp.Client
	var err error

	// Port 465 uses implicit TLS, port 587 uses STARTTLS
	if s.config.SMTPPort == "465" {
		// Implicit TLS (SMTPS)
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server (TLS): %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, s.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
	} else {
		// STARTTLS (port 587 or 25)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, s.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}

		// Upgrade to TLS
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	// Set sender and recipient
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}
	_, err = w.Write([]byte(message.String()))
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	client.Quit()

	fmt.Printf("Email sent successfully to %s\n", to)
	return nil
}
