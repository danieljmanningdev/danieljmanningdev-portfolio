package database

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBackupDoesNotCreateMissingSource(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "missing.db")
	if err := Backup(context.Background(), source, filepath.Join(root, "backup.db")); err == nil {
		t.Fatal("missing source accepted")
	}
	if _, err := os.Stat(source); !os.IsNotExist(err) {
		t.Fatal("missing source was created")
	}
}

func TestBackupAndRestoreLeaveUnrelatedTemporaryFilesAlone(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	source := filepath.Join(root, "source.db")
	backup := filepath.Join(root, "backup.db")
	restored := filepath.Join(root, "restored.db")
	db, err := Open(ctx, source)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL.Exec("CREATE TABLE example (id INTEGER PRIMARY KEY)"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{backup + ".tmp", restored + ".restore.tmp"} {
		if err := os.WriteFile(path, []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := Backup(ctx, source, backup); err != nil {
		t.Fatal(err)
	}
	if err := Restore(ctx, backup, restored, false); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{backup + ".tmp", restored + ".restore.tmp"} {
		body, err := os.ReadFile(path)
		if err != nil || string(body) != "keep" {
			t.Fatalf("unrelated file altered: %s", path)
		}
	}
	if err := os.WriteFile(restored+"-wal", []byte("pending"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Restore(ctx, backup, restored, true); err == nil {
		t.Fatal("restore accepted recovery sidecar")
	}
}

func TestMigrationNamesHaveUniquePositiveVersions(t *testing.T) {
	for _, files := range [][]string{{"001_a.sql", "001_b.sql"}, {"000_a.sql"}, {"-1_a.sql"}, {"001_.sql"}} {
		root := t.TempDir()
		for _, file := range files {
			if err := os.WriteFile(filepath.Join(root, file), []byte("SELECT 1;"), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := discoverMigrations(root); err == nil {
			t.Fatalf("accepted invalid migration set: %v", files)
		}
	}
}
