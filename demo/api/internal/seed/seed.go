package seed

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oullin/inertia-go/demo/api/internal/database"
)

func Run(db *sql.DB) error {
	if err := database.Truncate(db); err != nil {
		return fmt.Errorf("seed truncate: %w", err)
	}

	now := time.Now()

	if err := seedInvites(db, now); err != nil {
		return fmt.Errorf("seed invites: %w", err)
	}

	if err := seedUploads(db, now); err != nil {
		return fmt.Errorf("seed uploads: %w", err)
	}

	if err := seedApprovals(db, now); err != nil {
		return fmt.Errorf("seed approvals: %w", err)
	}

	if err := database.SetCounter(db, "priority_escalations", 18); err != nil {
		return fmt.Errorf("seed counter: %w", err)
	}

	return nil
}

func seedInvites(db *sql.DB, now time.Time) error {
	invites := []struct {
		id, name, email, role, status string
		ago                           time.Duration
	}{
		{"invite_100", "Aria Lim", "aria@northstarhq.test", "Operator", "Accepted", 9 * time.Minute},
		{"invite_101", "Noah Chen", "noah@northstarhq.test", "Manager", "Pending", 41 * time.Minute},
		{"invite_102", "Priya Kapoor", "priya@northstarhq.test", "Executive", "Accepted", 1 * time.Hour},
		{"invite_103", "Leo Ruiz", "leo@northstarhq.test", "Analyst", "Pending", 2 * time.Hour},
		{"invite_104", "Sana Okafor", "sana@northstarhq.test", "Operator", "Accepted", 3 * time.Hour},
		{"invite_105", "Kai Tanaka", "kai@northstarhq.test", "Manager", "Expired", 5 * time.Hour},
		{"invite_106", "Dani Alves", "dani@northstarhq.test", "Analyst", "Accepted", 6 * time.Hour},
		{"invite_107", "Remi Dubois", "remi@northstarhq.test", "Operator", "Pending", 8 * time.Hour},
		{"invite_108", "Ines Moreno", "ines@northstarhq.test", "Executive", "Pending", 24 * time.Hour},
		{"invite_109", "Jules Park", "jules@northstarhq.test", "Manager", "Accepted", 48 * time.Hour},
		{"invite_110", "Maya Tan", "maya@northstarhq.test", "Executive", "Accepted", 3 * 24 * time.Hour},
		{"invite_111", "Ava Gomez", "ava@northstarhq.test", "Manager", "Accepted", 4 * 24 * time.Hour},
		{"invite_112", "Zara Hussain", "zara@northstarhq.test", "Analyst", "Expired", 5 * 24 * time.Hour},
		{"invite_113", "Tomás Reyes", "tomas@northstarhq.test", "Operator", "Accepted", 6 * 24 * time.Hour},
		{"invite_114", "Mila Novak", "mila@northstarhq.test", "Executive", "Pending", 7 * 24 * time.Hour},
		{"invite_115", "Ethan Osei", "ethan@northstarhq.test", "Analyst", "Accepted", 8 * 24 * time.Hour},
		{"invite_116", "Lena Vogt", "lena@northstarhq.test", "Operator", "Expired", 10 * 24 * time.Hour},
		{"invite_117", "Rohan Mehta", "rohan@northstarhq.test", "Manager", "Pending", 12 * 24 * time.Hour},
		{"invite_118", "Cleo Andersen", "cleo@northstarhq.test", "Analyst", "Accepted", 14 * 24 * time.Hour},
		{"invite_119", "Felix Braun", "felix@northstarhq.test", "Executive", "Accepted", 18 * 24 * time.Hour},
	}

	for _, inv := range invites {
		if err := database.CreateInviteAt(db, inv.id, inv.name, inv.email, inv.role, inv.status, now.Add(-inv.ago)); err != nil {
			return err
		}
	}

	return nil
}

func seedUploads(db *sql.DB, now time.Time) error {
	uploads := []struct {
		id, label, filename, size, status string
		ago                               time.Duration
	}{
		{"upload_s00", "March collections export", "collections-march.csv", "84 KB", "Processed", 12 * time.Minute},
		{"upload_s01", "Sales handoff packet", "handoff-q2.pdf", "1.2 MB", "Ready", 2 * time.Hour},
		{"upload_s02", "Q1 revenue reconciliation", "revenue-q1.xlsx", "320 KB", "Processed", 3 * time.Hour},
		{"upload_s03", "Partner channel attribution", "channel-attr.csv", "156 KB", "Queued", 4 * time.Hour},
		{"upload_s04", "Northwind contract redlines", "northwind-redlines.pdf", "2.4 MB", "Processed", 6 * time.Hour},
		{"upload_s05", "Atlas compliance artifacts", "atlas-compliance.zip", "4.1 MB", "Ready", 8 * time.Hour},
		{"upload_s06", "Cedar renewal terms", "cedar-renewal.docx", "92 KB", "Queued", 24 * time.Hour},
		{"upload_s07", "Helio onboarding checklist", "helio-onboarding.pdf", "540 KB", "Processed", 48 * time.Hour},
		{"upload_s08", "Juniper SLA benchmarks", "juniper-sla.xlsx", "210 KB", "Ready", 3 * 24 * time.Hour},
		{"upload_s09", "Mistral pricing model", "mistral-pricing.json", "18 KB", "Processed", 4 * 24 * time.Hour},
		{"upload_s10", "Summit invoice batch", "summit-invoices-q1.csv", "1.8 MB", "Queued", 5 * 24 * time.Hour},
		{"upload_s11", "Polar integration spec", "polar-api-spec.pdf", "3.2 MB", "Ready", 6 * 24 * time.Hour},
		{"upload_s12", "Orchid deal room export", "orchid-deal-room.zip", "7.6 MB", "Processed", 8 * 24 * time.Hour},
		{"upload_s13", "Bluebird NDA packet", "bluebird-nda.pdf", "420 KB", "Queued", 10 * 24 * time.Hour},
		{"upload_s14", "Beacon ops runbook", "beacon-runbook.docx", "1.1 MB", "Processed", 14 * 24 * time.Hour},
	}

	for _, upl := range uploads {
		if err := database.CreateUploadAt(db, upl.id, upl.label, upl.filename, upl.size, upl.status, now.Add(-upl.ago)); err != nil {
			return err
		}
	}

	return nil
}

func seedApprovals(db *sql.DB, now time.Time) error {
	approvals := []struct {
		id, label, status string
		ago               time.Duration
	}{
		{"approval_s00", "Northwind priority routing", "Approved", 16 * time.Minute},
		{"approval_s01", "Atlas expansion guardrail", "Queued", 47 * time.Minute},
		{"approval_s02", "Cedar payment override", "Approved", 1 * time.Hour},
		{"approval_s03", "Helio discount authorization", "Synced", 2 * time.Hour},
		{"approval_s04", "Juniper compliance waiver", "Queued", 4 * time.Hour},
		{"approval_s05", "Mistral seat expansion", "Approved", 5 * time.Hour},
		{"approval_s06", "Summit invoice fast-track", "Synced", 8 * time.Hour},
		{"approval_s07", "Polar onboarding exception", "Queued", 24 * time.Hour},
		{"approval_s08", "Orchid renewal acceleration", "Approved", 2 * 24 * time.Hour},
		{"approval_s09", "Bluebird SLA amendment", "Synced", 3 * 24 * time.Hour},
		{"approval_s10", "Beacon routing promotion", "Approved", 4 * 24 * time.Hour},
		{"approval_s11", "Northwind contract extension", "Queued", 5 * 24 * time.Hour},
		{"approval_s12", "Atlas budget reallocation", "Approved", 7 * 24 * time.Hour},
		{"approval_s13", "Cedar autopay enrollment", "Synced", 10 * 24 * time.Hour},
		{"approval_s14", "Helio partner tier upgrade", "Queued", 14 * 24 * time.Hour},
	}

	for _, apr := range approvals {
		if err := database.CreateApprovalAt(db, apr.id, apr.label, apr.status, now.Add(-apr.ago)); err != nil {
			return err
		}
	}

	return nil
}
