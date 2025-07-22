package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_NewUser(t *testing.T) {
	tests := []struct {
		testName string
		id       string
		email    string
		userName string
		picture  string
		want     *User
		wantErr  bool
	}{
		{
			testName: "正常なユーザー作成",
			id:       "test-id",
			email:    "test@example.com",
			userName: "Test User",
			picture:  "https://example.com/picture.jpg",
			want: &User{
				ID:      "test-id",
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "https://example.com/picture.jpg",
			},
			wantErr: false,
		},
		{
			testName: "空のIDでエラー",
			id:       "",
			email:    "test@example.com",
			userName: "Test User",
			picture:  "https://example.com/picture.jpg",
			wantErr:  true,
		},
		{
			testName: "空のEmailでエラー",
			id:       "test-id",
			email:    "",
			userName: "Test User",
			picture:  "https://example.com/picture.jpg",
			wantErr:  true,
		},
		{
			testName: "無効なEmailでエラー",
			id:       "test-id",
			email:    "invalid-email",
			userName: "Test User",
			picture:  "https://example.com/picture.jpg",
			wantErr:  true,
		},
		{
			testName: "空のNameでエラー",
			id:       "test-id",
			email:    "test@example.com",
			userName: "",
			picture:  "https://example.com/picture.jpg",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := NewUser(tt.id, tt.email, tt.userName, tt.picture)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.Picture, got.Picture)
				assert.WithinDuration(t, time.Now(), got.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
			}
		})
	}
}

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		testName string
		user     *User
		want     bool
	}{
		{
			testName: "有効なユーザー",
			user: &User{
				ID:      "test-id",
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "https://example.com/picture.jpg",
			},
			want: true,
		},
		{
			testName: "空のIDで無効",
			user: &User{
				ID:      "",
				Email:   "test@example.com",
				Name:    "Test User",
				Picture: "https://example.com/picture.jpg",
			},
			want: false,
		},
		{
			testName: "空のEmailで無効",
			user: &User{
				ID:      "test-id",
				Email:   "",
				Name:    "Test User",
				Picture: "https://example.com/picture.jpg",
			},
			want: false,
		},
		{
			testName: "無効なEmailで無効",
			user: &User{
				ID:      "test-id",
				Email:   "invalid-email",
				Name:    "Test User",
				Picture: "https://example.com/picture.jpg",
			},
			want: false,
		},
		{
			testName: "空のNameで無効",
			user: &User{
				ID:      "test-id",
				Email:   "test@example.com",
				Name:    "",
				Picture: "https://example.com/picture.jpg",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := tt.user.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	user := &User{
		ID:        "test-id",
		Email:     "test@example.com",
		Name:      "Test User",
		Picture:   "https://example.com/picture.jpg",
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	newName := "Updated User"
	newPicture := "https://example.com/new-picture.jpg"

	err := user.UpdateProfile(newName, newPicture)
	assert.NoError(t, err)
	assert.Equal(t, newName, user.Name)
	assert.Equal(t, newPicture, user.Picture)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
}
