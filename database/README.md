# Database Migrations

## Come Applicare le Migration

### Opzione 1: Manualmente via psql
```bash
# Connettiti al database PostgreSQL
docker exec -it <postgres_container_name> psql -U gecogreen -d gecogreen

# Applica le migration in ordine
\i /path/to/001_add_quantity_unit.sql
\i /path/to/002_add_billing_info.sql
```

### Opzione 2: Via Docker
```bash
# Copia gli script nel container
docker cp database/001_add_quantity_unit.sql <postgres_container_name>:/tmp/
docker cp database/002_add_billing_info.sql <postgres_container_name>:/tmp/

# Esegui gli script
docker exec -it <postgres_container_name> psql -U gecogreen -d gecogreen -f /tmp/001_add_quantity_unit.sql
docker exec -it <postgres_container_name> psql -U gecogreen -d gecogreen -f /tmp/002_add_billing_info.sql
```

### Opzione 3: Tramite Coolify
1. Accedi alla console PostgreSQL su Coolify
2. Copia e incolla il contenuto di ogni file SQL
3. Esegui in ordine numerico

## Migration Disponibili

| # | Nome | Descrizione | Status |
|---|------|-------------|--------|
| 001 | add_quantity_unit | Aggiunge unità di misura ai prodotti | ⏳ Pending |
| 002 | add_billing_info | Aggiunge dati fatturazione completi | ⏳ Pending |

## Note

- Eseguire sempre le migration in ordine numerico
- Fare backup del database prima di applicare migration in produzione
- Testare prima in locale/staging
