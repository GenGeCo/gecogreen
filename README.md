# GecoGreen

**Marketplace anti-spreco B2B/B2C**

> "Dai valore a ciò che resta"

**Dominio:** www.gecogreen.com

---

## Cos'è GecoGreen?

GecoGreen connette aziende con eccedenze (alimentari, materiali, prodotti) a consumatori e altre aziende, riducendo lo spreco e generando valore economico per tutti.

**Caratteristiche principali:**
- Multi-categoria (non solo cibo)
- B2B + B2C nella stessa piattaforma
- Sezione "Regalo" gratuita
- Sistema antifrode integrato (Escrow, QR, moderazione AI)
- Scalabile a livello europeo

---

## Documentazione

| Documento | Descrizione |
|-----------|-------------|
| [01_PROJECT.md](docs/01_PROJECT.md) | Vision, mission, target, competitors |
| [02_ARCHITECTURE.md](docs/02_ARCHITECTURE.md) | Stack tecnico, API, infrastruttura |
| [03_USER_FLOWS.md](docs/03_USER_FLOWS.md) | Ruoli utente, pagine, flussi |
| [04_BUSINESS_RULES.md](docs/04_BUSINESS_RULES.md) | Regole antifrode, dispute, moderazione |
| [05_FINANCIAL_PLAN.md](docs/05_FINANCIAL_PLAN.md) | Costi, revenue, proiezioni |
| [06_TODO_LAUNCH.md](docs/06_TODO_LAUNCH.md) | Checklist pre-lancio |
| [07_IMPACT_SYSTEM.md](docs/07_IMPACT_SYSTEM.md) | EcoCredits, certificati impatto, CSR reports |

**Schema Database:** [schemas/database.sql](schemas/database.sql)

---

## Stack Tecnologico

```
Frontend:  SvelteKit + TypeScript + TailwindCSS + DaisyUI
Mobile:    Capacitor (iOS/Android wrapper)
Backend:   Go (Fiber framework)
Database:  PostgreSQL
Cache:     Redis
Storage:   Cloudflare R2 (foto prodotti)
Email:     Resend (transazionali)
Hosting:   Hetzner VPS + Coolify + Docker
Pagamenti: Stripe Connect

UI Theme:  Verde Geco (#00C853)
```

---

## Quick Start (Sviluppo)

```bash
# 1. Clona il repo
git clone https://github.com/[tuo-username]/gecogreen.git
cd gecogreen

# 2. Avvia i servizi con Docker
docker-compose up -d

# 3. Esegui le migrazioni DB
psql -U gecogreen -d gecogreen -f schemas/database.sql

# 4. Avvia il backend (in un altro terminale)
cd backend && go run cmd/server/main.go

# 5. Avvia il frontend (in un altro terminale)
cd frontend && npm install && npm run dev
```

---

## Struttura Progetto

```
gecogreen/
├── docs/                    # Documentazione progetto
├── schemas/                 # Schema database SQL
├── backend/                 # API Go (da creare)
├── frontend/                # App SvelteKit + Capacitor (da creare)
├── docker-compose.yml       # Setup sviluppo locale
└── README.md
```

---

## Roadmap

- [x] Documentazione progetto
- [x] Schema database
- [ ] Setup infrastruttura
- [ ] Backend MVP
- [ ] Frontend MVP
- [ ] Testing
- [ ] Lancio beta
- [ ] Lancio pubblico

---

## Budget

- **Capitale iniziale:** 3.000€
- **Focus:** 67% marketing, 10% tech, 13% legale, 10% riserva
- **Break-even stimato:** ~45 ordini/mese (mese 4-5)

---

## Contatti

- **Progetto:** GecoGreen
- **Dominio:** www.gecogreen.com
- **Stato:** In sviluppo

---

*Creato con Claude Code + VS Code*
