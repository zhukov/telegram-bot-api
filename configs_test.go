package tgbotapi

import (
	"fmt"
	"testing"
)

func TestPrepareInputMediaForParams(t *testing.T) {
	tests := []struct {
		name               string
		inputMedia         []InputMedia
		expectedMediaPaths []string
		expectedThumbPaths []string
	}{
		{
			name: "photo that needs upload",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FilePath("tests/image.jpg"),
					},
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{""},
		},
		{
			name: "video with thumbnail that need upload",
			inputMedia: []InputMedia{
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FilePath("tests/video.mp4"),
					},
					Thumb: FilePath("tests/image.jpg"),
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{"attach://file-0-thumb"},
		},
		{
			name: "multiple media items",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FilePath("tests/image.jpg"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FileURL("https://example.com/video.mp4"), // URL doesn't need upload
					},
				},
				&InputMediaDocument{
					BaseInputMedia: BaseInputMedia{
						Type:  "document",
						Media: FilePath("tests/audio.mp3"),
					},
					Thumb: FilePath("tests/image.jpg"),
				},
			},
			expectedMediaPaths: []string{"attach://file-0", "https://example.com/video.mp4", "attach://file-2"},
			expectedThumbPaths: []string{"", "", "attach://file-2-thumb"},
		},
		{
			name: "items that don't need upload",
			inputMedia: []InputMedia{
				&InputMediaPhoto{
					BaseInputMedia: BaseInputMedia{
						Type:  "photo",
						Media: FileID("photo123"),
					},
				},
				&InputMediaVideo{
					BaseInputMedia: BaseInputMedia{
						Type:  "video",
						Media: FileURL("https://example.com/video.mp4"),
					},
					Thumb: FileID("thumb123"),
				},
			},
			expectedMediaPaths: []string{"photo123", "https://example.com/video.mp4"},
			expectedThumbPaths: []string{"", "thumb123"},
		},
		{
			name: "paid media",
			inputMedia: []InputMedia{
				&InputPaidMedia{
					Type:  "video",
					Media: FilePath("tests/video.mp4"),
					Thumb: FilePath("tests/image.jpg"),
				},
			},
			expectedMediaPaths: []string{"attach://file-0"},
			expectedThumbPaths: []string{"attach://file-0-thumb"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prepareInputMediaForParams(tt.inputMedia)

			if len(result) != len(tt.inputMedia) {
				t.Errorf("Expected result length %d, got %d", len(tt.inputMedia), len(result))
			}

			for i, media := range result {
				mediaPath := media.getMedia().SendData()
				if mediaPath != tt.expectedMediaPaths[i] {
					t.Errorf("Media path at index %d: expected %s, got %s", i, tt.expectedMediaPaths[i], mediaPath)
				}

				thumb := media.getThumb()
				var thumbPath string
				if thumb != nil {
					thumbPath = thumb.SendData()
				}

				expectedThumb := tt.expectedThumbPaths[i]
				if thumbPath != expectedThumb {
					t.Errorf("Thumb path at index %d: expected %s, got %s", i, expectedThumb, thumbPath)
				}
			}

			if &result[0] == &tt.inputMedia[0] {
				t.Error("Result should be a deep copy, not a reference to the original slice")
			}
		})
	}
}

func TestPrepareInputMediaForFiles(t *testing.T) {
	inputMedia := []InputMedia{
		&InputMediaPhoto{
			BaseInputMedia: BaseInputMedia{
				Type:  "photo",
				Media: FilePath("tests/image.jpg"),
			},
		},
		&InputMediaVideo{
			BaseInputMedia: BaseInputMedia{
				Type:  "video",
				Media: FileURL("https://example.com/video.mp4"), // This doesn't need upload
			},
		},
		&InputMediaVideo{
			BaseInputMedia: BaseInputMedia{
				Type:  "video",
				Media: FilePath("tests/video.mp4"),
			},
			Thumb: FilePath("tests/image.jpg"),
		},
		&InputPaidMedia{
			Type:  "video",
			Media: FilePath("tests/video.mp4"),
			Thumb: FilePath("tests/image.jpg"),
		},
	}

	files := prepareInputMediaForFiles(inputMedia)

	expectedFiles := 5
	if len(files) != expectedFiles {
		t.Errorf("Expected %d files, got %d", expectedFiles, len(files))
	}

	expectedNames := map[string]bool{
		"file-0":       false,
		"file-2":       false,
		"file-2-thumb": false,
		"file-3":       false,
		"file-3-thumb": false,
	}

	for _, file := range files {
		expectedNames[file.Name] = true
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected file %s not found in result: %v", name, files)
		}
	}
}

func TestCloneMediaSlice(t *testing.T) {
	// Test that the cloning function creates proper deep copies
	inputMedia := []InputMedia{
		&InputMediaPhoto{
			BaseInputMedia: BaseInputMedia{
				Type:    "photo",
				Media:   FilePath("photo.jpg"),
				Caption: "Original photo caption",
			},
		},
		&InputMediaVideo{
			BaseInputMedia: BaseInputMedia{
				Type:    "video",
				Media:   FilePath("video.mp4"),
				Caption: "Original video caption",
			},
			Thumb:    FilePath("thumb.jpg"),
			Duration: 30,
		},
		&InputMediaAnimation{
			BaseInputMedia: BaseInputMedia{
				Type:    "animation",
				Media:   FilePath("anim.gif"),
				Caption: "Original animation caption",
			},
			Thumb: FilePath("anim-thumb.jpg"),
		},
		&InputMediaAudio{
			BaseInputMedia: BaseInputMedia{
				Type:    "audio",
				Media:   FilePath("audio.mp3"),
				Caption: "Original audio caption",
			},
			Thumb: FilePath("audio-thumb.jpg"),
		},
		&InputMediaDocument{
			BaseInputMedia: BaseInputMedia{
				Type:    "document",
				Media:   FilePath("doc.pdf"),
				Caption: "Original document caption",
			},
			Thumb: FilePath("doc-thumb.jpg"),
		},
		&InputPaidMedia{
			Type:  "paid_media",
			Media: FilePath("paid.jpg"),
			Thumb: FilePath("paid-thumb.jpg"),
		},
	}

	cloned := cloneMediaSlice(inputMedia)

	if len(cloned) != len(inputMedia) {
		t.Fatalf("Expected cloned slice length %d, got %d", len(inputMedia), len(cloned))
	}

	// Test that we have a deep copy (modifying one doesn't affect the other)
	for i := 0; i < len(inputMedia); i++ {
		if fmt.Sprintf("%T", cloned[i]) != fmt.Sprintf("%T", inputMedia[i]) {
			t.Errorf("Type mismatch at index %d: expected %T, got %T",
				i, inputMedia[i], cloned[i])
		}

		// Test deep copy by modifying caption in cloned and checking original
		switch media := cloned[i].(type) {
		case *InputMediaPhoto:
			originalMedia := inputMedia[i].(*InputMediaPhoto)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaVideo:
			originalMedia := inputMedia[i].(*InputMediaVideo)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaAnimation:
			originalMedia := inputMedia[i].(*InputMediaAnimation)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaAudio:
			originalMedia := inputMedia[i].(*InputMediaAudio)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		case *InputMediaDocument:
			originalMedia := inputMedia[i].(*InputMediaDocument)
			media.Caption = "Modified caption"
			if originalMedia.Caption == media.Caption {
				t.Errorf("Expected deep copy, but caption was modified in original")
			}
		}
	}
}
