# GecoGreen - Sistema Impatto Ambientale

## Overview

GecoGreen non vende "crediti di carbonio" certificati (troppo costoso e burocratico).
Invece, implementiamo un **Sistema di Impatto** che include:

1. **Eco-Scontrino** - Certificato visuale post-acquisto
2. **EcoCredits** - Valuta virtuale gamificata
3. **Report CSR** - Documento per aziende (bilancio sostenibilità)
4. **Community Impact** - Contatore alberi piantati

---

## 1. Eco-Scontrino (Impact Certificate)

### Cos'è
Un PDF/immagine generato automaticamente dopo ogni ordine completato.
Condivisibile su Instagram/LinkedIn per massimizzare la viralità.

### Contenuto del Certificato

```
┌────────────────────────────────────────────────────────────┐
│                                                            │
│     🌱 CERTIFICATO DI RISPARMIO AMBIENTALE 🌱              │
│                                                            │
│  ─────────────────────────────────────────────────────     │
│                                                            │
│  Grazie, Mario Rossi!                                      │
│                                                            │
│  Hai salvato dalla discarica:                              │
│  ┌────────────────────────────────────────────────┐        │
│  │  🪑  1x Sedia da ufficio ergonomica            │        │
│  │      Categoria: Arredamento                     │        │
│  └────────────────────────────────────────────────┘        │
│                                                            │
│  IL TUO IMPATTO:                                           │
│                                                            │
│     🏭  25 kg di CO₂ risparmiati                           │
│     💧  1.500 litri d'acqua risparmiati                    │
│     🗑️  15 kg di rifiuti evitati                           │
│                                                            │
│  EQUIVALE A:                                               │
│     🚗  100 km in auto non percorsi                        │
│     🌳  1 albero che assorbe CO₂ per 1 anno                │
│     💡  50 ore di lampadina LED spenta                     │
│                                                            │
│  ─────────────────────────────────────────────────────     │
│                                                            │
│  +250 EcoCredits guadagnati!                               │
│  Livello attuale: 🌿 Eco-Warrior (1.250 punti)             │
│                                                            │
│  ─────────────────────────────────────────────────────     │
│                                                            │
│  📅 12 Dicembre 2024                                       │
│  🔗 gecogreen.com/impact/ABC123                             │
│                                                            │
│  [Logo GecoGreen]                                          │
│  "Dai valore a ciò che resta"                              │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### Calcolo Impatto

**Sorgente dati:** Valori medi per categoria (non serve precisione scientifica, serve ordine di grandezza).

| Categoria | CO₂ (kg) | Acqua (L) | Rifiuti (kg) |
|-----------|----------|-----------|--------------|
| Alimentari Freschi | 2 | 500 | 1 |
| Alimentari Confezionati | 1 | 200 | 0.5 |
| Bevande | 0.5 | 100 | 0.3 |
| Frutta e Verdura | 1 | 300 | 0.5 |
| Surgelati | 3 | 400 | 1 |
| Cosmetici | 5 | 1000 | 0.5 |
| Cura Persona | 3 | 800 | 0.3 |
| Detergenza Casa | 4 | 600 | 1 |
| Pet Food | 2 | 400 | 0.5 |
| Giardinaggio | 5 | 200 | 2 |
| Materiali Tecnici | 15 | 500 | 5 |
| Automotive | 30 | 2000 | 10 |
| Sicurezza/DPI | 10 | 800 | 3 |
| HORECA | 3 | 500 | 1 |

**Equivalenze standard:**
- 1 kg CO₂ = 4 km in auto
- 1 kg CO₂ = 2 ore di lampadina LED
- 25 kg CO₂ = 1 albero/anno di assorbimento

### Formati Output

1. **PDF** - Per download/email
2. **PNG Square** - Per Instagram (1080x1080)
3. **PNG Story** - Per Instagram Stories (1080x1920)
4. **Link condivisibile** - Pagina web pubblica con dati

---

## 2. EcoCredits (Valuta Virtuale)

### Come si guadagnano

| Azione | EcoCredits |
|--------|------------|
| **Acquisto** | 1€ speso = 1 punto |
| **Vendita completata** | 1 kg CO₂ risparmiata = 10 punti |
| **Regalo donato** | 50 punti fissi |
| **Ritiro a mano** (no spedizione) | +50 punti bonus |
| **Prima recensione** | 20 punti |
| **Invita un amico** (che completa ordine) | 100 punti |
| **Profilo completo** | 30 punti |

### Livelli Utente (Gamification)

| Livello | Punti | Badge | Benefici |
|---------|-------|-------|----------|
| 🌱 Germoglio | 0-99 | Starter | - |
| 🌿 Eco-Curious | 100-499 | Bronze | Badge visibile |
| 🌳 Eco-Warrior | 500-1999 | Silver | Badge + 1 Boost gratuito/mese |
| 🌲 Eco-Champion | 2000-4999 | Gold | Badge + 2 Boost/mese + priorità ricerca |
| 🏆 Eco-Legend | 5000+ | Platinum | Badge + 3 Boost/mese + priorità supporto |

### Come si spendono

| Reward | Costo | Costo reale | Note |
|--------|-------|-------------|------|
| Boost visibilità prodotto (24h) | 200 punti | 0€ | Algoritmo |
| Donazione a Tree-Nation (1 albero) | 300 punti | ~1€ | Marketing |
| Badge personalizzato | 1000 punti | 0€ | Fidelizzazione |

**Nota economica:** Tutti i premi sono stati validati per sostenibilità economica.
Per guadagnare 300 punti servono ~300€ di acquisti → 30€ commissioni per la piattaforma.
Costo albero: 1€. Margine garantito: 29€.

### Regole Anti-Abuso

- I punti scadono dopo 24 mesi di inattività
- Max 1000 punti/giorno da acquisti
- Ordini annullati/rimborsati: punti revocati
- Account sospesi: punti congelati

---

## 3. Report CSR per Aziende

### Cos'è
Report PDF professionale per il **Bilancio di Sostenibilità** aziendale.
Generato automaticamente a fine anno (o su richiesta).

### Target
- Aziende che smaltiscono: PC, arredi ufficio, macchinari, eccedenze magazzino
- Invece di pagare per smaltire → vendono/regalano su GecoGreen → ottengono report

### Contenuto Report

```
┌────────────────────────────────────────────────────────────┐
│                                                            │
│  REPORT DI ECONOMIA CIRCOLARE 2024                         │
│  Azienda: TechCorp Srl                                     │
│                                                            │
│  ══════════════════════════════════════════════════════    │
│                                                            │
│  EXECUTIVE SUMMARY                                         │
│  Nel 2024, TechCorp ha reimmesso nel mercato 127 oggetti   │
│  evitando la produzione di nuovi beni e lo smaltimento     │
│  in discarica.                                             │
│                                                            │
│  ══════════════════════════════════════════════════════    │
│                                                            │
│  IMPATTO AMBIENTALE TOTALE                                 │
│                                                            │
│  🏭 CO₂ Risparmiata:        3.250 kg                       │
│  💧 Acqua Risparmiata:      85.000 litri                   │
│  🗑️ Rifiuti Evitati:        450 kg                         │
│                                                            │
│  ══════════════════════════════════════════════════════    │
│                                                            │
│  DETTAGLIO PER CATEGORIA                                   │
│                                                            │
│  │ Categoria      │ Qty │ CO₂ (kg) │ Valore €  │           │
│  │────────────────│─────│──────────│───────────│           │
│  │ Elettronica    │  45 │    2.250 │    4.500  │           │
│  │ Arredamento    │  32 │      960 │    2.100  │           │
│  │ Altro          │  50 │       40 │      800  │           │
│  │────────────────│─────│──────────│───────────│           │
│  │ TOTALE         │ 127 │    3.250 │    7.400  │           │
│                                                            │
│  ══════════════════════════════════════════════════════    │
│                                                            │
│  EQUIVALENZE                                               │
│                                                            │
│  La CO₂ risparmiata equivale a:                            │
│  • 13.000 km in automobile non percorsi                    │
│  • 130 alberi che assorbono CO₂ per un anno                │
│  • 6.500 ore di lampadina LED spente                       │
│                                                            │
│  ══════════════════════════════════════════════════════    │
│                                                            │
│  CERTIFICAZIONE                                            │
│                                                            │
│  Questo report è stato generato automaticamente da         │
│  GecoGreen sulla base delle transazioni registrate.        │
│                                                            │
│  ID Report: CSR-2024-TECHCORP-ABC123                       │
│  Data generazione: 31/12/2024                              │
│  Verificabile su: gecogreen.com/csr/ABC123                  │
│                                                            │
│  ══════════════════════════════════════════════════════    │
│                                                            │
│  [Logo GecoGreen]        [QR Code Verifica]                │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### Pricing Report CSR

| Piano | Prezzo | Include |
|-------|--------|---------|
| **Base** | Gratis | Report annuale automatico |
| **Pro** | 99€/anno | Report trimestrale + logo aziendale + export Excel |
| **Enterprise** | 299€/anno | Report mensile + API dati + consulenza |

---

## 4. Community Impact (Contatore Globale)

### Homepage Widget

```
┌────────────────────────────────────────────────────────────┐
│                                                            │
│  🌍 IMPATTO DELLA COMMUNITY GECOGREEN                      │
│                                                            │
│     🏭  125.430 kg    💧  3.2M litri    🌳  520           │
│        CO₂ salvata       Acqua salvata     Alberi         │
│                                                            │
│  [Scopri come contribuisci anche tu →]                     │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### Integrazione Tree-Nation

**Come funziona:**
- Ogni 100 ordini completati → GecoGreen pianta 1 albero vero
- Integrazione API Tree-Nation (o Ecologi)
- Certificato albero visibile nella pagina community

**Costo:** ~1€/albero = ~1% delle commissioni

---

## 5. Implementazione Tecnica

### Nuove Colonne Database

```sql
-- Categorie: valori impatto stimati
ALTER TABLE categories ADD COLUMN estimated_co2_kg DECIMAL(10,2) DEFAULT 0;
ALTER TABLE categories ADD COLUMN estimated_water_l DECIMAL(10,2) DEFAULT 0;
ALTER TABLE categories ADD COLUMN estimated_waste_kg DECIMAL(10,2) DEFAULT 0;

-- Utenti: statistiche aggregate e gamification
ALTER TABLE users ADD COLUMN total_co2_saved DECIMAL(10,2) DEFAULT 0;
ALTER TABLE users ADD COLUMN total_water_saved DECIMAL(10,2) DEFAULT 0;
ALTER TABLE users ADD COLUMN eco_credits INT DEFAULT 0;
ALTER TABLE users ADD COLUMN eco_level VARCHAR(50) DEFAULT 'Germoglio';

-- Ordini: impatto specifico ordine
ALTER TABLE orders ADD COLUMN co2_saved DECIMAL(10,2);
ALTER TABLE orders ADD COLUMN water_saved DECIMAL(10,2);
ALTER TABLE orders ADD COLUMN eco_credits_earned INT;
ALTER TABLE orders ADD COLUMN impact_certificate_url VARCHAR(500);
```

### Nuova Tabella: impact_logs

```sql
CREATE TABLE impact_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    order_id UUID REFERENCES orders(id),

    -- Tipo azione
    action_type VARCHAR(50) NOT NULL,  -- 'PURCHASE', 'SALE', 'GIFT', 'BONUS'

    -- Impatto
    co2_saved DECIMAL(10,2) DEFAULT 0,
    water_saved DECIMAL(10,2) DEFAULT 0,

    -- Punti
    eco_credits_earned INT DEFAULT 0,
    eco_credits_spent INT DEFAULT 0,

    -- Descrizione
    description TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_impact_logs_user ON impact_logs(user_id);
CREATE INDEX idx_impact_logs_order ON impact_logs(order_id);
```

### Nuova Tabella: csr_reports (per aziende)

```sql
CREATE TABLE csr_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES users(id),

    -- Periodo
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,

    -- Statistiche aggregate
    total_items INT DEFAULT 0,
    total_co2_kg DECIMAL(10,2) DEFAULT 0,
    total_water_l DECIMAL(10,2) DEFAULT 0,
    total_waste_kg DECIMAL(10,2) DEFAULT 0,
    total_value DECIMAL(10,2) DEFAULT 0,

    -- Dettaglio per categoria (JSON)
    category_breakdown JSONB,

    -- File generato
    pdf_url VARCHAR(500),
    report_code VARCHAR(50) UNIQUE,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### API Endpoints

```
GET  /api/impact/my-stats           # Statistiche utente corrente
GET  /api/impact/certificate/:id    # Singolo certificato
GET  /api/impact/leaderboard        # Top utenti per CO2 salvata
GET  /api/impact/community          # Stats globali community
POST /api/impact/redeem             # Spendi EcoCredits

# Solo Seller
GET  /api/seller/csr-report         # Genera report CSR
GET  /api/seller/csr-reports        # Lista report passati
```

---

## 6. Roadmap Implementazione

### MVP (Fase 1)
- [x] Campi impatto nel DB
- [ ] Calcolo automatico CO2 per ordine
- [ ] EcoCredits base (guadagno)
- [ ] Widget stats utente in dashboard

### Fase 2
- [ ] Generazione PDF Eco-Scontrino
- [ ] Condivisione social (PNG)
- [ ] Livelli utente
- [ ] Spending EcoCredits

### Fase 3
- [ ] Report CSR per aziende
- [ ] Integrazione Tree-Nation
- [ ] Leaderboard pubblica
- [ ] API partner

---

## 7. Note Marketing

### Claim da usare
- "Ogni acquisto salva il pianeta"
- "Trasforma lo spreco in impatto"
- "Il tuo shopping ha un'anima verde"

### Hashtag
- #GecoGreenImpact
- #ZeroSprecoZeroEmissioni
- #EcoCredits

### Press Kit
- Infografica impatto community
- Case study aziende (quando disponibili)
- Dati aggregati per giornalisti

---

*Documento creato: Dicembre 2024*
*Version: v1.0*
