package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ekholme/flexcreek"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite" // Import the sqlite driver for database/sql
)

// mustOpenDB opens an in-memory SQLite database for testing.
// It panics on error, which is acceptable for a test setup function.
func mustOpenDB(t *testing.T) *sql.DB {
	t.Helper()
	// Using ":memory:" creates a unique, private in-memory database for each test.
	// This is ideal for parallel tests as it prevents state from leaking between them.
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	return db
}

// createSchema creates the necessary tables for the tests.
func createSchema(t *testing.T, db *sql.DB) {
	t.Helper()
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		hashed_password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
}

// newTestUserService creates a new UserService with a clean in-memory database.
// It returns the service and a teardown function to close the DB.
func newTestUserService(t *testing.T) (flexcreek.UserService, func()) {
	t.Helper()

	db := mustOpenDB(t)
	createSchema(t, db)

	// Teardown function to clean up after tests.
	teardown := func() {
		db.Close()
	}

	// Since each test gets a fresh in-memory database, we don't need to
	// clear the tables manually.
	return NewUserService(db), teardown
}

// hashPassword is a test helper to bcrypt a password string.
func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hash)
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()
	service, teardown := newTestUserService(t)
	defer teardown()

	ctx := context.Background()
	user := &flexcreek.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		HashedPw:  hashPassword(t, "password123"),
	}

	id, err := service.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser() error = %v, want nil", err)
	}

	if id != 1 {
		t.Errorf("CreateUser() id = %d, want 1", id)
	}

	// Verify user was actually created by fetching it
	createdUser, err := service.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("GetUserByID() after create error = %v", err)
	}

	if createdUser.Email != user.Email {
		t.Errorf("created user email = %s, want %s", createdUser.Email, user.Email)
	}
}

func TestUserService_GetUser(t *testing.T) {
	t.Parallel()
	service, teardown := newTestUserService(t)
	defer teardown()

	ctx := context.Background()
	hashedPw := hashPassword(t, "password123")
	user := &flexcreek.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		HashedPw:  hashedPw,
	}

	id, err := service.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user for testing: %v", err)
	}

	t.Run("GetUserByID", func(t *testing.T) {
		foundUser, err := service.GetUserByID(ctx, id)
		if err != nil {
			t.Fatalf("GetUserByID() error = %v, want nil", err)
		}
		if foundUser.ID != id || foundUser.Email != user.Email {
			t.Errorf("GetUserByID() returned wrong user data")
		}
	})

	t.Run("GetUserByID_NotFound", func(t *testing.T) {
		_, err := service.GetUserByID(ctx, 999)
		if err != sql.ErrNoRows {
			t.Errorf("GetUserByID() with non-existent ID, error = %v, want %v", err, sql.ErrNoRows)
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		foundUser, err := service.GetUserByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("GetUserByEmail() error = %v, want nil", err)
		}
		if foundUser.ID != id || foundUser.Email != user.Email {
			t.Errorf("GetUserByEmail() returned wrong user data")
		}
	})

	t.Run("GetUserByEmail_NotFound", func(t *testing.T) {
		_, err := service.GetUserByEmail(ctx, "notfound@example.com")
		if err != sql.ErrNoRows {
			t.Errorf("GetUserByEmail() with non-existent email, error = %v, want %v", err, sql.ErrNoRows)
		}
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	t.Parallel()
	service, teardown := newTestUserService(t)
	defer teardown()

	ctx := context.Background()

	// Create two users
	user1 := &flexcreek.User{FirstName: "Alice", LastName: "Smith", Email: "alice@example.com", HashedPw: "hash1"}
	user2 := &flexcreek.User{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", HashedPw: "hash2"}
	_, _ = service.CreateUser(ctx, user1)
	_, _ = service.CreateUser(ctx, user2)

	users, err := service.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("GetAllUsers() error = %v, want nil", err)
	}

	if len(users) != 2 {
		t.Errorf("GetAllUsers() count = %d, want 2", len(users))
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Parallel()
	service, teardown := newTestUserService(t)
	defer teardown()

	ctx := context.Background()
	user := &flexcreek.User{FirstName: "Initial", LastName: "Name", Email: "update@example.com", HashedPw: "initial_hash"}
	id, err := service.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user for update test: %v", err)
	}
	// Fetch the initial user to compare timestamps later.
	initialUser, err := service.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to fetch user for update test: %v", err)
	}

	// Add a small delay to ensure CURRENT_TIMESTAMP will be different.
	time.Sleep(1 * time.Second)

	updatedUser := &flexcreek.User{ID: id, FirstName: "Updated", LastName: "User", Email: "updated.user@example.com", HashedPw: "updated_hash"}
	err = service.UpdateUser(ctx, updatedUser)
	if err != nil {
		t.Fatalf("UpdateUser() error = %v, want nil", err)
	}

	fetchedUser, err := service.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("GetUserByID() after update error = %v", err)
	}

	if fetchedUser.FirstName != "Updated" || fetchedUser.Email != "updated.user@example.com" {
		t.Errorf("user was not updated correctly. Got name: %s, email: %s", fetchedUser.FirstName, fetchedUser.Email)
	}

	if !fetchedUser.UpdatedAt.After(initialUser.UpdatedAt) {
		t.Errorf("expected UpdatedAt to be after initial value; initial=%v, updated=%v", initialUser.UpdatedAt, fetchedUser.UpdatedAt)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Parallel()
	service, teardown := newTestUserService(t)
	defer teardown()

	ctx := context.Background()
	user := &flexcreek.User{FirstName: "Delete", LastName: "Me", Email: "delete.me@example.com", HashedPw: "some_hash"}
	id, err := service.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user for delete test: %v", err)
	}

	err = service.DeleteUser(ctx, id)
	if err != nil {
		t.Fatalf("DeleteUser() error = %v, want nil", err)
	}

	_, err = service.GetUserByID(ctx, id)
	if err != sql.ErrNoRows {
		t.Errorf("GetUserByID() after delete, error = %v, want %v", err, sql.ErrNoRows)
	}
}

func TestUserService_Login(t *testing.T) {
	t.Parallel()
	service, teardown := newTestUserService(t)
	defer teardown()

	ctx := context.Background()
	email := "login.test@example.com"
	password := "S3cureP@ssw0rd!"
	hashedPw := hashPassword(t, password)

	user := &flexcreek.User{FirstName: "Login", LastName: "Test", Email: email, HashedPw: hashedPw}
	_, err := service.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user for login test: %v", err)
	}

	t.Run("Successful Login", func(t *testing.T) {
		loggedInUser, err := service.Login(ctx, email, password)
		if err != nil {
			t.Fatalf("Login() with correct credentials error = %v, want nil", err)
		}
		if loggedInUser == nil {
			t.Fatal("Login() returned nil user on success")
		}
		if loggedInUser.Email != email {
			t.Errorf("Login() returned user with wrong email. got %s, want %s", loggedInUser.Email, email)
		}
	})

	t.Run("Wrong Password", func(t *testing.T) {
		_, err := service.Login(ctx, email, "wrongpassword")
		if err != flexcreek.ErrInvalidCredentials {
			t.Errorf("Login() with wrong password, error = %v, want %v", err, flexcreek.ErrInvalidCredentials)
		}
	})

	t.Run("Non-existent Email", func(t *testing.T) {
		_, err := service.Login(ctx, "no.such.user@example.com", password)
		if err != flexcreek.ErrInvalidCredentials {
			t.Errorf("Login() with non-existent email, error = %v, want %v", err, flexcreek.ErrInvalidCredentials)
		}
	})
}
