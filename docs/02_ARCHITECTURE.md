# GecoGreen - Architettura Tecnica

## Stack Tecnologico

```
┌─────────────────────────────────────────────────────────────┐
│                        FRONTEND                              │
│        SvelteKit + TypeScript + TailwindCSS + DaisyUI        │
│              + Capacitor (iOS/Android wrapper)               │
│                    Colore: Verde Geco #00C853                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                         BACKEND                              │
│                    Go (Fiber Framework)                      │
│                  REST API + WebSocket Chat                   │
│                    (Docker su Hetzner)                       │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│   PostgreSQL     │ │   Redis          │ │  Cloudflare R2   │
│   (Database)     │ │   (Cache/Session)│ │   (Immagini)     │
└──────────────────┘ └──────────────────┘ └──────────────────┘

Dominio: www.gecogreen.com
```

---

## Dettaglio Componenti

### Frontend: SvelteKit + DaisyUI + Capacitor

**Perché SvelteKit:**
- Compilato (bundle piccolissimi, veloce)
- Server-Side Rendering nativo (SEO perfetto)
- Meno boilerplate di React/Next.js
- L'AI genera codice Svelte molto pulito

**UI Framework: TailwindCSS + DaisyUI**
- DaisyUI: Componenti pre-stilizzati per Tailwind
- Tema personalizzato: Verde Geco (#00C853)
- Dark mode supportato nativamente

**Mobile: Capacitor**
- Wrapper per iOS e Android
- Stesso codebase SvelteKit
- Accesso a API native (camera, notifiche push)

**Struttura:**
```
frontend/
├── src/
│   ├── routes/
│   │   ├── [[lang]]/             # Routing localizzato (opzionale)
│   │   │   ├── +page.svelte      # Homepage
│   │   │   ├── catalogo/
│   │   │   ├── prodotto/[id]/
│   │   │   ├── carrello/
│   │   │   ├── checkout/
│   │   │   └── ...
│   │   ├── auth/
│   │   │   ├── login/
│   │   │   └── register/
│   │   ├── dashboard/
│   │   │   ├── buyer/
│   │   │   ├── seller/
│   │   │   └── admin/
│   │   └── api/                  # API routes (proxy)
│   ├── lib/
│   │   ├── components/
│   │   ├── stores/
│   │   ├── i18n/                 # Sistema traduzioni
│   │   │   ├── index.ts          # Setup i18n
│   │   │   └── translations/
│   │   │       ├── it.json       # Italiano (default)
│   │   │       ├── en.json       # English
│   │   │       ├── de.json       # Deutsch
│   │   │       ├── fr.json       # Français
│   │   │       └── es.json       # Español
│   │   └── utils/
│   └── app.html
├── static/
├── tailwind.config.js
└── svelte.config.js
```

### Backend: Go + Fiber

**Perché Go:**
- Compilato, velocissimo
- Gestione concorrenza nativa (perfetto per ordini simultanei)
- Tipizzazione forte (l'AI non sbaglia i tipi)
- Deploy semplice (un binario)

**Struttura:**
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── products.go
│   │   ├── orders.go
│   │   ├── chat.go
│   │   └── admin.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── ratelimit.go
│   │   └── moderation.go
│   ├── models/
│   ├── repository/
│   ├── services/
│   └── utils/
├── pkg/
│   ├── stripe/
│   ├── moderation/
│   └── email/
├── migrations/
├── Dockerfile
└── go.mod
```

### Database: PostgreSQL

**Perché PostgreSQL:**
- Robusto per dati finanziari
- JSONB per dati flessibili (foto, metadata)
- Full-text search nativo
- Enum types nativi

**Schema completo:** Vedi `schemas/database.sql`

### Cache: Redis

**Uso:**
- Sessioni utente
- Rate limiting
- Cache catalogo (prodotti frequenti)
- Coda messaggi chat

---

## Infrastruttura

### Ambiente Sviluppo (Locale)

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: gecogreen
      POSTGRES_USER: gecogreen
      POSTGRES_PASSWORD: dev_password
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    environment:
      DATABASE_URL: postgres://gecogreen:dev_password@postgres:5432/gecogreen
      REDIS_URL: redis://redis:6379

volumes:
  pgdata:
```

### Ambiente Produzione (Hetzner + Coolify)

**VPS Consigliato:** Hetzner CPX21
- 3 vCPU
- 4 GB RAM
- 80 GB SSD
- ~7€/mese

**Setup:**
1. VPS Hetzner con Ubuntu 22.04
2. Installazione Coolify (self-hosted PaaS)
3. Collegamento repo GitHub
4. Deploy automatico su push

**Coolify gestisce:**
- Certificati SSL (Let's Encrypt)
- Reverse proxy (Traefik)
- Container Docker
- Backup database

---

## Integrazioni Esterne

### Stripe Connect

**Flusso:**
```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Buyer   │────▶│ GecoGreen│────▶│  Stripe  │────▶│  Seller  │
│  paga    │     │ (passa)  │     │ (split)  │     │ (riceve) │
│  100€    │     │          │     │ 90€+10€  │     │   90€    │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
```

**Endpoints necessari:**
- `POST /api/stripe/onboard` - Onboarding seller
- `POST /api/stripe/checkout` - Crea sessione pagamento
- `POST /api/stripe/webhook` - Riceve eventi

### Moderazione Chat

**3 Livelli:**

1. **RegEx (Gratis):**
   - Blocca numeri telefono
   - Blocca email
   - Blocca parole chiave (whatsapp, telegram, contanti)

2. **OpenAI Moderation API (Gratis):**
   - Controlla hate speech
   - Controlla violenza
   - Controlla contenuti illegali

3. **AI Mini - Solo se necessario (~0.01€/1000 msg):**
   - Analisi contestuale messaggi sospetti
   - Solo per utenti nuovi o già flaggati

### File Storage: Cloudflare R2

**Perché Cloudflare R2:**
- S3-compatible (facile migrazione)
- Nessun costo egress (download gratuiti)
- CDN integrato (immagini veloci)
- 10 GB gratis/mese

**Uso:**
- Foto prodotti (resize automatico)
- PDF certificati impatto
- Avatar utenti

### Email Transazionali: Resend

**Provider scelto: Resend**
- 3.000 email/mese gratis
- API semplice
- DNS: Aruba (solo inbox), Resend per invio
- Ottima deliverability

**Email da inviare:**
- Conferma registrazione
- Reset password
- Conferma ordine
- Promemoria ritiro
- Fattura commissioni (PDF allegato)

---

## API Structure

### Pubbliche (No Auth)
```
GET  /api/products              # Lista prodotti (paginata)
GET  /api/products/:id          # Dettaglio prodotto
GET  /api/categories            # Lista categorie
POST /api/auth/register         # Registrazione
POST /api/auth/login            # Login
POST /api/auth/forgot-password  # Reset password
```

### Buyer (Auth Required)
```
GET  /api/buyer/orders          # I miei ordini
POST /api/buyer/orders          # Crea ordine
GET  /api/buyer/orders/:id      # Dettaglio ordine
POST /api/buyer/orders/:id/confirm  # Conferma ricezione
POST /api/buyer/orders/:id/dispute  # Apri contestazione
GET  /api/buyer/chat/:orderId   # Messaggi ordine
POST /api/buyer/chat/:orderId   # Invia messaggio
```

### Seller (Auth + Role)
```
GET    /api/seller/products         # I miei prodotti
POST   /api/seller/products         # Nuovo prodotto
PUT    /api/seller/products/:id     # Modifica prodotto
DELETE /api/seller/products/:id     # Elimina prodotto
GET    /api/seller/orders           # Ordini ricevuti
PUT    /api/seller/orders/:id       # Aggiorna stato
POST   /api/seller/orders/:id/qr    # Scansiona QR ritiro
GET    /api/seller/wallet           # Saldo e movimenti
GET    /api/seller/invoices         # Fatture commissioni
```

### Admin (Auth + Admin Role)
```
GET  /api/admin/users           # Lista utenti
PUT  /api/admin/users/:id       # Modifica utente (ban, etc)
GET  /api/admin/disputes        # Lista dispute aperte
PUT  /api/admin/disputes/:id    # Risolvi disputa
GET  /api/admin/stats           # Dashboard statistiche
GET  /api/admin/logs            # Audit log
```

---

## Sicurezza

### Autenticazione
- JWT tokens (access + refresh)
- HttpOnly cookies per refresh token
- Rate limiting su login (5 tentativi/15 min)

### OAuth / Social Login

**Provider supportati:**
- Google Sign-In (obbligatorio)
- Apple Sign-In (obbligatorio per App Store)
- Facebook Login (opzionale, futuro)

**Flusso OAuth:**
```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│  User    │──▶│ Frontend │──▶│ Provider │──▶│ Backend  │
│  click   │   │ redirect │   │ (Google) │   │ callback │
└──────────┘   └──────────┘   └──────────┘   └──────────┘
                                                  │
                                                  ▼
                                          ┌──────────────┐
                                          │ Trova/Crea   │
                                          │ utente by    │
                                          │ google_id    │
                                          └──────────────┘
                                                  │
                                                  ▼
                                          ┌──────────────┐
                                          │ Genera JWT   │
                                          │ + redirect   │
                                          └──────────────┘
```

**Librerie:**
- Backend Go: `golang.org/x/oauth2`
- Frontend: SDK nativi Google/Apple

**Note implementazione:**
- Se utente esiste con stessa email ma senza google_id → Collegare account
- Se nuovo utente → Creare con email_verified=true (già verificato da Google)
- Avatar URL salvato automaticamente da profilo OAuth
- Password non richiesta per utenti OAuth (campo nullable)

### Autorizzazione
- Role-based (BUYER, SELLER, ADMIN)
- Middleware per controllo permessi
- Seller può vedere solo i propri prodotti/ordini

### Dati Sensibili
- Password: bcrypt con cost 12
- Dati pagamento: Mai salvati, gestiti da Stripe
- PII: Criptati a riposo (AES-256)

### Protezione API
- CORS configurato
- Helmet headers
- Input validation (Fiber validator)
- SQL injection: ORM con prepared statements

---

## Monitoring

### Logging
- Structured logging (JSON)
- Livelli: DEBUG, INFO, WARN, ERROR
- Aggregazione: Loki o file locale

### Metriche
- Prometheus metrics endpoint
- Grafana dashboard (opzionale)

### Alerting
- Uptime: UptimeRobot (gratis)
- Errori: Sentry (gratis tier)

---

## Internazionalizzazione (i18n)

### Strategia Multilingua

**Fase 1 (MVP):** Solo italiano
**Fase 2:** Inglese + altre lingue EU

**Libreria scelta: Paraglide.js (Inlang)**

Perché Paraglide invece di svelte-i18n:
- Typesafe (errori a compile time)
- Tree-shakeable (solo traduzioni usate nel bundle)
- Integrazione nativa con SvelteKit
- IDE support (VS Code extension)
- Extraction automatica delle chiavi

### Struttura File Traduzioni

```
frontend/src/lib/i18n/
├── translations/
│   ├── it.json          # ~500 chiavi (Italiano - default)
│   ├── en.json          # English
│   ├── de.json          # Deutsch
│   ├── fr.json          # Français
│   └── es.json          # Español
└── index.ts             # Setup e export
```

### Organizzazione Chiavi (per file mantenibile)

Struttura JSON gerarchica per sezione:

```json
{
  "common": {
    "save": "Salva",
    "cancel": "Annulla",
    "loading": "Caricamento...",
    "error": "Si è verificato un errore"
  },
  "auth": {
    "login": "Accedi",
    "register": "Registrati",
    "logout": "Esci",
    "forgotPassword": "Password dimenticata?"
  },
  "products": {
    "title": "Prodotti",
    "addToCart": "Aggiungi al carrello",
    "outOfStock": "Esaurito",
    "expiresIn": "Scade tra {days} giorni",
    "dutchAuction": {
      "priceDrops": "Il prezzo scende di {amount}€ ogni {hours}h",
      "minPrice": "Prezzo minimo: {price}€",
      "nextDrop": "Prossimo ribasso tra: {time}"
    }
  },
  "checkout": {
    "title": "Checkout",
    "shipping": {
      "pickup": "Ritiro in sede",
      "sellerShips": "Spedizione venditore",
      "platformManaged": "Spedizione GecoGreen"
    }
  },
  "seller": {
    "dashboard": "Dashboard Venditore",
    "newProduct": "Nuovo Prodotto",
    "orders": "Ordini"
  },
  "errors": {
    "required": "Campo obbligatorio",
    "invalidEmail": "Email non valida",
    "minLength": "Minimo {min} caratteri"
  }
}
```

### Uso nei Componenti Svelte

```svelte
<script>
  import { t } from '$lib/i18n';
</script>

<!-- Testo semplice -->
<button>{$t('common.save')}</button>

<!-- Con parametri -->
<p>{$t('products.expiresIn', { days: 3 })}</p>

<!-- Plurali -->
<p>{$t('cart.items', { count: items.length })}</p>
```

### URL Localizzati

**Opzione A: Prefisso lingua (consigliata per SEO)**
```
gecogreen.com/it/catalogo     → Italiano
gecogreen.com/en/catalog      → English
gecogreen.com/de/katalog      → Deutsch
```

**Opzione B: Dominio/subdomain**
```
gecogreen.it                  → Italiano
gecogreen.com                 → English
de.gecogreen.com              → Deutsch
```

**Scelta:** Opzione A con routing `[[lang]]` di SvelteKit

### Rilevamento Lingua

```typescript
// Ordine di priorità:
1. URL path (/en/catalog)
2. Cookie preferenza utente
3. Accept-Language header browser
4. Default: 'it' (italiano)
```

### SEO Multilingua

```html
<!-- In app.html o +layout.svelte -->
<link rel="alternate" hreflang="it" href="https://gecogreen.com/it/..." />
<link rel="alternate" hreflang="en" href="https://gecogreen.com/en/..." />
<link rel="alternate" hreflang="de" href="https://gecogreen.com/de/..." />
<link rel="alternate" hreflang="x-default" href="https://gecogreen.com/it/..." />
```

### Backend i18n

Il backend invia sempre dati grezzi. La traduzione avviene solo nel frontend.

**Eccezioni (tradotte lato server):**
- Email transazionali (template per lingua)
- PDF certificati impatto
- Notifiche push

```go
// Email service
func SendOrderConfirmation(order Order, lang string) {
    template := loadTemplate("order_confirmation", lang) // order_confirmation_it.html
    // ...
}
```

### Workflow Traduzioni

```
1. Dev scrive UI in italiano
2. Estrazione automatica chiavi (Inlang CLI)
3. Export JSON per traduttori
4. Traduttori completano via Inlang Editor (web)
5. Import traduzioni nel repo
6. CI verifica chiavi mancanti
```

### Lingue Supportate (Roadmap)

| Lingua | Codice | Fase | Priorità |
|--------|--------|------|----------|
| Italiano | it | MVP | 🟢 Attiva |
| English | en | 2 | 🟡 Pianificata |
| Deutsch | de | 2 | 🟡 Pianificata |
| Français | fr | 3 | ⚪ Futura |
| Español | es | 3 | ⚪ Futura |

---

*Documento creato: Dicembre 2024*
*Stack version: v1.0*
