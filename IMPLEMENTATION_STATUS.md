# GecoGreen - Implementation Status

**Data aggiornamento**: 2025-12-15

## ✅ Completato - Frontend

### 1. Unità di Misura Quantità
- [x] Tipi TypeScript: `QuantityUnit` (PIECE, KG, G, L, ML, CUSTOM)
- [x] Form inserimento prodotto: select unità + campo custom
- [x] Validazione: campo custom obbligatorio se CUSTOM selezionato

### 2. Asta Olandese (Dutch Auction)
- [x] Form completo con toggle on/off
- [x] Campi: prezzo iniziale, riduzione, intervallo ore, prezzo minimo
- [x] Preview dinamica dell'asta
- [x] Solo per listing type = SALE

### 3. Digital Freight Forwarders
- [x] Opzione aggiunta al menu spedizioni
- [x] Messaggio "Coming Soon" quando selezionata
- [x] Costo spedizione nascosto (sarà gestito dal forwarder)

### 4. Separazione Foto Prodotto/Scadenza
- [x] Due campi upload separati
- [x] Endpoint API separati
- [x] Form chiaro e intuitivo

---

## ⏳ Da Completare - Backend

### 1. Database Migrations
**File pronti**: `database/001_add_quantity_unit.sql`, `database/002_add_billing_info.sql`

**Istruzioni applicazione**:
```bash
# Via Docker
docker cp database/001_add_quantity_unit.sql <postgres_container>:/tmp/
docker exec -it <postgres_container> psql -U gecogreen -d gecogreen -f /tmp/001_add_quantity_unit.sql

docker cp database/002_add_billing_info.sql <postgres_container>:/tmp/
docker exec -it <postgres_container> psql -U gecogreen -d gecogreen -f /tmp/002_add_billing_info.sql
```

**Oppure via Coolify**:
1. Accedi al database PostgreSQL su Coolify
2. Copia il contenuto di `001_add_quantity_unit.sql` ed esegui
3. Copia il contenuto di `002_add_billing_info.sql` ed esegui

### 2. Modelli Go - Aggiornati ✅
- [x] `product.go`: aggiunti QuantityUnit, QuantityUnitCustom, ShippingDigitalForwarders
- [x] `user.go`: aggiunti FiscalCode, SDICode, PECEmail, EUVatID, BillingAddress, ecc.
- [x] UpdateProfileRequest: aggiunto AccountType per permettere cambio PRIVATE ↔ BUSINESS

### 3. Repository & Handlers - ⚠️ DA AGGIORNARE
Dopo aver applicato le migrations, aggiornare:

**`product_repository.go`**:
- [ ] Aggiungere `quantity_unit`, `quantity_unit_custom` alle query INSERT
- [ ] Aggiornare query SELECT per includerli
- [ ] Aggiornare query UPDATE

**`profile_handler.go`** / **`user_repository.go`**:
- [ ] Aggiungere campi billing_info alle query UPDATE profilo
- [ ] Validare cambio AccountType (PRIVATE → BUSINESS richiede P.IVA)
- [ ] Aggiungere warning quando si passa da BUSINESS a PRIVATE

### 4. Job Automatico Asta Olandese - 📋 Roadmap
- [ ] Cron job che gira ogni ora
- [ ] Controlla prodotti con `is_dutch_auction = true`
- [ ] Calcola ore passate da `dutch_started_at`
- [ ] Riduce prezzo se necessario
- [ ] Ferma quando raggiunge `dutch_min_price`
- [ ] Notifica venditore quando prezzo scende

---

## 🎯 Prossimi Step

### Immediato (Prima del deploy)
1. ✅ Applicare migrations al database produzione
2. ⚠️ Aggiornare repository Go per nuovi campi
3. ⚠️ Testare form inserimento prodotto
4. ⚠️ Testare cambio account type nel profilo

### Short-term (Q1 2026)
1. Implementare job asta olandese
2. UI per gestire cambio account type nel profilo
3. Sezione "Dati Fatturazione" nel profilo utente
4. Digital Freight Forwarders integration

---

## 📝 Note Tecniche

### Unità di Misura
```typescript
// Frontend
quantityUnit: 'KG' | 'G' | 'L' | 'ML' | 'PIECE' | 'CUSTOM'
quantityUnitCustom?: string // Se CUSTOM, es. "confezioni"

// Backend (Go)
QuantityUnit string // enum quantity_unit
QuantityUnitCustom *string
```

### Cambio Account Type
```go
// Consentito nel UpdateProfileRequest
AccountType *AccountType `json:"account_type,omitempty"`

// Validazione richiesta:
// PRIVATE → BUSINESS: richiede business_name + vat_number
// BUSINESS → PRIVATE: warning che dati aziendali rimangono ma nascosti
```

### Asta Olandese
```go
// Campi già presenti nel DB
IsDutchAuction bool
DutchStartPrice *float64
DutchDecreaseAmount *float64
DutchDecreaseHours *int
DutchMinPrice *float64
DutchStartedAt *time.Time

// Prezzo corrente calcolato con GetCurrentPrice()
func (p *Product) GetCurrentPrice() float64 {
    if !p.IsDutchAuction { return p.Price }
    hoursPassed := time.Since(*p.DutchStartedAt).Hours()
    intervals := int(hoursPassed) / (*p.DutchDecreaseHours)
    currentPrice := *p.DutchStartPrice - (float64(intervals) * (*p.DutchDecreaseAmount))
    if currentPrice < *p.DutchMinPrice {
        return *p.DutchMinPrice
    }
    return currentPrice
}
```

---

## 🐛 Issues Known

Nessuno al momento. Testare dopo deploy.

---

## ✅ Ready per Deploy?

**Frontend**: ✅ SI
**Backend Models**: ✅ SI
**Database**: ⏳ NO - Migrations da applicare
**Handlers/Repository**: ⏳ NO - Da aggiornare dopo migrations

**Raccomandazione**: Applicare migrations e aggiornare handlers prima del deploy produzione.
