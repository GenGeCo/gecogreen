# GecoGreen - Piano Finanziario

## Budget Iniziale

**Capitale disponibile:** 3.000€

---

## Allocazione Budget

### Distribuzione Consigliata

| Voce | Importo | % | Note |
|------|---------|---|------|
| **Pubblicità/Marketing** | 2.000€ | 67% | Focus principale |
| **Infrastruttura Tech** | 300€ | 10% | Primo anno |
| **Legale/Admin** | 400€ | 13% | T&C, Privacy, Commercialista |
| **Riserva Emergenza** | 300€ | 10% | Buffer imprevisti |
| **TOTALE** | **3.000€** | **100%** | |

---

## Dettaglio Costi

### 1. Infrastruttura Tech (300€/anno)

| Voce | Costo Mensile | Costo Annuo |
|------|---------------|-------------|
| VPS Hetzner (CPX21) | 7€ | 84€ |
| Dominio (.com o .io) | - | 15-50€ |
| Email professionale (opzionale) | 0-5€ | 0-60€ |
| Backup extra (opzionale) | 3€ | 36€ |
| **Subtotale Tech** | **~15€** | **~180€** |

**Servizi Gratuiti Inclusi:**
- Coolify (self-hosted, gratis)
- Supabase free tier (se usato)
- Stripe (gratis fino a transazione)
- OpenAI Moderation API (gratis)
- GitHub (gratis)
- Let's Encrypt SSL (gratis)

### 2. Costi Variabili per Transazione

| Voce | Costo |
|------|-------|
| Stripe fee | 1.4% + 0.25€ per transazione |
| Stripe Connect (payout) | Incluso |
| AI Moderation (Livello 3) | ~0.01€ / 1000 messaggi |

**Esempio su 100€ di vendita:**
```
Prezzo:                    100,00€
Stripe fee:                 -1,65€
Commissione GecoGreen (10%): 10,00€
---------------------------------
Netto al Seller:            88,35€
Guadagno GecoGreen:         10,00€
Costo Stripe su commissione: -0,39€
---------------------------------
Margine Netto GecoGreen:     9,61€
```

### 3. Marketing (2.000€)

#### Fase 1: Pre-lancio (500€)
| Voce | Budget | Obiettivo |
|------|--------|-----------|
| Logo e brand identity | 100-200€ | Fiverr/99designs |
| Materiale social | 100€ | Template Canva Pro |
| Landing page "coming soon" | 0€ | Fatto con stack |
| Outreach manuale venditori | 0€ | Tempo |
| **Subtotale** | **200-300€** | |

#### Fase 2: Lancio (800€)
| Canale | Budget | CPL stimato | Lead stimati |
|--------|--------|-------------|--------------|
| Facebook/Instagram Ads | 400€ | 2-5€ | 80-200 |
| Google Ads (search) | 300€ | 3-8€ | 40-100 |
| Influencer micro (barter) | 100€ | - | Brand awareness |
| **Subtotale** | **800€** | | **120-300 lead** |

#### Fase 3: Crescita (700€)
| Voce | Budget | Note |
|------|--------|------|
| Retargeting | 300€ | Chi ha visitato ma non comprato |
| Content marketing | 200€ | Articoli SEO, guest post |
| Partnership locali | 200€ | Eventi, associazioni |
| **Subtotale** | **700€** | |

### 4. Legale/Amministrativo (400€)

| Voce | Costo | Note |
|------|-------|------|
| T&C e Privacy Policy | 100-200€ | Template + revisione legale |
| Consulenza commercialista | 100-150€ | Setup fiscale, fatturazione |
| SCIA/Autorizzazioni | 50-100€ | Se richieste |
| **Subtotale** | **250-450€** | |

---

## Modello di Revenue

### Stream di Ricavi

| Fonte | % Revenue | Note |
|-------|-----------|------|
| Commissioni vendite | 90% | Core business |
| Fee servizio regali | 5% | 0.50€ su spedizioni gratis |
| Premium seller (futuro) | 5% | Visibilità, analytics |

### Commissione Vendite

**Opzioni da testare:**

| Modello | Pro | Contro |
|---------|-----|--------|
| **10% fisso** | Semplice, prevedibile | Può essere alto per margini bassi |
| **8% + 0.50€** | Bilancia piccoli/grandi ordini | Più complesso |
| **5-15% per categoria** | Ottimizzato per margini | Complesso da comunicare |
| **7% primi 6 mesi** | Attrae venditori | Meno revenue iniziale |

**Raccomandazione:** Iniziare con **10% fisso**, poi ottimizzare con i dati.

---

## Assunzioni del Piano

### Scenario A: Bootstrap (Piano Attuale)

**Questo piano assume:**
- Fondatore lavora SENZA stipendio (investe tempo, non soldi)
- Lavoro da casa (no affitto ufficio)
- Utenze personali già esistenti
- Sviluppo AI-assisted (no dev esterni)

**Adatto a:** Side-project, validazione iniziale, solopreneur

---

### Scenario B: Business "Vero" (Con Stipendi)

Se vuoi pagarti uno stipendio e avere costi realistici:

| Voce | Mensile | Annuo |
|------|---------|-------|
| Stipendio fondatore (netto) | 1.500€ | 18.000€ |
| Costo azienda (+40% contributi) | 2.100€ | 25.200€ |
| Affitto ufficio/coworking | 300€ | 3.600€ |
| Utenze (luce, internet) | 100€ | 1.200€ |
| Commercialista | 150€ | 1.800€ |
| Assicurazione | 50€ | 600€ |
| **TOTALE COSTI FISSI** | **2.700€** | **32.400€** |

**Break-even con stipendio:**
```
Costi fissi mensili:     2.700€
Margine per ordine:      3,38€ (su ordine medio 35€)
Ordini necessari:        2.700 / 3,38 = ~800 ordini/mese
GMV necessario:          800 × 35€ = 28.000€/mese
```

**Conclusione Scenario B:**
- Servono ~800 ordini/mese solo per pareggiare
- Con il piano attuale, raggiungi 800 ordini/mese solo nell'Anno 2
- **Anno 1 in perdita di ~25.000€** invece che in profitto

---

### Quale Scenario Scegliere?

| Fase | Raccomandazione |
|------|-----------------|
| **MVP (0-6 mesi)** | Scenario A - Bootstrap, no stipendio |
| **Validazione (6-12 mesi)** | Scenario A - Reinvesti i profitti |
| **Crescita (Anno 2)** | Transizione a B quando GMV > 30k€/mese |
| **Scale (Anno 3+)** | Scenario B completo + team |

**Il piano seguente usa lo Scenario A (Bootstrap).**
Per lo Scenario B, moltiplica i costi per ~10x e sposta il break-even all'Anno 2.

---

## Timeline: Quando Posso Permettermi Cosa?

### Profitto Cumulato nel Tempo

| Periodo | Profit Periodo | Profit Cumulato | Milestone Raggiungibile |
|---------|----------------|-----------------|------------------------|
| Mese 1-3 | -818€ | -818€ | Investimento iniziale |
| Mese 4-6 | +1.090€ | **+272€** | **Recupero investimento 3k€** |
| Mese 7-9 | +2.760€ | +3.032€ | Riserva emergenza |
| Mese 10-12 | +4.440€ | **+7.472€** | **Posso valutare ufficio** |
| Anno 2 Q1 | +6.620€ | +14.092€ | Marketing boost |
| Anno 2 Q2 | +10.680€ | **+24.772€** | **Posso assumere part-time** |
| Anno 2 Q3 | +14.240€ | +39.012€ | App mobile |
| Anno 2 Q4 | +17.800€ | **+56.812€** | **Posso assumere full-time** |
| Anno 3 | +143.325€ | **+200.137€** | **Team + Espansione EU** |

---

### Quando Pareggio l'Investimento?

```
Investimento iniziale:     3.000€
Pareggiato al:             Mese 5-6
Tempo per recupero:        ~5 mesi
```

---

### Quando Posso Permettermi un Ufficio?

**Costo ufficio/coworking:** 300€/mese = 3.600€/anno

| Periodo | Profit Annuo | Profit - Ufficio | Sostenibile? |
|---------|--------------|------------------|--------------|
| Anno 1 | 7.472€ | 3.872€ | ⚠️ Stretto |
| Anno 2 | 49.340€ | 45.740€ | ✅ Sì |
| Anno 3 | 143.325€ | 139.725€ | ✅ Abbondante |

**Risposta:** Puoi permetterti un ufficio da **inizio Anno 2** (mese 13).
Prima è rischioso - meglio lavorare da casa.

---

### Quando Posso Assumere?

**Costo dipendente:**
- Part-time (20h/sett): ~800€/mese netto = ~12.000€/anno lordo
- Full-time: ~1.500€/mese netto = ~25.000€/anno lordo

| Periodo | Profit Annuo | Part-time (12k) | Full-time (25k) |
|---------|--------------|-----------------|-----------------|
| Anno 1 | 7.472€ | ❌ No | ❌ No |
| Anno 2 | 49.340€ | ✅ Sì (+37k) | ✅ Sì (+24k) |
| Anno 3 | 143.325€ | ✅ (+131k) | ✅ (+118k) |

**Risposta:**
- **Part-time (support):** da **Anno 2 Q1** (mese 13)
- **Full-time (dev/support):** da **Anno 2 Q3** (mese 19)

---

### Piano Progressivo Realistico

```
MESE 1-6:   Bootstrap totale
            └── Tu da casa, niente spese extra
            └── Focus: validare modello

MESE 7-12:  Primi reinvestimenti
            └── +500€/mese in marketing
            └── Ancora da casa
            └── Profit: ~7.500€

MESE 13-18: Prima assunzione
            └── Part-time support (12k/anno)
            └── Coworking opzionale (3.6k/anno)
            └── Profit residuo: ~20k

MESE 19-24: Team minimo
            └── 1 full-time (25k/anno)
            └── Ufficio piccolo (3.6k/anno)
            └── Profit residuo: ~20k

MESE 25-36: Scaling
            └── 2-3 dipendenti
            └── Ufficio vero
            └── Profit: ~100k+ per crescita
```

---

### Riepilogo Milestone

| Cosa | Quando | Profit Cumulato Necessario |
|------|--------|---------------------------|
| Recupero investimento 3k | **Mese 5-6** | 3.000€ |
| Primo ufficio/coworking | **Mese 13** | ~10.000€ |
| Prima assunzione part-time | **Mese 13-15** | ~15.000€ |
| Prima assunzione full-time | **Mese 19-21** | ~40.000€ |
| Primo stipendio a te | **Mese 24** | ~55.000€ |
| Team di 3 persone | **Mese 30** | ~100.000€ |

---

## Business Plan Triennale (Scenario Bootstrap)

### Ipotesi di Crescita

| Parametro | Anno 1 | Anno 2 | Anno 3 |
|-----------|--------|--------|--------|
| Crescita seller | +5/mese | +15/mese | +30/mese |
| Crescita ordini | +10%/mese | +8%/mese | +5%/mese |
| AOV (valore medio ordine) | 35€ | 40€ | 45€ |
| Commissione media | 10% | 10% | 10% |
| Repeat rate buyer | 20% | 35% | 45% |

---

### ANNO 1 - Lancio e Validazione

**Obiettivo:** Validare il modello, raggiungere break-even

| Trimestre | Seller | Ordini/mese | GMV | Commissioni | Costi | Profit |
|-----------|--------|-------------|-----|-------------|-------|--------|
| Q1 | 15 | 65 | 6.825€ | 682€ | 1.500€ | **-818€** |
| Q2 | 30 | 180 | 18.900€ | 1.890€ | 800€ | **+1.090€** |
| Q3 | 45 | 320 | 33.600€ | 3.360€ | 600€ | **+2.760€** |
| Q4 | 60 | 480 | 50.400€ | 5.040€ | 600€ | **+4.440€** |
| **TOTALE ANNO 1** | **60** | **~3.000** | **109.725€** | **10.972€** | **3.500€** | **+7.472€** |

**Investimento iniziale:** 3.000€
**ROI Anno 1:** +149% (7.472€ / 3.000€ - 1)
**Break-even:** Mese 4-5

---

### ANNO 2 - Crescita

**Obiettivo:** Scalare, assumere, espandere categorie

| Trimestre | Seller | Ordini/mese | GMV | Commissioni | Costi | Profit |
|-----------|--------|-------------|-----|-------------|-------|--------|
| Q1 | 105 | 800 | 96.000€ | 9.120€ | 2.500€ | **+6.620€** |
| Q2 | 150 | 1.200 | 144.000€ | 13.680€ | 3.000€ | **+10.680€** |
| Q3 | 195 | 1.600 | 192.000€ | 18.240€ | 4.000€ | **+14.240€** |
| Q4 | 240 | 2.000 | 240.000€ | 22.800€ | 5.000€ | **+17.800€** |
| **TOTALE ANNO 2** | **240** | **~20.000** | **672.000€** | **63.840€** | **14.500€** | **+49.340€** |

**Reinvestimenti Anno 2:**
| Voce | Importo | Note |
|------|---------|------|
| Part-time support | 6.000€ | 500€/mese |
| Marketing scaling | 5.000€ | Ads, influencer |
| Sviluppo app mobile | 3.000€ | PWA → Native |
| Legale/Compliance | 500€ | GDPR, T&C update |
| **Totale reinvestito** | **14.500€** | |

**Utile netto Anno 2:** 49.340€

---

### ANNO 3 - Espansione

**Obiettivo:** Leadership Italia, test Europa

| Trimestre | Seller | Ordini/mese | GMV | Commissioni | Costi | Profit |
|-----------|--------|-------------|-----|-------------|-------|--------|
| Q1 | 330 | 2.800 | 378.000€ | 34.020€ | 8.000€ | **+26.020€** |
| Q2 | 420 | 3.500 | 472.500€ | 42.525€ | 10.000€ | **+32.525€** |
| Q3 | 510 | 4.200 | 567.000€ | 51.030€ | 12.000€ | **+39.030€** |
| Q4 | 600 | 5.000 | 675.000€ | 60.750€ | 15.000€ | **+45.750€** |
| **TOTALE ANNO 3** | **600** | **~55.000** | **2.092.500€** | **188.325€** | **45.000€** | **+143.325€** |

**Reinvestimenti Anno 3:**
| Voce | Importo | Note |
|------|---------|------|
| Team (2 FTE) | 24.000€ | Support + Dev |
| Marketing Europa | 10.000€ | FR, DE, ES |
| Infrastruttura | 5.000€ | Scaling server |
| Compliance EU | 6.000€ | Multi-paese |
| **Totale reinvestito** | **45.000€** | |

**Utile netto Anno 3:** 143.325€

---

### Riepilogo Triennale

| Metrica | Anno 1 | Anno 2 | Anno 3 | Totale |
|---------|--------|--------|--------|--------|
| **GMV** | 109.725€ | 672.000€ | 2.092.500€ | **2.874.225€** |
| **Commissioni** | 10.972€ | 63.840€ | 188.325€ | **263.137€** |
| **Costi/Reinvest** | 3.500€ | 14.500€ | 45.000€ | **63.000€** |
| **Utile Netto** | 7.472€ | 49.340€ | 143.325€ | **200.137€** |
| **Seller attivi** | 60 | 240 | 600 | - |
| **Ordini/anno** | 3.000 | 20.000 | 55.000 | **78.000** |

### Investimento vs Ritorno

```
INVESTIMENTO INIZIALE:     3.000€

RITORNO:
├── Anno 1:   +7.472€  (ROI 149%)
├── Anno 2:  +49.340€  (cumulato: 56.812€)
└── Anno 3: +143.325€  (cumulato: 200.137€)

RITORNO TOTALE 3 ANNI:   200.137€
ROI COMPLESSIVO:         6.571% (200k su 3k investiti)
```

---

### Scenari Alternativi

#### Scenario Pessimistico (-40%)

| Anno | GMV | Commissioni | Utile |
|------|-----|-------------|-------|
| 1 | 66.000€ | 6.600€ | +3.100€ |
| 2 | 400.000€ | 38.000€ | +23.500€ |
| 3 | 1.250.000€ | 112.500€ | +67.500€ |
| **Totale** | **1.716.000€** | **157.100€** | **+94.100€** |

#### Scenario Ottimistico (+50%)

| Anno | GMV | Commissioni | Utile |
|------|-----|-------------|-------|
| 1 | 165.000€ | 16.500€ | +13.000€ |
| 2 | 1.000.000€ | 95.000€ | +80.500€ |
| 3 | 3.100.000€ | 280.000€ | +235.000€ |
| **Totale** | **4.265.000€** | **391.500€** | **+328.500€** |

---

### Milestone di Crescita

| Milestone | Quando | Trigger per Azione |
|-----------|--------|-------------------|
| 🟢 Break-even | Mese 4-5 | Aumentare budget ads |
| 🟢 100 ordini/mese | Mese 6 | Assumere supporto part-time |
| 🟢 50 seller | Mese 8 | Sviluppo app mobile |
| 🟡 500 ordini/mese | Anno 2 Q1 | Team dedicato |
| 🟡 100.000€ GMV/mese | Anno 2 Q3 | Valutare funding |
| 🔴 1M€ GMV annuo | Anno 3 | Espansione Europa |

---

## Proiezioni Finanziarie (Dettaglio Mensile Anno 1)

### Scenario Conservativo (12 mesi)

**Assunzioni:**
- 30 seller attivi a fine anno
- 100 ordini/mese a regime (mese 12)
- Valore medio ordine: 35€
- Commissione: 10%

| Mese | Seller | Ordini | GMV | Commissioni | Costi | Profit |
|------|--------|--------|-----|-------------|-------|--------|
| 1 | 5 | 10 | 350€ | 35€ | 200€ | -165€ |
| 2 | 8 | 20 | 700€ | 70€ | 200€ | -130€ |
| 3 | 12 | 35 | 1.225€ | 122€ | 250€ | -128€ |
| 4 | 15 | 50 | 1.750€ | 175€ | 200€ | -25€ |
| 5 | 18 | 60 | 2.100€ | 210€ | 200€ | +10€ |
| 6 | 22 | 70 | 2.450€ | 245€ | 200€ | +45€ |
| 7 | 25 | 80 | 2.800€ | 280€ | 150€ | +130€ |
| 8 | 27 | 85 | 2.975€ | 297€ | 150€ | +147€ |
| 9 | 28 | 90 | 3.150€ | 315€ | 150€ | +165€ |
| 10 | 29 | 95 | 3.325€ | 332€ | 150€ | +182€ |
| 11 | 30 | 98 | 3.430€ | 343€ | 150€ | +193€ |
| 12 | 30 | 100 | 3.500€ | 350€ | 150€ | +200€ |
| **TOTALE** | | **793** | **27.755€** | **2.774€** | **2.150€** | **+624€** |

### Scenario Ottimistico (12 mesi)

**Assunzioni:**
- 100 seller attivi a fine anno
- 500 ordini/mese a regime
- Valore medio ordine: 40€

| Metrica | Valore |
|---------|--------|
| GMV Anno 1 | ~120.000€ |
| Commissioni | ~12.000€ |
| Costi | ~3.000€ |
| **Profit** | **~9.000€** |

---

## Break-Even Analysis

### Costi Fissi Mensili

| Voce | Importo |
|------|---------|
| Hosting | 15€ |
| Marketing minimo | 100€ |
| Varie | 35€ |
| **Totale Fisso** | **150€** |

### Break-Even Point

```
Break-Even = Costi Fissi / Margine per Ordine

Margine per ordine (su 35€):
  Commissione 10%: 3,50€
  - Stripe fee: ~0,12€
  = Margine netto: 3,38€

Break-Even mensile: 150€ / 3,38€ = ~45 ordini/mese
```

**Obiettivo:** Raggiungere 45+ ordini/mese entro il mese 4.

---

## Metriche da Monitorare (KPI)

### Acquisition

| Metrica | Target M6 | Target M12 |
|---------|-----------|------------|
| CAC (Cost per Acquisition) | < 20€ | < 15€ |
| Seller attivi | 20 | 50 |
| Buyer registrati | 500 | 2.000 |
| Conversion rate (visit→order) | 2% | 3% |

### Revenue

| Metrica | Target M6 | Target M12 |
|---------|-----------|------------|
| GMV mensile | 2.000€ | 5.000€ |
| Commissioni mensili | 200€ | 500€ |
| AOV (Average Order Value) | 30€ | 40€ |
| Orders/month | 70 | 150 |

### Retention

| Metrica | Target |
|---------|--------|
| Repeat buyer rate | > 30% |
| Seller churn | < 10%/mese |
| NPS (Net Promoter Score) | > 40 |

### Quality

| Metrica | Target |
|---------|--------|
| Dispute rate | < 2% |
| Resolution time | < 72h |
| Avg rating seller | > 4.2 |

---

## Rischi Finanziari

| Rischio | Probabilità | Impatto | Mitigazione |
|---------|-------------|---------|-------------|
| Pochi venditori | Alta | Alto | Focus outreach iniziale |
| Frodi/chargeback | Media | Alto | Sistema antifrode robusto |
| Costi marketing troppo alti | Media | Medio | A/B testing, ottimizzazione |
| Competitor copiano | Bassa | Medio | First mover, community |
| Bug critico | Media | Alto | Testing, staging environment |

---

## Piano di Scaling

### Trigger per Investimento Aggiuntivo

| Milestone | Azione | Budget Extra |
|-----------|--------|--------------|
| 100 ordini/mese stabili | Aumentare ads | +500€ |
| 50 seller attivi | Hire part-time support | +800€/mese |
| Break-even raggiunto | Sviluppo app mobile | +2.000€ |
| 1.000€/mese profit | Espansione Europa | +5.000€ |

### Opzioni Funding Futuro

1. **Bootstrapping:** Continuare con profitti reinvestiti
2. **Business Angels:** 20-50k€ per accelerare
3. **Crowdfunding:** Validazione + capitale
4. **Grants EU:** Bandi economia circolare

---

## Checklist Finanziaria Pre-Lancio

- [ ] Aprire conto corrente aziendale (se necessario)
- [ ] Attivare Stripe Connect
- [ ] Setup fatturazione elettronica
- [ ] Definire commissione definitiva
- [ ] Preparare template fattura
- [ ] Consulenza commercialista
- [ ] Budget tracker (spreadsheet)

---

*Documento creato: Dicembre 2024*
*Version: v1.0*
*Nota: Tutti i numeri sono stime e dovranno essere validati con dati reali*
