package userservice

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"example.com/m/models"
	"example.com/m/repository"
	"example.com/m/utils"
)

// ---------- фейковый репозиторий ----------

// fakeRepo хранит профили в map и считает вызовы.
// Get/Create отдают и сохраняют КОПИИ, как настоящая БД: если сервис
// изменит объект и вернёт ошибку, "база" не должна поменяться.
type fakeRepo struct {
	data map[uuid.UUID]models.UserProfile

	getCalls, createCalls, updateCalls, deleteCalls int

	// если не nil, соответствующий метод вернёт эту ошибку
	errGet, errCreate, errUpdate, errDelete error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{data: make(map[uuid.UUID]models.UserProfile)}
}

func (f *fakeRepo) GetUserProfileByID(_ context.Context, id uuid.UUID) (*models.UserProfile, error) {
	f.getCalls++
	if f.errGet != nil {
		return nil, f.errGet
	}
	p, ok := f.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &p, nil // p — копия
}

func (f *fakeRepo) CreateUserProfile(_ context.Context, p *models.UserProfile) error {
	f.createCalls++
	if f.errCreate != nil {
		return f.errCreate
	}
	for _, other := range f.data {
		if other.Username == p.Username {
			return repository.ErrAlreadyExists
		}
	}
	f.data[p.UserID] = *p
	return nil
}

func (f *fakeRepo) UpdateUserProfile(_ context.Context, p *models.UserProfile) error {
	f.updateCalls++
	if f.errUpdate != nil {
		return f.errUpdate
	}
	if _, ok := f.data[p.UserID]; !ok {
		return repository.ErrNotFound
	}
	f.data[p.UserID] = *p
	return nil
}

func (f *fakeRepo) DeleteUserProfile(_ context.Context, id uuid.UUID) error {
	f.deleteCalls++
	if f.errDelete != nil {
		return f.errDelete
	}
	if _, ok := f.data[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.data, id)
	return nil
}

// ---------- хелперы ----------

func strPtr(s string) *string { return &s }

// seed кладёт в "базу" готового пользователя через сам сервис.
func seed(t *testing.T, svc *UserService) uuid.UUID {
	t.Helper()
	id := uuid.New()
	_, err := svc.CreateUserProfile(context.Background(), id, "azat", "Азат", "Шакиев", "+79625677812")
	if err != nil {
		t.Fatalf("seed: unexpected error: %v", err)
	}
	return id
}

// ---------- Get ----------

func TestGetUserProfileByID(t *testing.T) {
	ctx := context.Background()

	t.Run("nil id", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)

		_, err := svc.GetUserProfileByID(ctx, uuid.Nil)

		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("want ErrInvalidUserID, got %v", err)
		}
		if repo.getCalls != 0 {
			t.Errorf("repo must not be called on invalid id, calls=%d", repo.getCalls)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := NewUserService(newFakeRepo())

		_, err := svc.GetUserProfileByID(ctx, uuid.New())

		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("want ErrNotFound, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		svc := NewUserService(newFakeRepo())
		id := seed(t, svc)

		p, err := svc.GetUserProfileByID(ctx, id)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.UserID != id || p.Username != "azat" {
			t.Errorf("unexpected profile: %+v", p)
		}
	})
}

// ---------- Create ----------

func TestCreateUserProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("nil id", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)

		_, err := svc.CreateUserProfile(ctx, uuid.Nil, "azat", "Азат", "", "+79625677812")

		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("want ErrInvalidUserID, got %v", err)
		}
		if repo.createCalls != 0 {
			t.Error("repo must not be called")
		}
	})

	t.Run("validation errors never reach repo", func(t *testing.T) {
		tests := []struct {
			name                         string
			username, first, last, phone string
		}{
			{"bad username", "bad name!", "Азат", "", "+79625677812"},
			{"empty username", "", "Азат", "", "+79625677812"},
			{"empty name", "azat", "", "", "+79625677812"},
			{"name with digits", "azat", "Аз4т", "", "+79625677812"},
			{"empty phone", "azat", "Азат", "", ""},
			{"short phone", "azat", "Азат", "", "+7962"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := newFakeRepo()
				svc := NewUserService(repo)

				_, err := svc.CreateUserProfile(ctx, uuid.New(), tt.username, tt.first, tt.last, tt.phone)

				if !errors.Is(err, utils.ErrValidation) {
					t.Fatalf("want ErrValidation, got %v", err)
				}
				if repo.createCalls != 0 {
					t.Error("repo must not be called when validation fails")
				}
			})
		}
	})

	t.Run("success returns and stores profile", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := uuid.New()

		p, err := svc.CreateUserProfile(ctx, id, "azat", "Азат", "Шакиев", "+79625677812")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.UserID != id {
			t.Errorf("returned wrong id")
		}
		if _, ok := repo.data[id]; !ok {
			t.Error("profile was not saved in repo")
		}
	})

	t.Run("duplicate username", func(t *testing.T) {
		svc := NewUserService(newFakeRepo())
		seed(t, svc)

		_, err := svc.CreateUserProfile(ctx, uuid.New(), "azat", "Другой", "", "+79990000000")

		if !errors.Is(err, repository.ErrAlreadyExists) {
			t.Fatalf("want ErrAlreadyExists, got %v", err)
		}
	})

	t.Run("repo error is wrapped, not swallowed", func(t *testing.T) {
		repo := newFakeRepo()
		boom := errors.New("db is down")
		repo.errCreate = boom
		svc := NewUserService(repo)

		_, err := svc.CreateUserProfile(ctx, uuid.New(), "azat", "Азат", "", "+79625677812")

		if !errors.Is(err, boom) {
			t.Fatalf("want wrapped boom, got %v", err)
		}
	})
}

// ---------- Update ----------

func TestUpdateUserProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("nil input", func(t *testing.T) {
		svc := NewUserService(newFakeRepo())

		_, err := svc.UpdateUserProfile(ctx, nil)

		if !errors.Is(err, ErrNilUpdate) {
			t.Fatalf("want ErrNilUpdate, got %v", err)
		}
	})

	t.Run("nil id", func(t *testing.T) {
		svc := NewUserService(newFakeRepo())

		_, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: uuid.Nil, Name: strPtr("Новое")})

		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("want ErrInvalidUserID, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		svc := NewUserService(newFakeRepo())

		_, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: uuid.New(), Name: strPtr("Новое")})

		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("want ErrNotFound, got %v", err)
		}
	})

	t.Run("partial update changes only passed field", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := seed(t, svc)

		p, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: id, Name: strPtr("Азатик")})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Name != "Азатик" {
			t.Errorf("name not changed: %q", p.Name)
		}
		// nil-поля не должны затираться
		if p.Username != "azat" || p.Surname != "Шакиев" || p.Phone != "+79625677812" {
			t.Errorf("untouched fields changed: %+v", p)
		}
		if repo.data[id].Name != "Азатик" {
			t.Error("change was not saved in repo")
		}
	})

	t.Run("update bumps UpdatedAt", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := seed(t, svc)
		before := repo.data[id].UpdatedAt

		p, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: id, Name: strPtr("Азатик")})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !p.UpdatedAt.After(before) {
			t.Errorf("UpdatedAt not bumped: before=%v after=%v", before, p.UpdatedAt)
		}
	})

	t.Run("no real changes do not hit repo.Update", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := seed(t, svc)

		// все поля nil
		_, err1 := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: id})
		// поле передано, но значение то же самое
		_, err2 := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: id, Name: strPtr("Азат")})

		if err1 != nil || err2 != nil {
			t.Fatalf("unexpected errors: %v, %v", err1, err2)
		}
		if repo.updateCalls != 0 {
			t.Errorf("repo.Update must not be called, calls=%d", repo.updateCalls)
		}
	})

	t.Run("invalid value: error and DB stays untouched", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := seed(t, svc)

		// имя валидное, логин нет: имя в копии уже изменено, но сохраняться ничего не должно
		_, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{
			UserID:   id,
			Name:     strPtr("Новое"),
			Username: strPtr("bad name!"),
		})

		if !errors.Is(err, utils.ErrValidation) {
			t.Fatalf("want ErrValidation, got %v", err)
		}
		if repo.updateCalls != 0 {
			t.Error("repo.Update must not be called after a validation error")
		}
		if got := repo.data[id].Name; got != "Азат" {
			t.Errorf("DB changed despite error: name=%q", got)
		}
	})

	t.Run("clear description with empty string", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := seed(t, svc)
		if _, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: id, Description: strPtr("привет")}); err != nil {
			t.Fatal(err)
		}

		p, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: id, Description: strPtr("")})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Description != "" {
			t.Errorf("description not cleared: %q", p.Description)
		}
	})

	t.Run("repo update error is returned", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := seed(t, svc)
		boom := errors.New("db is down")
		repo.errUpdate = boom

		_, err := svc.UpdateUserProfile(ctx, &UpdateUserProfileInput{UserID: id, Name: strPtr("Азатик")})

		if !errors.Is(err, boom) {
			t.Fatalf("want wrapped boom, got %v", err)
		}
	})
}

// ---------- Delete ----------

func TestDeleteUserProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("nil id", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)

		err := svc.DeleteUserProfile(ctx, uuid.Nil)

		if !errors.Is(err, ErrInvalidUserID) {
			t.Fatalf("want ErrInvalidUserID, got %v", err)
		}
		if repo.deleteCalls != 0 {
			t.Error("repo must not be called")
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := NewUserService(newFakeRepo())

		err := svc.DeleteUserProfile(ctx, uuid.New())

		if !errors.Is(err, repository.ErrNotFound) {
			t.Fatalf("want ErrNotFound, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewUserService(repo)
		id := seed(t, svc)

		if err := svc.DeleteUserProfile(ctx, id); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := repo.data[id]; ok {
			t.Error("profile still in repo")
		}
	})
}
