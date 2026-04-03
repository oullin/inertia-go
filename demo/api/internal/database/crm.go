package database

import (
	"database/sql"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int64
	Name       string
	Email      string
	Password   string
	VerifiedAt *time.Time
}

type Organization struct {
	ID            int64
	Name          string
	ContactsCount int
}

type Contact struct {
	ID               int64
	OrganizationID   *int64
	OrganizationName string
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	Address          string
	City             string
	Region           string
	Country          string
	PostalCode       string
	IsFavorite       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Note struct {
	ID          int64
	ContactID   int64
	ContactName string
	UserID      int64
	UserName    string
	Body        string
	CreatedAt   time.Time
}

func CreateUser(db *sql.DB, name, email, password string, verifiedAt *time.Time) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return 0, err
	}

	result, err := db.Exec(
		"INSERT INTO users (name, email, password, verified_at) VALUES (?, ?, ?, ?)",
		name,
		email,
		string(hashedPassword),
		verifiedAt,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func FindUserByEmail(db *sql.DB, email string) (*User, error) {
	row := db.QueryRow(`
		SELECT id, name, email, password, verified_at
		FROM users
		WHERE lower(email) = lower(?)
		LIMIT 1
	`, email)

	return scanUser(row)
}

func FindUserByID(db *sql.DB, id int64) (*User, error) {
	row := db.QueryRow(`
		SELECT id, name, email, password, verified_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id)

	return scanUser(row)
}

func ListOrganizations(db *sql.DB, search string) ([]Organization, error) {
	query := `
		SELECT o.id, o.name, COUNT(c.id) AS contacts_count
		FROM organizations o
		LEFT JOIN contacts c ON c.organization_id = o.id
	`

	var args []any

	if search = strings.TrimSpace(search); search != "" {
		query += " WHERE lower(o.name) LIKE lower(?)"
		args = append(args, "%"+search+"%")
	}

	query += " GROUP BY o.id, o.name ORDER BY o.name ASC"

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Organization, error) {
		var org Organization

		if err := scan(&org.ID, &org.Name, &org.ContactsCount); err != nil {
			return Organization{}, err
		}

		return org, nil
	}), nil
}

func GetOrganization(db *sql.DB, id int64) (*Organization, error) {
	row := db.QueryRow(`
		SELECT o.id, o.name, COUNT(c.id) AS contacts_count
		FROM organizations o
		LEFT JOIN contacts c ON c.organization_id = o.id
		WHERE o.id = ?
		GROUP BY o.id, o.name
	`, id)

	var org Organization

	if err := row.Scan(&org.ID, &org.Name, &org.ContactsCount); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &org, nil
}

func CreateOrganization(db *sql.DB, name string) (int64, error) {
	result, err := db.Exec("INSERT INTO organizations (name) VALUES (?)", strings.TrimSpace(name))

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateOrganization(db *sql.DB, id int64, name string) error {
	_, err := db.Exec(
		"UPDATE organizations SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		strings.TrimSpace(name),
		id,
	)

	return err
}

func CountOrganizations(db *sql.DB) int {
	var count int

	db.QueryRow("SELECT COUNT(*) FROM organizations").Scan(&count)

	return count
}

func ListContacts(db *sql.DB, search string, favoritesOnly bool) ([]Contact, error) {
	query := `
		SELECT
			c.id,
			c.organization_id,
			COALESCE(o.name, ''),
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.address,
			c.city,
			c.region,
			c.country,
			c.postal_code,
			c.is_favorite,
			c.created_at,
			c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE 1 = 1
	`

	var args []any

	if search = strings.TrimSpace(search); search != "" {
		query += `
			AND (
				lower(c.first_name) LIKE lower(?)
				OR lower(c.last_name) LIKE lower(?)
				OR lower(c.email) LIKE lower(?)
				OR lower(COALESCE(o.name, '')) LIKE lower(?)
			)
		`
		like := "%" + search + "%"
		args = append(args, like, like, like, like)
	}

	if favoritesOnly {
		query += " AND c.is_favorite = 1"
	}

	query += " ORDER BY c.first_name ASC, c.last_name ASC"

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Contact, error) {
		return scanContact(scan)
	}), nil
}

func ListContactsByOrganization(db *sql.DB, organizationID int64) ([]Contact, error) {
	rows, err := db.Query(`
		SELECT
			c.id,
			c.organization_id,
			COALESCE(o.name, ''),
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.address,
			c.city,
			c.region,
			c.country,
			c.postal_code,
			c.is_favorite,
			c.created_at,
			c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE c.organization_id = ?
		ORDER BY c.first_name ASC, c.last_name ASC
	`, organizationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Contact, error) {
		return scanContact(scan)
	}), nil
}

func GetContact(db *sql.DB, id int64) (*Contact, error) {
	row := db.QueryRow(`
		SELECT
			c.id,
			c.organization_id,
			COALESCE(o.name, ''),
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.address,
			c.city,
			c.region,
			c.country,
			c.postal_code,
			c.is_favorite,
			c.created_at,
			c.updated_at
		FROM contacts c
		LEFT JOIN organizations o ON o.id = c.organization_id
		WHERE c.id = ?
		LIMIT 1
	`, id)

	contact, err := scanContact(row.Scan)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &contact, nil
}

func CreateContact(db *sql.DB, contact Contact) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO contacts (
			organization_id,
			first_name,
			last_name,
			email,
			phone,
			address,
			city,
			region,
			country,
			postal_code,
			is_favorite
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		nullableInt64(contact.OrganizationID),
		strings.TrimSpace(contact.FirstName),
		strings.TrimSpace(contact.LastName),
		strings.TrimSpace(contact.Email),
		strings.TrimSpace(contact.Phone),
		strings.TrimSpace(contact.Address),
		strings.TrimSpace(contact.City),
		strings.TrimSpace(contact.Region),
		strings.TrimSpace(contact.Country),
		strings.TrimSpace(contact.PostalCode),
		boolToInt(contact.IsFavorite),
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func UpdateContact(db *sql.DB, id int64, contact Contact) error {
	_, err := db.Exec(`
		UPDATE contacts
		SET organization_id = ?,
			first_name = ?,
			last_name = ?,
			email = ?,
			phone = ?,
			address = ?,
			city = ?,
			region = ?,
			country = ?,
			postal_code = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		nullableInt64(contact.OrganizationID),
		strings.TrimSpace(contact.FirstName),
		strings.TrimSpace(contact.LastName),
		strings.TrimSpace(contact.Email),
		strings.TrimSpace(contact.Phone),
		strings.TrimSpace(contact.Address),
		strings.TrimSpace(contact.City),
		strings.TrimSpace(contact.Region),
		strings.TrimSpace(contact.Country),
		strings.TrimSpace(contact.PostalCode),
		id,
	)

	return err
}

func ToggleContactFavorite(db *sql.DB, id int64) error {
	_, err := db.Exec(`
		UPDATE contacts
		SET is_favorite = CASE WHEN is_favorite = 1 THEN 0 ELSE 1 END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)

	return err
}

func CountContacts(db *sql.DB) int {
	var count int

	db.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&count)

	return count
}

func CreateNote(db *sql.DB, contactID, userID int64, body string) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO notes (contact_id, user_id, body) VALUES (?, ?, ?)",
		contactID,
		userID,
		strings.TrimSpace(body),
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func ListContactNotes(db *sql.DB, contactID int64) ([]Note, error) {
	rows, err := db.Query(`
		SELECT
			n.id,
			n.contact_id,
			(c.first_name || ' ' || c.last_name) AS contact_name,
			n.user_id,
			u.name,
			n.body,
			n.created_at
		FROM notes n
		INNER JOIN contacts c ON c.id = n.contact_id
		INNER JOIN users u ON u.id = n.user_id
		WHERE n.contact_id = ?
		ORDER BY n.created_at DESC
	`, contactID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Note, error) {
		var note Note

		if err := scan(
			&note.ID,
			&note.ContactID,
			&note.ContactName,
			&note.UserID,
			&note.UserName,
			&note.Body,
			&note.CreatedAt,
		); err != nil {
			return Note{}, err
		}

		return note, nil
	}), nil
}

func ListRecentNotes(db *sql.DB, limit int) ([]Note, error) {
	rows, err := db.Query(`
		SELECT
			n.id,
			n.contact_id,
			(c.first_name || ' ' || c.last_name) AS contact_name,
			n.user_id,
			u.name,
			n.body,
			n.created_at
		FROM notes n
		INNER JOIN contacts c ON c.id = n.contact_id
		INNER JOIN users u ON u.id = n.user_id
		ORDER BY n.created_at DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows, func(scan func(...any) error) (Note, error) {
		var note Note

		if err := scan(
			&note.ID,
			&note.ContactID,
			&note.ContactName,
			&note.UserID,
			&note.UserName,
			&note.Body,
			&note.CreatedAt,
		); err != nil {
			return Note{}, err
		}

		return note, nil
	}), nil
}

func CountNotes(db *sql.DB) int {
	var count int

	db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&count)

	return count
}

func scanUser(row interface{ Scan(...any) error }) (*User, error) {
	var user User

	var verifiedAt sql.NullTime

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &verifiedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if verifiedAt.Valid {
		user.VerifiedAt = &verifiedAt.Time
	}

	return &user, nil
}

func scanContact(scan func(...any) error) (Contact, error) {
	var contact Contact

	var organizationID sql.NullInt64

	var favorite int

	if err := scan(
		&contact.ID,
		&organizationID,
		&contact.OrganizationName,
		&contact.FirstName,
		&contact.LastName,
		&contact.Email,
		&contact.Phone,
		&contact.Address,
		&contact.City,
		&contact.Region,
		&contact.Country,
		&contact.PostalCode,
		&favorite,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	); err != nil {
		return Contact{}, err
	}

	if organizationID.Valid {
		id := organizationID.Int64
		contact.OrganizationID = &id
	}

	contact.IsFavorite = favorite == 1

	return contact, nil
}

func nullableInt64(v *int64) any {
	if v == nil {
		return nil
	}

	return *v
}

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}
