# GecoGreen - Sistema Impatto Ambientale

## Overview

GecoGreen non vende "crediti di carbonio" certificati (troppo costoso e burocratico).
Invece, implementiamo un **Sistema di Impatto** che include:

1. **Eco-Scontrino** - Certificato visuale post-acquisto
2. **EcoCredits** - Valuta virtuale gamificata
3. **Classifiche & Premi** - Leaderboard settimanali/mensili/annuali
4. **Report CSR** - Documento per aziende (bilancio sostenibilità)
5. **Community Impact** - Contatore alberi piantati
6. **YouTube Content** - Interviste ai top seller

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

#### Per BUYER (chi compra)

| Azione | EcoCredits |
|--------|------------|
| **Ogni € speso** | 1 punto |
| **Primo acquisto** | 100 punti (welcome bonus) |
| **Ritiro a mano** (no spedizione) | +50 punti bonus |
| **Acquisto "Last Chance"** (<24h scadenza) | +30% punti |
| **Prima recensione** | 20 punti |
| **5 ordini completati** | 200 punti |
| **Condivide Eco-Scontrino sui social** | 50 punti |
| **Invita un amico** (che completa ordine) | 100 punti |
| **Profilo completo** | 30 punti |

#### Per SELLER (chi vende)

| Azione | EcoCredits |
|--------|------------|
| **Vendita completata** | 1 kg CO₂ risparmiata = 10 punti |
| **Regalo donato** | 50 punti fissi |
| **Prima vendita** | 100 punti |
| **10 prodotti salvati** | 200 punti |
| **50 prodotti salvati** | 500 punti |

### Livelli Utente (Gamification)

| Livello | Punti | Badge | Benefici |
|---------|-------|-------|----------|
| 🌱 Germoglio | 0-99 | Starter | - |
| 🌿 Eco-Curious | 100-499 | Bronze | Badge visibile |
| 🌳 Eco-Warrior | 500-1999 | Silver | Badge + 1 Boost gratuito/mese |
| 🌲 Eco-Champion | 2000-4999 | Gold | Badge + 2 Boost/mese + priorità ricerca |
| 🏆 Eco-Legend | 5000+ | Platinum | Badge + 3 Boost/mese + priorità supporto |

### Come si spendono (SOLO COSE A COSTO ZERO)

| Reward | Costo | Costo Reale | Note |
|--------|-------|-------------|------|
| **Boost visibilità 24h** | 200 punti | 0€ | Algoritmo |
| **Pianta 1 albero** | 300 punti | ~1€ | Margine garantito (vedi sotto) |
| **Boost visibilità 7 giorni** | 500 punti | 0€ | Algoritmo |
| **Posizione Top in categoria** | 800 punti | 0€ | Algoritmo |
| **Badge personalizzato** | 1000 punti | 0€ | Solo grafica |

**⚠️ NON diamo MAI:**
- Sconti (erode margine)
- Spedizioni gratis (erode margine)
- Cashback (erode margine)

**Nota economica - Albero:**
Per guadagnare 300 punti servono ~300€ di acquisti → 30€ commissioni per la piattaforma.
Costo albero: 1€. **Margine garantito: 29€.**

### Regole Anti-Abuso

- I punti scadono dopo 24 mesi di inattività
- Max 1000 punti/giorno da acquisti
- Ordini annullati/rimborsati: punti revocati
- Account sospesi: punti congelati

---

## 3. Classifiche & Premi (Hall of Fame)

### Classifiche Pubbliche

Le classifiche sono visibili a tutti e aggiornate in tempo reale.

```
┌─────────────────────────────────────────────────────────────┐
│                    🏆 TOP ECO-SELLER                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  📅 QUESTA SETTIMANA                                        │
│  ┌────┬──────────────────────┬──────────┬────────────┐      │
│  │ #  │ Azienda              │ CO₂ (kg) │ Prodotti   │      │
│  ├────┼──────────────────────┼──────────┼────────────┤      │
│  │ 🥇 │ Supermercato Verde   │ 450      │ 127        │      │
│  │ 🥈 │ Bio Factory Srl      │ 320      │ 89         │      │
│  │ 🥉 │ Orto di Maria        │ 280      │ 156        │      │
│  │ 4  │ Tech Recycle         │ 250      │ 45         │      │
│  │ 5  │ Panificio Rossi      │ 180      │ 203        │      │
│  └────┴──────────────────────┴──────────┴────────────┘      │
│                                                             │
│  [Settimana] [Mese] [Anno] [Sempre]                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Periodi Classifiche

| Periodo | Reset | Premi |
|---------|-------|-------|
| **Settimanale** | Ogni lunedì | Badge "Top Week" |
| **Mensile** | Primo del mese | Badge + Intervista YouTube |
| **Annuale** | 1 Gennaio | Badge + Feature homepage + Articolo blog |
| **All-time** | Mai | Badge permanente "Leggenda" |

### Premi Automatici (Costo Zero)

| Posizione | Premio | Valore Percepito |
|-----------|--------|------------------|
| 🥇 **1° Mensile** | Intervista YouTube + Badge Gold | ALTO |
| 🥈 **2° Mensile** | Menzione YouTube + Badge Silver | MEDIO |
| 🥉 **3° Mensile** | Badge Bronze | BASSO |
| **New Entry del Mese** | Intervista YouTube | ALTO |
| **Record del Mese** | Feature homepage | ALTO |

---

## 4. YouTube GecoGreen - Contenuti Mensili

### Strategia Contenuti

| Contenuto | Frequenza | Costo | Valore per Azienda |
|-----------|-----------|-------|---------------------|
| **Intervista Eco-Champion** | 1/mese | 0€ (facciamo noi) | ENORME (visibilità) |
| **Intervista New Entry** | 1/mese | 0€ | ENORME |
| **"Dietro le Quinte"** | 2/mese | 0€ | Storytelling |
| **Tips & Tricks** | 4/mese | 0€ | Educational |

### Hall of Fame Mensile

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  🌟 DICEMBRE 2024 - HALL OF FAME                            │
│                                                             │
│  ══════════════════════════════════════════════════════     │
│                                                             │
│  🏆 ECO-CHAMPION DEL MESE                                   │
│  Supermercato Verde Srl                                     │
│  "450 kg di CO₂ risparmiati"                                │
│  [Guarda l'intervista →]                                    │
│                                                             │
│  ──────────────────────────────────────────────────────     │
│                                                             │
│  🚀 NEW ENTRY DEL MESE                                      │
│  Panificio Rossi                                            │
│  "Da 0 a 203 prodotti salvati in 30 giorni!"                │
│  [Guarda l'intervista →]                                    │
│                                                             │
│  ──────────────────────────────────────────────────────     │
│                                                             │
│  📈 RECORD DEL MESE                                         │
│  Bio Factory Srl                                            │
│  "89 prodotti venduti in un singolo giorno!"                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Perché Funziona (Win-Win)

- **Per l'azienda:** Marketing gratuito, visibilità, credibilità
- **Per GecoGreen:** Contenuti gratuiti, testimonianze, SEO
- **Costo:** Solo tempo per registrare/editare video

---

## 5. Admin Dashboard - Gestione Premi

### Reminder Automatici per Admin

Il sistema genera automaticamente task per l'admin:

```
┌─────────────────────────────────────────────────────────────┐
│  📋 ADMIN TODO - CONTENUTI DA CREARE                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ⏰ URGENTE (scade tra 3 giorni)                             │
│  ┌────────────────────────────────────────────────────┐     │
│  │ 🎬 Intervista Eco-Champion Novembre                 │     │
│  │    Azienda: Supermercato Verde Srl                  │     │
│  │    Contatto: mario@supermercatoverde.it             │     │
│  │    CO₂ salvata: 450 kg                              │     │
│  │    [Segna come Completato] [Contatta] [Rimanda]     │     │
│  └────────────────────────────────────────────────────┘     │
│                                                             │
│  📅 QUESTA SETTIMANA                                        │
│  ┌────────────────────────────────────────────────────┐     │
│  │ 🎬 Intervista New Entry Novembre                    │     │
│  │    Azienda: Panificio Rossi                         │     │
│  │    Prima vendita: 15/11/2024                        │     │
│  │    Prodotti salvati: 203                            │     │
│  │    [Segna come Completato] [Contatta] [Rimanda]     │     │
│  └────────────────────────────────────────────────────┘     │
│                                                             │
│  ✅ COMPLETATI QUESTO MESE: 3                               │
│  ⏳ IN ATTESA: 2                                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Workflow Premi

1. **Fine mese:** Sistema calcola automaticamente classifiche
2. **Giorno 1:** Admin riceve notifica con vincitori
3. **Giorno 1-3:** Admin contatta aziende per interviste
4. **Giorno 3-7:** Registrazione interviste
5. **Giorno 7-10:** Pubblicazione su YouTube
6. **Giorno 10:** Aggiornamento Hall of Fame sul sito

### Funzionalità Admin

- **Dashboard classifiche:** Vista in tempo reale
- **Gestione premi:** Assegna badge, boost, menzioni
- **Calendario contenuti:** Pianifica interviste
- **Email template:** Contatta vincitori con un click
- **Storico:** Archivio tutti i premi assegnati

---

## 6. Report CSR per Aziende

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
│  CERTIFICAZIONE                                            │
│                                                            │
│  Questo report è stato generato automaticamente da         │
│  GecoGreen sulla base delle transazioni registrate.        │
│                                                            │
│  ID Report: CSR-2024-TECHCORP-ABC123                       │
│  Data generazione: 31/12/2024                              │
│  Verificabile su: gecogreen.com/csr/ABC123                  │
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

## 7. Community Impact (Contatore Globale)

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

## 8. Implementazione Tecnica

### Nuove Tabelle Database

```sql
-- Classifiche mensili
CREATE TABLE leaderboard_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),

    -- Periodo
    period_type VARCHAR(20) NOT NULL, -- 'WEEKLY', 'MONTHLY', 'YEARLY'
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,

    -- Metriche
    total_co2_saved DECIMAL(10,2) DEFAULT 0,
    total_water_saved DECIMAL(10,2) DEFAULT 0,
    total_products_sold INT DEFAULT 0,
    total_orders INT DEFAULT 0,

    -- Posizione
    rank INT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_leaderboard_period ON leaderboard_snapshots(period_type, period_start);
CREATE INDEX idx_leaderboard_user ON leaderboard_snapshots(user_id);

-- Premi assegnati
CREATE TABLE awards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),

    -- Tipo premio
    award_type VARCHAR(50) NOT NULL, -- 'ECO_CHAMPION', 'NEW_ENTRY', 'RECORD', 'TOP_WEEK'
    period_type VARCHAR(20), -- 'WEEKLY', 'MONTHLY', 'YEARLY'
    period_start DATE,
    period_end DATE,

    -- Dettagli
    title VARCHAR(200),
    description TEXT,
    badge_url VARCHAR(500),

    -- Contenuti associati
    youtube_url VARCHAR(500),
    interview_status VARCHAR(20) DEFAULT 'PENDING', -- 'PENDING', 'SCHEDULED', 'RECORDED', 'PUBLISHED'
    interview_scheduled_at TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP
);

CREATE INDEX idx_awards_user ON awards(user_id);
CREATE INDEX idx_awards_type ON awards(award_type, period_start);

-- Task Admin per contenuti
CREATE TABLE admin_content_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Riferimenti
    award_id UUID REFERENCES awards(id),
    user_id UUID REFERENCES users(id),

    -- Task
    task_type VARCHAR(50) NOT NULL, -- 'INTERVIEW', 'YOUTUBE_PUBLISH', 'SOCIAL_POST'
    title VARCHAR(200) NOT NULL,
    description TEXT,

    -- Stato
    status VARCHAR(20) DEFAULT 'PENDING', -- 'PENDING', 'IN_PROGRESS', 'COMPLETED', 'SKIPPED'
    priority VARCHAR(20) DEFAULT 'NORMAL', -- 'URGENT', 'HIGH', 'NORMAL', 'LOW'
    due_date DATE,

    -- Completamento
    completed_at TIMESTAMP,
    completed_by UUID REFERENCES users(id),
    notes TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_admin_tasks_status ON admin_content_tasks(status, due_date);
CREATE INDEX idx_admin_tasks_award ON admin_content_tasks(award_id);

-- Impact logs (già esistente, espanso)
CREATE TABLE impact_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    order_id UUID REFERENCES orders(id),

    -- Tipo azione
    action_type VARCHAR(50) NOT NULL,  -- 'PURCHASE', 'SALE', 'GIFT', 'BONUS', 'FIRST_PURCHASE', etc.

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
CREATE INDEX idx_impact_logs_date ON impact_logs(created_at);

-- CSR Reports
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
# Classifiche pubbliche
GET  /api/leaderboard                    # Classifica corrente
GET  /api/leaderboard/:period            # Classifica per periodo (weekly/monthly/yearly/alltime)
GET  /api/leaderboard/hall-of-fame       # Hall of Fame con vincitori

# Impact utente
GET  /api/impact/my-stats                # Statistiche utente corrente
GET  /api/impact/certificate/:id         # Singolo certificato
GET  /api/impact/community               # Stats globali community

# EcoCredits
GET  /api/credits/balance                # Saldo punti
GET  /api/credits/history                # Storico punti
POST /api/credits/redeem                 # Spendi EcoCredits

# Admin - Classifiche e Premi
GET  /api/admin/leaderboard/current      # Classifica corrente con dettagli
POST /api/admin/leaderboard/snapshot     # Salva snapshot fine periodo
GET  /api/admin/awards                   # Lista premi
POST /api/admin/awards                   # Assegna premio
PATCH /api/admin/awards/:id              # Aggiorna premio (es. YouTube URL)

# Admin - Task Contenuti
GET  /api/admin/tasks                    # Lista task
GET  /api/admin/tasks/pending            # Task da fare
PATCH /api/admin/tasks/:id               # Aggiorna task
POST /api/admin/tasks/:id/complete       # Segna come completato

# Seller
GET  /api/seller/csr-report              # Genera report CSR
GET  /api/seller/csr-reports             # Lista report passati
GET  /api/seller/my-awards               # I miei premi
```

---

## 9. Roadmap Implementazione

### MVP (Fase 1)
- [x] Campi impatto nel DB
- [ ] Calcolo automatico CO2 per ordine
- [ ] EcoCredits base (guadagno)
- [ ] Widget stats utente in dashboard

### Fase 2
- [ ] Generazione PDF Eco-Scontrino
- [ ] Condivisione social (PNG)
- [ ] Livelli utente
- [ ] Spending EcoCredits (Boost, Alberi)

### Fase 3
- [ ] **Classifiche pubbliche (Settimanali/Mensili/Annuali)**
- [ ] **Hall of Fame**
- [ ] **Admin Dashboard per premi**
- [ ] **Sistema reminder task contenuti**

### Fase 4
- [ ] Report CSR per aziende
- [ ] Integrazione Tree-Nation
- [ ] API partner
- [ ] YouTube integration (link video ai premi)

---

## 10. Note Marketing

### Claim da usare
- "Ogni acquisto salva il pianeta"
- "Trasforma lo spreco in impatto"
- "Il tuo shopping ha un'anima verde"
- "Diventa Eco-Champion del mese!"

### Hashtag
- #GecoGreenImpact
- #ZeroSprecoZeroEmissioni
- #EcoCredits
- #EcoChampion

### Press Kit
- Infografica impatto community
- Case study aziende (interviste YouTube)
- Dati aggregati per giornalisti
- Hall of Fame annuale

---

*Documento creato: Dicembre 2024*
*Version: v2.0 - Aggiunto sistema classifiche, premi e admin dashboard*
