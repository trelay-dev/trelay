package handler

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/domain"
	"github.com/aftaab/trelay/internal/core/link"
)

type ImportHandler struct {
	linkService *link.Service
}

func NewImportHandler(linkService *link.Service) *ImportHandler {
	return &ImportHandler{
		linkService: linkService,
	}
}

type ImportResult struct {
	Total    int           `json:"total"`
	Imported int           `json:"imported"`
	Skipped  int           `json:"skipped"`
	Failed   int           `json:"failed"`
	Errors   []ImportError `json:"errors,omitempty"`
}

type ImportError struct {
	Row     int    `json:"row"`
	URL     string `json:"url,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Message string `json:"message"`
}

func (h *ImportHandler) Import(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.BadRequest(w, "failed to parse form: "+err.Error())
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, "file is required")
		return
	}
	defer file.Close()

	format := r.FormValue("format")
	if format == "" {
		format = "generic"
	}

	skipDuplicates := r.FormValue("skip_duplicates") != "false"

	result, err := h.processCSV(r.Context(), file, format, skipDuplicates)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *ImportHandler) ImportJSON(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Links          []ImportLink `json:"links"`
		SkipDuplicates bool         `json:"skip_duplicates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	if len(req.Links) == 0 {
		response.BadRequest(w, "no links provided")
		return
	}

	result := h.importLinks(r.Context(), req.Links, req.SkipDuplicates)
	response.JSON(w, http.StatusOK, result)
}

type ImportLink struct {
	URL      string   `json:"url"`
	Slug     string   `json:"slug,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Password string   `json:"password,omitempty"`
	TTLHours int      `json:"ttl_hours,omitempty"`
	FolderID *int64   `json:"folder_id,omitempty"`
}

func (h *ImportHandler) processCSV(ctx context.Context, file io.Reader, format string, skipDuplicates bool) (*ImportResult, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &ImportResult{}, nil
	}

	header := records[0]
	colMap := h.detectColumns(header, format)

	var links []ImportLink
	for i, record := range records[1:] {
		link, err := h.parseRecord(record, colMap, i+2)
		if err != nil {
			continue
		}
		if link.URL != "" {
			links = append(links, link)
		}
	}

	return h.importLinks(ctx, links, skipDuplicates), nil
}

type columnMap struct {
	URL     int
	Slug    int
	Tags    int
	Title   int
	Created int
	Clicks  int
}

func (h *ImportHandler) detectColumns(header []string, format string) columnMap {
	colMap := columnMap{URL: -1, Slug: -1, Tags: -1, Title: -1, Created: -1, Clicks: -1}

	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))

		switch format {
		case "yourls":
			switch col {
			case "keyword":
				colMap.Slug = i
			case "url":
				colMap.URL = i
			case "title":
				colMap.Title = i
			case "timestamp":
				colMap.Created = i
			case "clicks":
				colMap.Clicks = i
			}

		case "shlink":
			switch col {
			case "shortcode", "short_code":
				colMap.Slug = i
			case "longurl", "long_url", "originalurl", "original_url":
				colMap.URL = i
			case "datecreated", "date_created", "createdat", "created_at":
				colMap.Created = i
			case "tags":
				colMap.Tags = i
			case "visitscount", "visits_count", "clicks":
				colMap.Clicks = i
			}

		case "bitly":
			switch col {
			case "link", "long_url", "destination":
				colMap.URL = i
			case "bitlink", "short_url", "shorturl":
				colMap.Slug = i
			case "title":
				colMap.Title = i
			case "created", "created_at":
				colMap.Created = i
			case "tags":
				colMap.Tags = i
			case "total_clicks", "clicks":
				colMap.Clicks = i
			}

		default:
			switch col {
			case "url", "original_url", "long_url", "destination", "target":
				colMap.URL = i
			case "slug", "short", "shortcode", "short_code", "keyword", "alias":
				colMap.Slug = i
			case "tags", "labels":
				colMap.Tags = i
			case "title", "name":
				colMap.Title = i
			case "created", "created_at", "timestamp", "date":
				colMap.Created = i
			case "clicks", "visits", "total_clicks":
				colMap.Clicks = i
			}
		}
	}

	if colMap.URL == -1 && len(header) > 0 {
		colMap.URL = 0
	}

	return colMap
}

func (h *ImportHandler) parseRecord(record []string, colMap columnMap, row int) (ImportLink, error) {
	link := ImportLink{}

	if colMap.URL >= 0 && colMap.URL < len(record) {
		link.URL = strings.TrimSpace(record[colMap.URL])
	}

	if colMap.Slug >= 0 && colMap.Slug < len(record) {
		slug := strings.TrimSpace(record[colMap.Slug])
		if strings.Contains(slug, "/") {
			parts := strings.Split(slug, "/")
			slug = parts[len(parts)-1]
		}
		link.Slug = slug
	}

	if colMap.Tags >= 0 && colMap.Tags < len(record) {
		tagsStr := strings.TrimSpace(record[colMap.Tags])
		if tagsStr != "" {
			tagsStr = strings.ReplaceAll(tagsStr, ";", ",")
			for _, tag := range strings.Split(tagsStr, ",") {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					link.Tags = append(link.Tags, tag)
				}
			}
		}
	}

	return link, nil
}

func (h *ImportHandler) importLinks(ctx context.Context, links []ImportLink, skipDuplicates bool) *ImportResult {
	result := &ImportResult{
		Total:  len(links),
		Errors: []ImportError{},
	}

	for i, link := range links {
		if link.URL == "" {
			result.Failed++
			result.Errors = append(result.Errors, ImportError{
				Row:     i + 1,
				Message: "URL is required",
			})
			continue
		}

		req := domain.CreateLinkRequest{
			URL:  link.URL,
			Slug: link.Slug,
			Tags: link.Tags,
		}

		if link.Password != "" {
			req.Password = link.Password
		}

		if link.TTLHours > 0 {
			req.TTLHours = link.TTLHours
		}

		if link.FolderID != nil {
			req.FolderID = link.FolderID
		}

		_, err := h.linkService.Create(ctx, req)
		if err != nil {
			if err == domain.ErrSlugTaken {
				if skipDuplicates {
					result.Skipped++
					continue
				}
				req.Slug = ""
				_, err = h.linkService.Create(ctx, req)
				if err != nil {
					result.Failed++
					result.Errors = append(result.Errors, ImportError{
						Row:     i + 1,
						URL:     link.URL,
						Slug:    link.Slug,
						Message: err.Error(),
					})
					continue
				}
			} else {
				result.Failed++
				result.Errors = append(result.Errors, ImportError{
					Row:     i + 1,
					URL:     link.URL,
					Slug:    link.Slug,
					Message: err.Error(),
				})
				continue
			}
		}

		result.Imported++
	}

	if len(result.Errors) > 50 {
		result.Errors = result.Errors[:50]
	}

	return result
}

func (h *ImportHandler) Export(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv"
	}

	filter := domain.ListLinksFilter{
		Limit: 10000,
	}

	if folderIDStr := r.URL.Query().Get("folder_id"); folderIDStr != "" {
		if folderID, err := strconv.ParseInt(folderIDStr, 10, 64); err == nil {
			filter.FolderID = &folderID
		}
	}

	links, err := h.linkService.List(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}

	switch format {
	case "json":
		h.exportJSON(w, links)
	default:
		h.exportCSV(w, links)
	}
}

func (h *ImportHandler) exportCSV(w http.ResponseWriter, links []*domain.Link) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=trelay-links.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"url", "slug", "tags", "has_password", "expires_at", "clicks", "created_at"})

	for _, link := range links {
		tags := strings.Join(link.Tags, ",")
		hasPassword := "false"
		if link.HasPassword {
			hasPassword = "true"
		}
		expiresAt := ""
		if link.ExpiresAt != nil {
			expiresAt = link.ExpiresAt.Format(time.RFC3339)
		}

		writer.Write([]string{
			link.OriginalURL,
			link.Slug,
			tags,
			hasPassword,
			expiresAt,
			strconv.FormatInt(link.ClickCount, 10),
			link.CreatedAt.Format(time.RFC3339),
		})
	}
}

func (h *ImportHandler) exportJSON(w http.ResponseWriter, links []*domain.Link) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=trelay-links.json")

	type exportLink struct {
		URL         string   `json:"url"`
		Slug        string   `json:"slug"`
		Tags        []string `json:"tags,omitempty"`
		HasPassword bool     `json:"has_password"`
		ExpiresAt   string   `json:"expires_at,omitempty"`
		Clicks      int64    `json:"clicks"`
		CreatedAt   string   `json:"created_at"`
	}

	var exportLinks []exportLink
	for _, link := range links {
		el := exportLink{
			URL:         link.OriginalURL,
			Slug:        link.Slug,
			Tags:        link.Tags,
			HasPassword: link.HasPassword,
			Clicks:      link.ClickCount,
			CreatedAt:   link.CreatedAt.Format(time.RFC3339),
		}
		if link.ExpiresAt != nil {
			el.ExpiresAt = link.ExpiresAt.Format(time.RFC3339)
		}
		exportLinks = append(exportLinks, el)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"links": exportLinks,
		"total": len(exportLinks),
	})
}
