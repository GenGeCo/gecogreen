# GecoGreen - Flussi Utente e Pagine

## Ruoli Utente

| Ruolo | Descrizione | Permessi |
|-------|-------------|----------|
| **GUEST** | Visitatore non registrato | Naviga catalogo, vede prezzi, non può comprare |
| **BUYER** | Cliente registrato | Compra, chat, gestisce ordini |
| **SELLER** | Commerciante verificato | Vende, gestisce inventario, riceve pagamenti |
| **ADMIN** | Amministratore piattaforma | Tutto + gestione utenti e dispute |

> Un utente può essere sia BUYER che SELLER contemporaneamente

---

## Mappa Pagine

### Pagine Pubbliche (GUEST)

```
/                       → Homepage (hero + prodotti in evidenza)
/catalogo               → Lista prodotti con filtri
/catalogo/[categoria]   → Prodotti per categoria
/prodotto/[id]          → Scheda prodotto dettagliata
/cerca                  → Ricerca avanzata
/come-funziona          → Spiegazione piattaforma
/diventa-venditore      → Landing per seller
/auth/login             → Login
/auth/registrati        → Registrazione
/auth/reset-password    → Reset password
```

### Dashboard BUYER

```
/dashboard                      → Overview (ordini recenti, messaggi)
/dashboard/ordini               → Lista ordini
/dashboard/ordini/[id]          → Dettaglio ordine + chat + QR
/dashboard/messaggi             → Inbox conversazioni
/dashboard/messaggi/[orderId]   → Chat specifica
/dashboard/profilo              → Dati personali
/dashboard/preferiti            → Prodotti salvati
```

### Dashboard SELLER

```
/seller                         → Overview (vendite, saldo, notifiche)
/seller/prodotti                → Lista miei prodotti
/seller/prodotti/nuovo          → Crea prodotto (seleziona sede)
/seller/prodotti/[id]/modifica  → Modifica prodotto
/seller/ordini                  → Ordini ricevuti (Kanban o lista)
/seller/ordini/[id]             → Dettaglio ordine + chat
/seller/ordini/[id]/qr          → Scanner QR per conferma ritiro
/seller/sedi                    → Lista punti vendita/ritiro
/seller/sedi/nuova              → Aggiungi nuova sede
/seller/sedi/[id]               → Modifica sede (orari, indirizzo)
/seller/wallet                  → Saldo, movimenti, payout
/seller/fatture                 → Fatture commissioni ricevute
/seller/impostazioni            → Dati azienda
```

### Gestione Multi-Sede (Seller)

**Caso d'uso:** Catena supermercati, azienda con più magazzini, ristorante con più sedi.

**Flusso creazione sede:**
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ /seller/sedi│───▶│ + Nuova Sede│───▶│ Form:       │
│             │    │             │    │ - Nome      │
│             │    │             │    │ - Indirizzo │
│             │    │             │    │ - Orari     │
│             │    │             │    │ - Telefono  │
└─────────────┘    └─────────────┘    └─────────────┘
```

**Flusso creazione prodotto (con sede):**
```
┌─────────────────────────────────────────────────┐
│         NUOVO PRODOTTO                          │
├─────────────────────────────────────────────────┤
│ Sede di ritiro: [▼ Seleziona sede]              │
│   ○ Milano Centro - Via Roma 1                  │
│   ○ Milano Nord - Via Zara 50                   │
│   ○ Magazzino Lodi                              │
├─────────────────────────────────────────────────┤
│ Titolo: [________________________]              │
│ Categoria: [▼ Alimentari Freschi]               │
│ ...                                             │
└─────────────────────────────────────────────────┘
```

**Vista Buyer (scheda prodotto):**
```
┌─────────────────────────────────────────────────┐
│ [Immagine Prodotto]                             │
├─────────────────────────────────────────────────┤
│ Yogurt Bio - Scadenza 15/12                     │
│ 2,50€  ███████ 5,00€                            │
├─────────────────────────────────────────────────┤
│ 📍 Ritiro presso:                               │
│    Supermercato XYZ - Sede Milano Centro        │
│    Via Roma 1, 20100 Milano                     │
│    Orari: Lun-Sab 9:00-19:00                    │
│    [Vedi su mappa]                              │
└─────────────────────────────────────────────────┘
```

### Dashboard ADMIN

```
/admin                          → Overview (stats, alert)
/admin/utenti                   → Lista utenti (cerca, filtra)
/admin/utenti/[id]              → Dettaglio utente (modifica, ban)
/admin/venditori                → Richieste verifica seller
/admin/dispute                  → Contestazioni aperte
/admin/dispute/[id]             → Gestione singola disputa
/admin/ordini                   → Tutti gli ordini (debug)
/admin/finanze                  → Revenue, commissioni, payout
/admin/log                      → Audit log azioni
```

---

## Flussi Dettagliati

### Flusso 1: Registrazione Buyer

**Opzione A: Social Login (Consigliata)**
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Homepage  │───▶│ "Continua   │───▶│  Dashboard  │
│   (CTA)     │    │ con Google" │    │  Buyer      │
└─────────────┘    └─────────────┘    └─────────────┘
```
- 1 click, nessun form
- Email già verificata
- Avatar importato automaticamente

**Opzione B: Email tradizionale**
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Homepage  │───▶│ Registrati  │───▶│ Verifica    │───▶│  Dashboard  │
│   (CTA)     │    │ (Form)      │    │ Email       │    │  Buyer      │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

**Schermata Login/Registrazione:**
```
┌────────────────────────────────────┐
│                                    │
│   [G] Continua con Google          │
│   [🍎] Continua con Apple          │
│                                    │
│   ──────── oppure ────────         │
│                                    │
│   Email: [________________]        │
│   Password: [______________]       │
│                                    │
│   [    ACCEDI    ]                 │
│                                    │
│   Non hai un account? Registrati   │
│                                    │
└────────────────────────────────────┘
```

**Campi form (se email):**
- Nome, Cognome
- Email
- Password (min 8 char, 1 numero)
- Città (per ricerche locali)
- Checkbox T&C

**Post-registrazione:**
- Email con link verifica (solo se registrazione email)
- Redirect a dashboard
- Popup "Completa profilo" (chiedi città se mancante)

---

### Flusso 2: Registrazione Seller

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Landing     │───▶│ Form Base   │───▶│ Form        │───▶│ Stripe      │───▶│ Attesa      │
│ Venditore   │    │ (come buyer)│    │ Aziendale   │    │ Onboarding  │    │ Verifica    │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                                                                                   │
                                                                                   ▼
                                                                           ┌─────────────┐
                                                                           │ Dashboard   │
                                                                           │ Seller      │
                                                                           └─────────────┘
```

**Campi aggiuntivi seller:**
- Ragione sociale
- P.IVA
- Indirizzo sede
- Telefono
- Categoria merceologica
- Descrizione attività

**Stripe Onboarding:**
- Redirect a Stripe Connect
- Inserimento IBAN
- Verifica identità (per payout)

**Verifica Admin:**
- Admin riceve notifica
- Controlla dati aziendali
- Approva o richiede documenti

---

### Flusso 2b: Inserimento Prodotto (Seller)

**Form nuovo prodotto:**
```
┌─────────────────────────────────────────────────────────────┐
│              NUOVO PRODOTTO                                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  INFORMAZIONI BASE                                           │
│  ─────────────────                                           │
│  Titolo: [_________________________________]                 │
│  Categoria: [▼ Alimentari Freschi]                           │
│  Descrizione: [                                  ]           │
│               [                                  ]           │
│                                                              │
│  FOTO PRODOTTO (minimo 3, massimo 10)                        │
│  ─────────────────                                           │
│  [+ Aggiungi foto]                                           │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐                             │
│  │ 📷1 │ │ 📷2 │ │ 📷3 │ │ ... │                             │
│  └─────┘ └─────┘ └─────┘ └─────┘                             │
│                                                              │
│  ⚠️ FOTO SCADENZA (obbligatoria per alimentari)              │
│  ─────────────────                                           │
│  [📷 Carica foto etichetta con data scadenza]                │
│  (La foto deve mostrare chiaramente la data di scadenza)     │
│                                                              │
│  PREZZI                                                      │
│  ─────────────────                                           │
│  ○ Prezzo fisso                                              │
│     Prezzo attuale: [_____] €                                │
│     Prezzo originale: [_____] €  (opzionale, per sconto)     │
│                                                              │
│  ○ Prezzo decrescente (Dutch Auction) 📉                     │
│     Prezzo iniziale: [100__] €                               │
│     Diminuzione: [1____] € ogni [24] ore                     │
│     Prezzo minimo: [50___] €  (si ferma qui)                 │
│     [i] Il prezzo scende automaticamente ogni giorno!        │
│                                                              │
│  SCADENZA                                                    │
│  ─────────────────                                           │
│  Data scadenza: [__/__/____]  (da etichetta)                 │
│                                                              │
│  CONSEGNA                                                    │
│  ─────────────────                                           │
│  Sede di ritiro: [▼ Milano Centro - Via Roma 1]              │
│                                                              │
│  Metodo:                                                     │
│  ☑ Ritiro in sede                                            │
│  ☐ Spedizione (costo: [___] €)                               │
│  ☐ Trasporto a carico acquirente                             │
│                                                              │
│  QUANTITÀ                                                    │
│  ─────────────────                                           │
│  Disponibili: [10__] unità                                   │
│  Peso unitario: [0.5_] kg  (per calcolo spedizione)          │
│                                                              │
│  [SALVA BOZZA]  [PUBBLICA]                                   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Dutch Auction (Prezzo Decrescente):**

Il prezzo scende automaticamente nel tempo per incentivare l'acquisto rapido.

```
Esempio: Yogurt in scadenza

Giorno 1: 10€ (prezzo iniziale)
     ▼ -1€ dopo 24h
Giorno 2: 9€
     ▼ -1€ dopo 24h
Giorno 3: 8€
     ...
Giorno 6: 5€ (prezzo minimo raggiunto)
Giorno 7: 5€ (resta così fino a vendita)
```

**Vista catalogo con Dutch Auction:**
```
┌─────────────────────────────┐
│ [Immagine]                  │
│ ────────────────────────────│
│ Yogurt Bio 6 pack           │
│ 7,00€  📉 -1€/giorno        │
│ Min: 4€ │ Scade tra 5gg     │
│ 🕐 Prossimo ribasso: 18:32  │
│ [⭐ 4.8] [❤️]               │
└─────────────────────────────┘
```

**Regole foto:**
- Minimo 3 foto del prodotto
- Massimo 10 foto
- Per alimentari: foto scadenza OBBLIGATORIA
- Formati: JPG, PNG, WebP
- Dimensione max: 5MB per foto
- Resize automatico a 1200px lato lungo

---

### Flusso 3: Acquisto con Ritiro

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Catalogo │──▶│ Scheda   │──▶│ Carrello │──▶│ Checkout │──▶│ Pagamento│
│          │   │ Prodotto │   │          │   │          │   │ Stripe   │
└──────────┘   └──────────┘   └──────────┘   └──────────┘   └──────────┘
                                                                  │
     ┌────────────────────────────────────────────────────────────┘
     ▼
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Conferma │──▶│ Sblocco  │──▶│ Ritiro   │──▶│ Conferma │
│ Ordine   │   │ Indirizzo│   │ + QR     │   │ Ricezione│
└──────────┘   └──────────┘   └──────────┘   └──────────┘
```

**Stati ordine ritiro:**
1. `CREATED` - Nel carrello
2. `AWAITING_PAYMENT` - Checkout iniziato
3. `PAID` - Pagato, indirizzo sbloccato
4. `READY_PICKUP` - Seller ha preparato (opzionale)
5. `COMPLETED` - QR scansionato, merce ritirata
6. `DISPUTED` - Contestazione aperta

**Dopo il pagamento:**
- Buyer vede: indirizzo completo, orari, telefono seller
- Buyer riceve: QR code univoco nell'app
- Timer: 7 giorni per ritirare

**Al ritiro:**
- Buyer mostra QR
- Seller scansiona con camera
- Sistema conferma automaticamente
- Soldi sbloccati al seller (dopo 48h buffer)

---

### Flusso 4: Acquisto con Spedizione

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Checkout │──▶│ Scelta   │──▶│ Pagamento│──▶│ Seller   │
│          │   │ Spediz.  │   │ Tot+Ship │   │ Spedisce │
└──────────┘   └──────────┘   └──────────┘   └──────────┘
                    │                              │
                    ▼                              ▼
            ┌──────────────┐              ┌──────────────┐
            │ A) Seller    │              │ Tracking     │
            │    spedisce  │              │ inserito     │
            │    (costo X) │              └──────────────┘
            ├──────────────┤                     │
            │ B) Buyer     │                     ▼
            │    organizza │              ┌──────────────┐
            │    trasporto │              │ Consegnato   │
            └──────────────┘              │ (48h buffer) │
                                          └──────────────┘
                                                 │
                                                 ▼
                                          ┌──────────────┐
                                          │ Soldi        │
                                          │ sbloccati    │
                                          └──────────────┘
```

**Opzioni spedizione:**
- **Seller Ships:** Seller indica costo, compreso nel totale
- **Buyer Arranges:** Buyer paga solo prodotto, organizza lui il ritiro/trasporto
- **Platform Managed (Futuro):** GecoGreen gestisce il corriere (vedi sotto)

**Post-spedizione:**
- Seller inserisce tracking
- Buyer segue tracking
- Alla consegna: 48h per contestare
- Poi: soldi al seller

---

### Flusso 4b: Spedizione Gestita GecoGreen (Futuro)

> ⚠️ **Funzionalità futura** - Sarà attivata dall'admin quando pronta

**Cos'è:**
GecoGreen gestisce direttamente la spedizione. Il buyer paga a noi, noi paghiamo il corriere.

**Vantaggi:**
- Tariffe negoziate con i corrieri
- Esperienza unificata per il buyer
- Meno lavoro per il seller (etichetta pre-generata)
- Tracking integrato nella piattaforma

**Flusso:**
```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Checkout │──▶│ Calcolo  │──▶│ Buyer    │──▶│ Etichetta│
│          │   │ Tariffa  │   │ Paga     │   │ Generata │
└──────────┘   └──────────┘   └──────────┘   └──────────┘
                    │                              │
                    ▼                              ▼
            ┌──────────────┐              ┌──────────────┐
            │ GecoGreen    │              │ Seller       │
            │ calcola:     │              │ riceve PDF   │
            │ base + kg    │              │ da stampare  │
            │ + markup 10% │              └──────────────┘
            └──────────────┘                     │
                                                 ▼
                                          ┌──────────────┐
                                          │ Corriere     │
                                          │ ritira       │
                                          │ dal seller   │
                                          └──────────────┘
```

**Pannello Admin:**
```
┌─────────────────────────────────────────────────────────────┐
│              GESTIONE CORRIERI                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ☑ Spedizioni GecoGreen: ATTIVE                              │
│                                                              │
│  CORRIERI ABILITATI                                          │
│  ─────────────────                                           │
│  │ Corriere │ Base  │ €/kg │ Max kg │ Stato    │             │
│  │──────────│───────│──────│────────│──────────│             │
│  │ BRT      │ 5,00€ │ 0,50 │ 30 kg  │ ✅ Attivo │             │
│  │ GLS      │ 4,50€ │ 0,60 │ 25 kg  │ ✅ Attivo │             │
│  │ DHL      │ 7,00€ │ 0,40 │ 50 kg  │ ❌ Disab. │             │
│  │ Poste    │ 4,00€ │ 0,80 │ 20 kg  │ ✅ Attivo │             │
│                                                              │
│  IMPOSTAZIONI                                                │
│  ─────────────────                                           │
│  Ordine minimo: [20___] €                                    │
│  Spedizione gratis sopra: [100__] €                          │
│  Markup piattaforma: [10___] %                               │
│  Assicurazione: [1____] % del valore                         │
│                                                              │
│  [SALVA IMPOSTAZIONI]                                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Vista Checkout (con opzione GecoGreen):**
```
┌─────────────────────────────────────────────────────────────┐
│              SCEGLI SPEDIZIONE                               │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ○ Ritiro in sede (gratuito)                                 │
│    📍 Via Roma 1, Milano                                     │
│                                                              │
│  ○ Spedizione Seller (8,00€)                                 │
│    🚚 Il venditore spedisce autonomamente                    │
│                                                              │
│  ● Spedizione GecoGreen (6,50€) ✨ CONSIGLIATA               │
│    🚚 BRT Express - 2-4 giorni lavorativi                    │
│    ✓ Tracking integrato                                      │
│    ✓ Assistenza GecoGreen                                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

### Flusso 5: Regalo (Gratis)

```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Catalogo │──▶│ Scheda   │──▶│ "Lo      │──▶│ Chat si  │
│ Regali   │   │ (0€)     │   │  Voglio" │   │ apre     │
└──────────┘   └──────────┘   └──────────┘   └──────────┘
                                                  │
                    ┌─────────────────────────────┘
                    ▼
             ┌────────────┐
             │ A) Ritiro  │───▶ Accordo in chat ───▶ QR ───▶ Fine
             ├────────────┤
             │ B) Spediz. │───▶ Buyer paga solo ───▶ Tracking ───▶ Fine
             │            │     spedizione
             └────────────┘
```

**Regole regalo:**
- Prezzo forzato a 0€
- Nessuna commissione
- Max 3 regali attivi per buyer
- Feedback obbligatorio

**Se spedito:**
- Buyer paga SOLO il costo spedizione
- Fee servizio 0.50€ (per coprire costi)

---

### Flusso 6: Contestazione (Disputa)

```
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│ Ordine       │──▶│ "Segnala     │──▶│ Form +       │
│ Completato   │   │  Problema"   │   │ Upload Foto  │
└──────────────┘   └──────────────┘   └──────────────┘
                                             │
                                             ▼
                                    ┌──────────────┐
                                    │ Seller       │
                                    │ Notificato   │
                                    └──────────────┘
                                             │
                   ┌─────────────────────────┼─────────────────────────┐
                   ▼                         ▼                         ▼
           ┌──────────────┐         ┌──────────────┐         ┌──────────────┐
           │ Seller       │         │ Seller       │         │ Nessuna      │
           │ Accetta      │         │ Propone      │         │ Risposta     │
           │ Rimborso     │         │ Sconto       │         │ (48h)        │
           └──────────────┘         └──────────────┘         └──────────────┘
                   │                         │                         │
                   ▼                         ▼                         ▼
           ┌──────────────┐         ┌──────────────┐         ┌──────────────┐
           │ Rimborso     │         │ Buyer        │         │ Admin        │
           │ Automatico   │         │ Accetta/Rif. │         │ Interviene   │
           └──────────────┘         └──────────────┘         └──────────────┘
```

**Motivi disputa:**
- `ITEM_NOT_RECEIVED` - Mai arrivato
- `ITEM_DAMAGED` - Danneggiato
- `ITEM_NOT_AS_DESCRIBED` - Diverso dalla descrizione
- `SELLER_NO_SHOW` - Venditore irreperibile

**Prove richieste:**
- Foto obbligatorie (max 5)
- Descrizione testuale
- Tutto dentro la chat ordine

**Risoluzione Admin:**
- Vede foto, chat, storico utenti
- Decide: rimborso totale/parziale o pagamento seller
- La decisione è finale

---

## Componenti UI Comuni

### Header
```
┌────────────────────────────────────────────────────────────┐
│ [Logo]  [Catalogo ▼]  [Cerca...]         [Carrello] [User] │
└────────────────────────────────────────────────────────────┘
```

### Scheda Prodotto (Card)
```
┌─────────────────────────────┐
│ [Immagine]                  │
│ ────────────────────────────│
│ Nome Prodotto               │
│ 12,50€  ██████ 25,00€       │
│ Scadenza: 3 giorni          │
│ Milano, Lambrate            │
│ [⭐ 4.8] [❤️]                │
└─────────────────────────────┘
```

### Stato Ordine (Badge)
```
🔵 PAGATO        → Blu
🟡 IN PREPARAZIONE → Giallo
🟢 COMPLETATO    → Verde
🔴 DISPUTA       → Rosso
⚫ ANNULLATO     → Grigio
```

### QR Code Ritiro e Identificazione

**Chi ha il QR?** Il BUYER, sul telefono nell'app.

**Come funziona:**
```
┌─────────────────────────────┐
│     ORDINE #LZ-2024-0001    │
│                             │
│      ┌─────────────┐        │
│      │ ▓▓▓▓▓▓▓▓▓▓▓ │        │
│      │ ▓▓▓▓▓▓▓▓▓▓▓ │        │
│      │ ▓▓▓ QR ▓▓▓▓ │        │
│      │ ▓▓▓▓▓▓▓▓▓▓▓ │        │
│      │ ▓▓▓▓▓▓▓▓▓▓▓ │        │
│      └─────────────┘        │
│                             │
│  Codice: ABC-123-XYZ        │
│  (se non funziona il QR)    │
│                             │
│  Mostra al venditore        │
│  Valido fino: 18/12/2024    │
│                             │
│  [Delega ritiro]            │
└─────────────────────────────┘
```

**Flusso al ritiro:**
```
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ Buyer    │──▶│ Seller   │──▶│ Sistema  │──▶│ Ordine   │
│ mostra   │   │ scansiona│   │ verifica │   │ COMPLETATO
│ QR       │   │ con app  │   │ firma    │   │          │
└──────────┘   └──────────┘   └──────────┘   └──────────┘
```

**Alternativa senza smartphone:**
- Il buyer può dettare il codice alfanumerico (ABC-123-XYZ)
- Il seller lo inserisce manualmente nell'app
- Stesso risultato

### Sistema Delega

**Caso d'uso:** "Non posso andare, mando mia moglie"

**Flusso delega:**
```
┌────────────────────────────────────────┐
│         DELEGA RITIRO                  │
├────────────────────────────────────────┤
│                                        │
│  Chi ritirerà al posto tuo?            │
│                                        │
│  Nome: [Mario_______________]          │
│  Cognome: [Rossi____________]          │
│                                        │
│  [  GENERA CODICE DELEGATO  ]          │
│                                        │
└────────────────────────────────────────┘

              ▼

┌────────────────────────────────────────┐
│         CODICE DELEGATO                │
├────────────────────────────────────────┤
│                                        │
│  Delegato: Mario Rossi                 │
│                                        │
│  Codice ritiro: DEL-987-ZZZ            │
│                                        │
│  Il delegato deve:                     │
│  1. Mostrare questo codice             │
│  2. Mostrare documento d'identità      │
│                                        │
│  [Condividi via WhatsApp]              │
│  [Copia codice]                        │
│                                        │
│  [Annulla delega]                      │
│                                        │
└────────────────────────────────────────┘
```

**Regole delega:**
- Un solo delegato per ordine
- Il delegato deve mostrare documento con nome corrispondente
- La delega può essere annullata fino al ritiro
- Il sistema logga chi ha effettivamente ritirato

**Vista Seller al ritiro (con delega):**
```
┌────────────────────────────────────────┐
│  VERIFICA RITIRO                       │
├────────────────────────────────────────┤
│                                        │
│  Ordine: #LZ-2024-0001                 │
│  Acquirente: Giuseppe Bianchi          │
│                                        │
│  ⚠️  RITIRO DELEGATO                   │
│  Delegato: Mario Rossi                 │
│                                        │
│  ✓ Verifica documento del delegato     │
│                                        │
│  [CONFERMA CONSEGNA]                   │
│                                        │
└────────────────────────────────────────┘
```

---

## Notifiche

### Email
| Evento | Destinatario | Oggetto |
|--------|--------------|---------|
| Registrazione | User | Benvenuto su GecoGreen! |
| Nuovo ordine | Seller | Hai ricevuto un ordine! |
| Pagamento ok | Buyer | Ordine confermato |
| Spedito | Buyer | Il tuo ordine è in viaggio |
| Promemoria ritiro | Buyer | Ricorda di ritirare entro X |
| Disputa aperta | Seller | Un cliente ha segnalato un problema |
| Disputa risolta | Entrambi | La contestazione è stata risolta |
| Payout | Seller | Hai ricevuto un pagamento |

### Push/In-App
- Nuovo messaggio in chat
- Ordine pronto al ritiro
- Timer ritiro in scadenza
- Disputa: risposta ricevuta

---

*Documento creato: Dicembre 2024*
*Version: v1.0*
