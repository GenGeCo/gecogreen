# GecoGreen - Roadmap & Feature Backlog

## 🚀 In Sviluppo (Current Sprint)

### Quantità con Unità di Misura
- [ ] Aggiungere campo `unit` alla tabella products (PIECE, KG, G, L, ML, CUSTOM)
- [ ] Aggiungere campo `unit_custom` per unità personalizzate
- [ ] Modificare form inserimento prodotto
- [ ] Aggiornare visualizzazione prodotti

### Dati Fatturazione Completi
- [ ] Aggiungere tabella `billing_info`:
  - Codice Fiscale (Italia)
  - Codice Univoco SDI / PEC (fattura elettronica)
  - VAT ID europeo
  - Indirizzo fatturazione completo
- [ ] Sezione fatturazione nel profilo utente
- [ ] Validazione dati fatturazione

### Cambio Tipo Account (Privato ↔ Azienda)
- [ ] Funzionalità nel profilo per passare da PRIVATE a BUSINESS
- [ ] Funzionalità nel profilo per passare da BUSINESS a PRIVATE
- [ ] Warning e conferme durante il cambio
- [ ] Validazione dati aziendali quando si passa a BUSINESS
- [ ] Gestione visibilità dati aziendali quando si torna PRIVATE

### Asta Olandese (Dutch Auction)
- [ ] Implementare form inserimento asta olandese:
  - Prezzo iniziale
  - Riduzione prezzo ogni X ore
  - Prezzo minimo finale
- [ ] Job automatico che riduce il prezzo ogni X ore
- [ ] Visualizzazione countdown prossima riduzione prezzo
- [ ] Notifiche quando prezzo scende

---

## 📦 Roadmap - Q1 2026

### Digital Freight Forwarders Integration
**Status**: 📋 Pianificato
**Priority**: Alta

Integrazione con servizi di spedizione digitali per gestire logistica complessa:
- [ ] Ricerca e selezione Digital Freight Forwarder partner
- [ ] API integration per quotazioni real-time
- [ ] Tracking spedizioni integrate
- [ ] Gestione documentazione doganale
- [ ] Calcolo automatico costi spedizione
- [ ] Multi-carrier comparison

**Note**: Attualmente disponibile come opzione "Coming Soon" nel form spedizione

---

## 🎯 Feature Future (Backlog)

### Pagamenti Integrati
- [ ] Integrazione Stripe/PayPal
- [ ] Pagamento in-app sicuro
- [ ] Gestione escrow per transazioni
- [ ] Split payment (commissioni GecoGreen)

### Sistema Recensioni e Rating
- [ ] Recensioni prodotti
- [ ] Rating venditore/acquirente
- [ ] Sistema reputazione
- [ ] Badge affidabilità

### Geolocalizzazione Avanzata
- [ ] Mappa prodotti nelle vicinanze
- [ ] Filtro per distanza
- [ ] Suggerimenti basati su posizione
- [ ] Punti ritiro multipli su mappa

### Notifiche Push
- [ ] Notifiche browser
- [ ] Email transazionali
- [ ] Promemoria scadenze prodotti
- [ ] Alert nuovi prodotti categorie preferite

### Analytics per Venditori
- [ ] Dashboard statistiche vendite
- [ ] Andamento prezzi per categoria
- [ ] Suggerimenti pricing
- [ ] Report fiscali

### App Mobile
- [ ] React Native / Flutter app
- [ ] Scan barcode prodotti
- [ ] Upload foto da fotocamera
- [ ] Notifiche push native

---

## ✅ Completate (Sprint Precedenti)

### Sprint 1 - MVP Base
- [x] Sistema autenticazione JWT
- [x] Gestione profilo utente
- [x] CRUD prodotti
- [x] Upload immagini su R2
- [x] Filtri e ricerca prodotti
- [x] Separazione foto prodotto/scadenza
- [x] Sistema moderazione immagini (OCR stub)
- [x] Account PRIVATE vs BUSINESS
- [x] Gestione più sedi per aziende
- [x] Deploy su Coolify
- [x] Domini configurati (gecogreen.com, api.gecogreen.com)

---

## 🐛 Bug Known

Nessun bug critico al momento.

---

## 💡 Idee da Valutare

- Gamification (punti per prodotti salvati dallo spreco)
- Abbonamenti premium per aziende
- Integrazione con inventory management systems
- QR Code per prodotti in negozio
- Social features (condividi prodotto)
- Wishlist e preferiti
- Sistema offerte/controfferte
- Bundle di prodotti
