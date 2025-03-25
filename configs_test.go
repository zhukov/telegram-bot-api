package tgbotapi

import (
	"testing"
)

func TestPrepareInputMediaFile(t *testing.T) {
	tests := []struct {
		name     string
		media    any
		wantFile bool
	}{
		{
			name: "InputMediaPhoto with file",
			media: InputMediaPhoto{
				BaseInputMedia: BaseInputMedia{Media: FileBytes{Bytes: []byte("test")}},
			},
			wantFile: true,
		},
		{
			name:     "InputMediaPhoto without file",
			media:    InputMediaPhoto{BaseInputMedia: BaseInputMedia{Media: FileURL("https://example.com/image.jpg")}},
			wantFile: false,
		},
		{
			name: "InputMediaVideo with file and thumb",
			media: InputMediaVideo{
				BaseInputMedia: BaseInputMedia{
					Media: FileBytes{Bytes: []byte("test")},
				},
				Thumb: &FileBytes{Bytes: []byte("thumb")},
			},
			wantFile: true,
		},
		{
			name: "InputMediaAudio with file",
			media: InputMediaAudio{
				BaseInputMedia: BaseInputMedia{Media: FileBytes{Bytes: []byte("test")}},
			},
			wantFile: true,
		},
		{
			name: "InputMediaDocument with file",
			media: InputMediaDocument{
				BaseInputMedia: BaseInputMedia{Media: FileBytes{Bytes: []byte("test")}},
			},
			wantFile: true,
		},
		{
			name: "InputMediaAnimation with file",
			media: InputMediaAnimation{
				BaseInputMedia: BaseInputMedia{Media: FileBytes{Bytes: []byte("test")}},
			},
			wantFile: true,
		},
		{
			name:     "Unsupported media type",
			media:    "string is not a valid media type",
			wantFile: false,
		},
		{
			name: "PaidMediaConfig with InputMediaPhoto",
			media: PaidMediaConfig{
				Media: []InputPaidMedia{
					{
						Type:  "photo",
						Media: FileBytes{Bytes: []byte("test")},
					},
				},
			},
			wantFile: true,
		},
		{
			name: "*PaidMediaConfig with InputMediaPhoto",
			media: &PaidMediaConfig{
				Media: []InputPaidMedia{
					{
						Type:  "photo",
						Media: FileBytes{Bytes: []byte("test")},
					},
				},
			},
			wantFile: true,
		},
	}

	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := prepareInputMediaFile(tt.media, idx)

			if tt.wantFile && len(files) == 0 {
				t.Errorf("prepareInputMediaFile() returned no files, expected at least one")
			}

			if !tt.wantFile && len(files) > 0 {
				t.Errorf("prepareInputMediaFile() returned files, expected none")
			}
		})
	}
}
