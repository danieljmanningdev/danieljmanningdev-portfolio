# SQLite backup and restore runbook

The application keeps production data in `/data/portfolio.db` on the Fly volume named `portfolio_data`.

Use two independent recovery layers:

1. **Fly volume snapshots** for infrastructure-level recovery.
2. **Logical SQLite backups** produced with `dbctl` for portable, integrity-checked copies.

A backup is not considered complete until it has been copied away from the production volume and verified.

## Local backups

Create a transactionally consistent SQLite backup:

```bash
go run ./cmd/dbctl \
  -action backup \
  -database ./data/app.db \
  -file ./backups/app-$(date -u +%Y%m%dT%H%M%SZ).db
```

Verify any backup before relying on it:

```bash
go run ./cmd/dbctl \
  -action verify \
  -file ./backups/app-20260829T190000Z.db
```

The backup command uses SQLite `VACUUM INTO`, writes through a temporary file, runs `PRAGMA quick_check`, applies `0600` permissions, and publishes the result with an atomic rename.

## Local restore rehearsal

Restore into a separate path first:

```bash
go run ./cmd/dbctl \
  -action restore \
  -file ./backups/app-20260829T190000Z.db \
  -database ./data/restore-test.db \
  -force
```

Start a temporary application instance against the restored database:

```bash
DATABASE_PATH=./data/restore-test.db \
APP_PORT=8081 \
go run ./cmd/server
```

Confirm that the journal, clients, projects, contracts, and administrator records are present before replacing any active database.

## Production logical backup

The production image includes `/dbctl`.

Start or wake the application Machine, then create a consistent backup inside the Machine:

```bash
BACKUP_NAME="portfolio-$(date -u +%Y%m%dT%H%M%SZ).db"

fly ssh console \
  -a danieljmanningdev-portfolio \
  -C "/dbctl -action backup -database /data/portfolio.db -file /tmp/${BACKUP_NAME}"
```

Copy it to an encrypted local backup directory:

```bash
mkdir -p ./backups/production

fly sftp get \
  "/tmp/${BACKUP_NAME}" \
  "./backups/production/${BACKUP_NAME}" \
  -a danieljmanningdev-portfolio
```

Verify the downloaded file locally:

```bash
go run ./cmd/dbctl \
  -action verify \
  -file "./backups/production/${BACKUP_NAME}"
```

Remove the temporary remote copy after verification:

```bash
fly ssh console \
  -a danieljmanningdev-portfolio \
  -C "rm -f /tmp/${BACKUP_NAME}"
```

Keep at least one recent verified copy outside Fly.io.

## Production restore

Do not replace `/data/portfolio.db` while the web process is running or while another Machine has the volume mounted.

Preferred recovery order:

1. Preserve the current volume and database before changing anything.
2. Verify the logical backup locally with `dbctl`.
3. Rehearse the restore against a separate local database.
4. Stop application writes.
5. Restore into a maintenance environment or a replacement volume.
6. Run `PRAGMA quick_check` and application smoke tests.
7. Attach the recovered volume to the production Machine only after validation.

For a logical restore in a stopped maintenance container:

```bash
/dbctl \
  -action restore \
  -file /data/restore-source.db \
  -database /data/portfolio.db \
  -force
```

`-force` is deliberately required. Before invoking it, create another copy of the current database whenever the file remains readable.

## Fly volume snapshots

Fly automatically creates daily snapshots for volumes. The default retention is five days, and Fly supports configuring retention from one to sixty days.

List the production volume and its snapshots:

```bash
fly volumes list -a danieljmanningdev-portfolio
fly volumes snapshots list <volume-id>
```

Create an on-demand snapshot before a risky production operation:

```bash
fly volumes snapshots create <volume-id>
```

Restore a snapshot into a new volume rather than overwriting the only existing volume:

```bash
fly volumes create portfolio_data_recovery \
  --snapshot-id <snapshot-id> \
  --size <existing-volume-size-gb> \
  --region lhr \
  -a danieljmanningdev-portfolio
```

Relevant Fly documentation:

- https://fly.io/docs/volumes/snapshots/
- https://fly.io/docs/volumes/volume-manage/
- https://fly.io/docs/flyctl/sftp/

## Suggested schedule

- Fly automatic snapshot: daily.
- Logical SQLite backup: before every migration or high-risk production change.
- Routine logical backup: at least weekly while the workspace is actively used.
- Restore rehearsal: monthly, and after any backup-process change.

Record the date, filename, verification result, and restore-test result for each retained backup.
