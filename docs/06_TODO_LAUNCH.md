# GecoGreen - Checklist Pre-Lancio

## Fase 0: Preparazione (Settimana 1)

### Identità e Dominio
- [x] ~~Verificare disponibilità dominio~~ → **www.gecogreen.com ACQUISTATO**
- [x] ~~Acquistare dominio~~ → Completato
- [ ] Registrare account social (Instagram, Facebook, LinkedIn)
- [ ] Verificare trademark "GecoGreen" (UIBM)

### Setup Legale
- [ ] Decidere forma societaria (Ditta individuale / SRL / etc.)
- [ ] Consulenza commercialista per setup fiscale
- [ ] Aprire P.IVA se necessario
- [ ] Aprire conto corrente dedicato

---

## Fase 1: Infrastruttura Tech (Settimana 2)

### Hosting e Server
- [ ] Creare account Hetzner
- [ ] Acquistare VPS (CPX21 consigliato)
- [ ] Installare Ubuntu 22.04
- [ ] Installare Coolify
- [ ] Configurare dominio → VPS
- [ ] Setup SSL automatico (Let's Encrypt)

### Repository e CI/CD
- [ ] Creare repo GitHub privato
- [ ] Collegare repo a Coolify
- [ ] Configurare deploy automatico
- [ ] Setup branch protection (main)

### Database
- [ ] Installare PostgreSQL su VPS
- [ ] Creare database "gecogreen"
- [ ] Eseguire script schema iniziale
- [ ] Configurare backup automatico

### Ambiente Locale
- [ ] Installare Docker Desktop
- [ ] Creare docker-compose.yml
- [ ] Verificare ambiente funzionante
- [ ] Documentare setup per team futuro

---

## Fase 2: Pagamenti (Settimana 2-3)

### Stripe Setup
- [ ] Creare account Stripe
- [ ] Attivare Stripe Connect
- [ ] Configurare webhook endpoint
- [ ] Testare flusso pagamento (test mode)
- [ ] Testare flusso payout (test mode)
- [ ] Configurare email ricevute Stripe

### Verifica Compliance
- [ ] Verificare requisiti KYC Stripe Italia
- [ ] Preparare documenti necessari
- [ ] Completare verifica account Stripe

---

## Fase 3: Documenti Legali (Settimana 3)

### Termini e Condizioni
- [ ] Redigere T&C base (o usare template + revisione)
- [ ] Includere clausole marketplace
- [ ] Includere clausole dispute/rimborsi
- [ ] Includere limitazioni responsabilità
- [ ] Far revisionare da legale

### Privacy Policy
- [ ] Redigere Privacy Policy GDPR compliant
- [ ] Elencare dati raccolti
- [ ] Specificare base legale trattamento
- [ ] Includere diritti utente
- [ ] Nominare DPO se necessario

### Cookie Policy
- [ ] Creare Cookie Policy
- [ ] Implementare cookie banner
- [ ] Gestione consensi

---

## Fase 4: Sviluppo MVP (Settimane 3-8)

### Backend Core
- [ ] Setup progetto Go + Fiber
- [ ] Implementare autenticazione JWT
- [ ] Implementare CRUD utenti
- [ ] Implementare CRUD prodotti
- [ ] Implementare sistema ordini
- [ ] Implementare chat base
- [ ] Integrare Stripe Connect
- [ ] Implementare QR code ritiro
- [ ] Implementare moderazione chat (Livello 1)
- [ ] Scrivere test API principali

### Frontend Core
- [ ] Setup progetto SvelteKit + TailwindCSS + DaisyUI
- [ ] Configurare Capacitor per mobile
- [ ] Implementare layout base (tema Verde Geco #00C853)
- [ ] Implementare autenticazione
- [ ] Pagina catalogo con filtri
- [ ] Pagina prodotto
- [ ] Carrello e checkout
- [ ] Dashboard buyer
- [ ] Dashboard seller
- [ ] Sistema chat
- [ ] Responsive mobile

### Admin Panel
- [ ] Dashboard statistiche
- [ ] Gestione utenti
- [ ] Gestione dispute
- [ ] Log sistema

---

## Fase 5: Testing (Settimana 8-9)

### Test Funzionali
- [ ] Test registrazione buyer
- [ ] Test registrazione seller
- [ ] Test flusso acquisto completo
- [ ] Test flusso ritiro con QR
- [ ] Test flusso spedizione
- [ ] Test apertura disputa
- [ ] Test chat
- [ ] Test pagamenti reali (piccoli importi)

### Test Sicurezza
- [ ] Verificare autenticazione robusta
- [ ] Testare rate limiting
- [ ] Verificare input validation
- [ ] Test CORS
- [ ] Verificare SQL injection protection

### Test Performance
- [ ] Load test base (50 utenti concorrenti)
- [ ] Verificare tempi risposta API < 200ms
- [ ] Ottimizzare query lente

---

## Fase 6: Pre-Lancio (Settimana 9-10)

### Contenuti
- [ ] Scrivere copy homepage
- [ ] Scrivere pagina "Come funziona"
- [ ] Scrivere FAQ
- [ ] Preparare email templates
- [ ] Creare immagini placeholder prodotti

### SEO Base
- [ ] Configurare meta tags
- [ ] Creare sitemap.xml
- [ ] Creare robots.txt
- [ ] Verificare su Google Search Console

### Monitoring
- [ ] Setup Sentry per error tracking
- [ ] Setup UptimeRobot per uptime monitoring
- [ ] Configurare alert email per errori critici

### Outreach Venditori
- [ ] Creare lista 50 potenziali venditori
- [ ] Preparare pitch email/telefonata
- [ ] Contattare primi 20 venditori
- [ ] Onboarding primi 5-10 venditori "beta"

---

## Fase 7: Lancio Soft (Settimana 10-11)

### Go-Live Checklist
- [ ] Backup database pre-lancio
- [ ] Deploy versione finale
- [ ] Verificare tutte le funzioni in produzione
- [ ] Monitorare errori prime 24h
- [ ] Rispondere rapidamente a bug critici

### Comunicazione
- [ ] Post social "Siamo online!"
- [ ] Email a beta tester
- [ ] Comunicare a venditori onboardati

---

## Fase 8: Marketing Launch (Settimane 11-12)

### Paid Ads
- [ ] Creare campagna Facebook/Instagram
- [ ] Creare campagna Google Ads
- [ ] Setup pixel tracking
- [ ] A/B test creatività
- [ ] Monitorare CPL giornalmente

### Content
- [ ] Pubblicare primo articolo blog
- [ ] Post social 3x settimana
- [ ] Rispondere a commenti/DM

### PR
- [ ] Contattare blog sostenibilità
- [ ] Contattare giornali locali
- [ ] Cercare podcast rilevanti

---

## Checklist Giornaliera Post-Lancio

### Ogni Giorno
- [ ] Controllare errori Sentry
- [ ] Controllare uptime
- [ ] Rispondere a ticket/email
- [ ] Moderare dispute aperte
- [ ] Controllare metriche base

### Ogni Settimana
- [ ] Review KPI
- [ ] Analisi costi marketing
- [ ] Feedback call con venditori
- [ ] Planning settimana successiva

### Ogni Mese
- [ ] Report finanziario
- [ ] Fatturazione commissioni
- [ ] Review roadmap prodotto
- [ ] Analisi churn

---

## Contatti Utili

| Servizio | Link | Note |
|----------|------|------|
| Hetzner | hetzner.com | Hosting |
| Stripe | stripe.com | Pagamenti |
| Namecheap | namecheap.com | Domini |
| Coolify | coolify.io | PaaS |
| Sentry | sentry.io | Error tracking |
| UptimeRobot | uptimerobot.com | Monitoring |
| Canva | canva.com | Grafica |
| Fiverr | fiverr.com | Logo/design |

---

## Note Importanti

1. **Non lanciare senza T&C e Privacy Policy** - Rischi legali alti
2. **Non lanciare senza Stripe verificato** - Non puoi incassare
3. **Testare pagamenti reali prima del lancio** - Usa piccoli importi
4. **Avere almeno 5 venditori attivi al lancio** - Serve contenuto
5. **Budget marketing: non spendere tutto subito** - Testa e ottimizza

---

*Documento creato: Dicembre 2024*
*Da aggiornare man mano che si completano i task*
