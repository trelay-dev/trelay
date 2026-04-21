// Seed inserts demo links and folders for manual UI testing.
// Run from the trelay repo root so .env loads: go run ./cmd/seed  (or: make seed)
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aftaab/trelay/internal/config"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/folder"
	"github.com/aftaab/trelay/internal/core/link"
	"github.com/aftaab/trelay/internal/storage/sqlite"
)

// All slugs use the seed- prefix so we can remove them on re-run.
var seedSlugs = []string{
	"seed-welcome",
	"seed-protected",
	"seed-expires-soon",
	"seed-expires-week",
	"seed-onetime",
	"seed-preview",
	"seed-other",
	"seed-trashed",
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sqlite.Open(cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		log.Fatal("migrate: ", err)
	}

	lr := sqlite.NewLinkRepository(db)
	cr := sqlite.NewClickRepository(db)
	linkSvc := link.NewService(lr, cfg.App.SlugLength, cfg.App.CustomDomains)
	folderSvc := folder.NewService(sqlite.NewFolderRepository(db))
	ctx := context.Background()

	for _, slug := range seedSlugs {
		_ = linkSvc.HardDelete(ctx, slug)
	}

	projID := mustFolder(ctx, folderSvc, "Demo Projects")
	archID := mustFolder(ctx, folderSvc, "Demo Archive")

	type mk struct {
		req domain.CreateLinkRequest
	}
	creates := []mk{
		{domain.CreateLinkRequest{
			URL:      "https://example.com/welcome",
			Slug:     "seed-welcome",
			Tags:     []string{"demo", "marketing"},
			FolderID: projID,
		}},
		{domain.CreateLinkRequest{
			URL:      "https://example.com/private-doc",
			Slug:     "seed-protected",
			Password: "demo",
			Tags:     []string{"demo", "security"},
			FolderID: projID,
		}},
		{domain.CreateLinkRequest{
			URL:      "https://example.com/flash-sale",
			Slug:     "seed-expires-soon",
			Tags:     []string{"demo"},
			TTLHours: 2,
			FolderID: projID,
		}},
		{domain.CreateLinkRequest{
			URL:      "https://example.com/webinar",
			Slug:     "seed-expires-week",
			Tags:     []string{"demo"},
			TTLHours: 24 * 4,
			FolderID: archID,
		}},
		{domain.CreateLinkRequest{
			URL:       "https://example.com/one-shot",
			Slug:      "seed-onetime",
			Tags:      []string{"demo"},
			IsOneTime: true,
		}},
		{domain.CreateLinkRequest{
			URL:             "https://example.com/long-article",
			Slug:            "seed-preview",
			Tags:            []string{"demo", "content"},
			FolderID:        projID,
			OGTitle:         "Custom preview title (seed)",
			OGDescription:   "This text overrides what we would fetch from the target page.",
			OGImageURL:      "https://picsum.photos/seed/trelaypreview/640/360",
		}},
		{domain.CreateLinkRequest{
			URL:      "https://example.org/other",
			Slug:     "seed-other",
			Tags:     []string{"other", "demo"},
			FolderID: archID,
		}},
	}

	for _, c := range creates {
		if _, err := linkSvc.Create(ctx, c.req); err != nil {
			log.Fatalf("create %s: %v", c.req.Slug, err)
		}
	}

	welcome, err := lr.GetBySlug(ctx, "seed-welcome")
	if err != nil {
		log.Fatal(err)
	}
	// Bump link.click_count (shown in list) and insert matching rows into `clicks`
	// so daily analytics / dashboard chart see the same activity (IncrementClickCount alone does not).
	for i := 0; i < 12; i++ {
		_ = lr.IncrementClickCount(ctx, welcome.ID)
	}
	now := time.Now().UTC()
	for i := 0; i < 12; i++ {
		daysAgo := i % 7
		ts := now.AddDate(0, 0, -daysAgo)
		ts = time.Date(ts.Year(), ts.Month(), ts.Day(), 10+(i%8), (i*5)%60, 0, 0, time.UTC)
		click := &domain.Click{
			LinkID:    welcome.ID,
			Timestamp: ts,
			Referrer:  "direct",
			UserAgent: "Mozilla/5.0 (seed demo)",
		}
		if err := cr.Record(ctx, click); err != nil {
			log.Fatalf("record seed click: %v", err)
		}
	}

	tr, err := linkSvc.Create(ctx, domain.CreateLinkRequest{
		URL:  "https://example.com/old-campaign",
		Slug: "seed-trashed",
		Tags: []string{"demo"},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := linkSvc.Delete(ctx, tr.Slug); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Seed finished OK.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Short link password:  demo   →  open /seed-protected (or use ?p=demo)")
	fmt.Fprintln(os.Stderr, "  seed-welcome: 12 demo clicks + analytics rows (dashboard Click Activity)")
	fmt.Fprintln(os.Stderr, "  Tags used: demo, marketing, security, content, other")
	fmt.Fprintln(os.Stderr, "  Folders: Demo Projects, Demo Archive")
	fmt.Fprintln(os.Stderr, "")
}

func mustFolder(ctx context.Context, fs *folder.Service, name string) *int64 {
	list, err := fs.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range list {
		if f.Name == name {
			id := f.ID
			return &id
		}
	}
	f, err := fs.Create(ctx, domain.CreateFolderRequest{Name: name})
	if err != nil {
		log.Fatal(err)
	}
	id := f.ID
	return &id
}
