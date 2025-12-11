# GecoGreen - Regole di Business e Antifrode

## Sistema Pagamenti

### Stripe Connect (Escrow)

Il denaro NON passa mai direttamente al venditore. Usiamo Stripe Connect in modalità **Escrow**.

```
┌──────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────┐
│  BUYER   │────▶│  STRIPE      │────▶│  HOLD 48h    │────▶│  SELLER  │
│  paga    │     │  (incassa)   │     │  (buffer)    │     │  riceve  │
│  100€    │     │              │     │              │     │  ~87€    │
└──────────┘     └──────────────┘     └──────────────┘     └──────────┘
                        │
                        ▼
                 ┌──────────────┐
                 │  GECOGREEN   │
                 │  commissione │
                 │  ~10€        │
                 └──────────────┘
```

### Flusso Denaro Dettagliato

**Acquisto 100€:**
| Voce | Importo |
|------|---------|
| Prezzo prodotto | 100,00€ |
| Fee Stripe (1.4% + 0.25€) | -1,65€ |
| Commissione GecoGreen (10%) | -10,00€ |
| **Netto al Seller** | **88,35€** |

### Tempistiche Payout

| Evento | Tempo |
|--------|-------|
| Pagamento buyer | T0 |
| Conferma ricezione (o timeout) | T0 + max 7gg |
| Buffer contestazione | +48h |
| Payout a seller | T0 + 9gg max |

> Se il buyer conferma subito, il payout parte dopo 48h dalla conferma.

---

## Logistica

### Opzione 1: Ritiro (PICKUP)

**Regole:**
- Indirizzo nascosto fino al pagamento
- QR Code generato per il buyer
- Seller DEVE scansionare QR per confermare consegna
- Timer: 7 giorni per ritirare

**Privacy:**
- Pre-pagamento: "Disponibile a Milano, zona Lambrate" (solo CAP)
- Post-pagamento: Via completa, civico, orari, telefono

### Opzione 2: Spedizione Seller (SELLER_SHIPS)

**Regole:**
- Seller indica costo spedizione nel prodotto
- Buyer paga: prezzo + spedizione
- Seller organizza e paga il corriere
- Tracking OBBLIGATORIO

**Flusso:**
1. Buyer paga tutto
2. Seller spedisce entro 3 giorni lavorativi
3. Seller inserisce numero tracking
4. Sistema monitora tracking (se API disponibile)
5. Consegna → 48h buffer → Payout

### Opzione 3: Trasporto Buyer (BUYER_ARRANGES)

**Regole:**
- Buyer paga solo il prodotto
- Buyer organizza ritiro (suo corriere/camion)
- Seller deve solo preparare la merce
- Deadline concordata in chat

**Uso tipico:** Merce pesante, pallet, macchinari

---

## Sistema Antifrode

### 1. QR Code Ritiro

**Problema:** Come dimostrare che il ritiro è avvenuto?

**Soluzione:**
```
┌────────────────────────────────────────┐
│           QR CODE STRUCTURE            │
├────────────────────────────────────────┤
│ order_id: UUID                         │
│ buyer_id: UUID                         │
│ created_at: timestamp                  │
│ expires_at: timestamp (+7 giorni)      │
│ signature: HMAC-SHA256                 │
└────────────────────────────────────────┘
```

- QR firmato crittograficamente
- Valido solo per quell'ordine
- Scadenza automatica
- Non riutilizzabile

**Scansione:**
- Seller apre scanner nell'app
- Inquadra QR del buyer
- Sistema verifica firma
- Ordine passa a COMPLETED
- Soldi sbloccati (dopo buffer)

### 2. Timer e Scadenze

| Timer | Durata | Conseguenza |
|-------|--------|-------------|
| Ritiro | 7 giorni | Penale 1% al buyer, rimborso 99% |
| Spedizione seller | 3 giorni | Alert, poi possibile annullamento |
| Contestazione | 48h da consegna | Dopo: soldi sbloccati |
| Risposta disputa | 48h | Escalation ad admin |

### 3. Penale No-Show

**Scenario:** Buyer paga, non ritira, non risponde.

**Regola:**
```
SE (ordine.stato == PAID_UNCOLLECTED)
E (oggi > ordine.deadline)
E (buyer non ha aperto disputa)
ALLORA:
  - Trattenere 1% come penale
  - Penale va al seller (per il disturbo)
  - Rimborsare 99% al buyer
  - Prodotto torna disponibile
  - +1 Strike al buyer
```

**Strike System:**
- 1 strike: Warning email
- 2 strike: Limite 1 ordine attivo
- 3 strike: Ban temporaneo 30 giorni
- 5 strike: Ban permanente

### 4. Sistema Dispute

**Motivi validi:**
```go
type DisputeReason string

const (
    ItemNotReceived     DisputeReason = "ITEM_NOT_RECEIVED"
    ItemDamaged         DisputeReason = "ITEM_DAMAGED"
    ItemNotAsDescribed  DisputeReason = "ITEM_NOT_AS_DESCRIBED"
    SellerNoShow        DisputeReason = "SELLER_NO_SHOW"
    ScamAttempt         DisputeReason = "SCAM_ATTEMPT"
)
```

**Prove richieste:**
- Minimo 1 foto (max 5)
- Descrizione testuale (min 50 caratteri)
- Data/ora del problema

**Workflow disputa:**
```
OPEN → SELLER_RESPONSE → BUYER_REVIEW → [RESOLVED | ADMIN_REVIEW]
```

**Risoluzioni possibili:**
| Decisione | Azione |
|-----------|--------|
| `REFUND_FULL` | 100% al buyer |
| `REFUND_PARTIAL` | X% al buyer, resto al seller |
| `PAYOUT_SELLER` | 100% al seller (buyer ha torto) |
| `SPLIT` | 50/50 o altra % |

### 5. Truffe Comuni e Contromisure

#### Truffa "Pacco Vuoto"
**Scenario:** Buyer dice che il pacco era vuoto/con un mattone.

**Contromisure:**
- Per ordini > 200€: Seller deve caricare video chiusura pacco
- Assicurazione corriere obbligatoria
- Storico feedback visibile

#### Truffa "Non l'ho mai ricevuto" (Friendly Fraud)
**Scenario:** Buyer riceve ma nega.

**Contromisure:**
- Tracking obbligatorio
- Firma alla consegna per ordini > 100€
- Storico dispute del buyer visibile all'admin

#### Truffa "Merce Scadente"
**Scenario:** Merce diversa da foto/descrizione.

**Contromisure:**
- Foto multiple obbligatorie (min 3)
- Descrizione dettagliata obbligatoria
- Se seller ha > 3 dispute "NOT_AS_DESCRIBED": review manuale nuovi annunci

#### Elusione Commissioni
**Scenario:** "Facciamo fuori, ti do i soldi a mano"

**Contromisure:**
- AI Moderation chat (vedi sotto)
- Nessuna tutela se pagano fuori
- Warning automatico se tentano

---

## Moderazione Chat

### Architettura 3 Livelli

```
┌─────────────────────────────────────────────────────────────┐
│                    MESSAGGIO IN ARRIVO                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  LIVELLO 1: RegEx (Gratis, <1ms)                            │
│  - Numeri telefono: /(\+39)?[\s]?3[0-9]{8,9}/               │
│  - Email: /[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/  │
│  - Keywords: whatsapp, telegram, contanti, bonifico, iban   │
│  AZIONE: Blocca messaggio, mostra warning                   │
└─────────────────────────────────────────────────────────────┘
                              │
                        (se passa)
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  LIVELLO 2: OpenAI Moderation (Gratis)                      │
│  Endpoint: https://api.openai.com/v1/moderations            │
│  Controlla: hate, violence, self-harm, sexual, harassment   │
│  AZIONE: Se flagged → Blocca + Alert admin                  │
└─────────────────────────────────────────────────────────────┘
                              │
                    (se utente sospetto)
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  LIVELLO 3: AI Mini (Gemini Flash / GPT-4o-mini)            │
│  Costo: ~0.01€ / 1000 messaggi                              │
│  Attivazione: Solo se utente < 5 ordini O già flaggato      │
│  Prompt: "Questo messaggio tenta di eludere pagamenti?"     │
│  AZIONE: Se sì → Shadow ban + Alert admin                   │
└─────────────────────────────────────────────────────────────┘
```

### Shadow Banning

**Cos'è:** L'utente vede il messaggio come "inviato", ma il destinatario non lo riceve.

**Quando si attiva:**
- 3 tentativi di elusione in 1 ora
- Messaggio flaggato come "SCAM" dal Livello 3

**Perché funziona:** Lo spammer pensa che l'app sia morta o l'utente maleducato. Si stufa e se ne va. Nessun intervento manuale richiesto.

### Keywords Bloccate

```go
var blockedKeywords = []string{
    // Contatti
    "whatsapp", "telegram", "signal", "messenger",
    "chiamami", "scrivimi", "il mio numero",

    // Pagamenti esterni
    "contanti", "cash", "bonifico", "iban",
    "paypal", "satispay", "postepay",
    "fuori app", "fuori piattaforma",

    // Sospetti
    "evitare commissioni", "risparmiare commissioni",
    "direttamente", "di persona",
}
```

---

## Sistema Feedback

### Struttura Rating

```
┌────────────────────────────────────────┐
│           FEEDBACK ORDINE              │
├────────────────────────────────────────┤
│ order_id: UUID                         │
│ from_user: UUID                        │
│ to_user: UUID                          │
│ rating: 1-5 stelle                     │
│ comment: text (opzionale)              │
│ is_anonymous: boolean                  │
│ created_at: timestamp                  │
└────────────────────────────────────────┘
```

### Regole Feedback

- Feedback possibile solo dopo ordine COMPLETED
- Window: 14 giorni dalla conclusione
- Buyer valuta Seller (e viceversa)
- Commenti moderati (Livello 1 + 2)
- Non modificabile dopo 24h

### Impatto Rating

| Rating Medio | Conseguenza |
|--------------|-------------|
| > 4.5 | Badge "Top Seller" |
| 4.0 - 4.5 | Normale |
| 3.0 - 4.0 | Warning, meno visibilità |
| < 3.0 | Review manuale, possibile sospensione |
| < 2.0 (5+ ordini) | Sospensione automatica |

---

## Limiti e Soglie

### Limiti Utente

| Tipo | Nuovo Utente | Utente Verificato |
|------|--------------|-------------------|
| Ordini attivi | 3 | 10 |
| Regali attivi | 1 | 3 |
| Messaggi/ora | 20 | 100 |
| Prodotti (seller) | 10 | Illimitati |
| Valore singolo ordine | 500€ | 5.000€ |

### Soglie Sicurezza

| Soglia | Azione |
|--------|--------|
| Ordine > 200€ | Video pacco richiesto |
| Ordine > 500€ (nuovo utente) | Verifica manuale |
| Seller nuovo, primo payout | Review manuale |
| > 3 dispute in 30 giorni | Sospensione temporanea |
| Chargeback Stripe | Ban immediato + indagine |

---

## Compliance e Fiscale

### Fatturazione Commissioni

**Flusso mensile:**
1. Fine mese: Sistema calcola commissioni per seller
2. Genera PDF fattura (GecoGreen → Seller)
3. Invia email con allegato
4. Seller scarica da dashboard

**Dati fattura:**
- Intestatario: Dati azienda seller
- Periodo: Mese di riferimento
- Dettaglio: Lista ordini, commissione per ordine
- Totale: Somma commissioni + IVA 22%

### Record Keeping

Conservare per 10 anni:
- Tutte le transazioni
- Fatture emesse/ricevute
- Log modifiche dati fiscali
- Comunicazioni ufficiali

---

*Documento creato: Dicembre 2024*
*Version: v1.0*
