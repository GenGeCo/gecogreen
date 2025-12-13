package moderation

import (
	"context"
	"regexp"
	"strings"
)

// OCRResult holds the results of text detection
type OCRResult struct {
	DetectedText  string  `json:"detected_text"`
	DetectedPhone bool    `json:"detected_phone"`
	DetectedEmail bool    `json:"detected_email"`
	DetectedURL   bool    `json:"detected_url"`
	IsSuspicious  bool    `json:"is_suspicious"`
	Confidence    float64 `json:"confidence"`
}

// OCRService handles text detection in images
// Currently a stub - Google Cloud Vision can be added later
type OCRService struct {
	enabled bool
}

// NewOCRService creates a new OCR service
// Currently returns a disabled service - OCR can be enabled later with Google Cloud Vision
func NewOCRService(ctx context.Context) (*OCRService, error) {
	// OCR is disabled by default until Google Cloud Vision is configured
	// To enable: set GOOGLE_APPLICATION_CREDENTIALS and uncomment the vision client code
	return &OCRService{enabled: false}, nil
}

// Close closes the client connection
func (s *OCRService) Close() error {
	return nil
}

// IsEnabled returns whether OCR is actually functional
func (s *OCRService) IsEnabled() bool {
	return s.enabled
}

// AnalyzeImageFromURL analyzes an image from a URL
// Currently returns empty result - OCR disabled
func (s *OCRService) AnalyzeImageFromURL(ctx context.Context, imageURL string) (*OCRResult, error) {
	if !s.enabled {
		return &OCRResult{Confidence: 0}, nil
	}
	// TODO: Implement with Google Cloud Vision when Go 1.24+ is available
	return &OCRResult{Confidence: 0}, nil
}

// AnalyzeImageFromBytes analyzes an image from bytes
// Currently returns empty result - OCR disabled
func (s *OCRService) AnalyzeImageFromBytes(ctx context.Context, imageBytes []byte) (*OCRResult, error) {
	if !s.enabled {
		return &OCRResult{Confidence: 0}, nil
	}
	// TODO: Implement with Google Cloud Vision when Go 1.24+ is available
	return &OCRResult{Confidence: 0}, nil
}

// AnalyzeText analyzes text for suspicious content (phone, email, URLs)
// This works without Google Cloud Vision - useful for analyzing text fields
func AnalyzeText(text string) *OCRResult {
	result := &OCRResult{
		DetectedText: text,
		Confidence:   1.0,
	}

	result.DetectedPhone = containsPhoneNumber(text)
	result.DetectedEmail = containsEmail(text)
	result.DetectedURL = containsURL(text)
	result.IsSuspicious = result.DetectedPhone || result.DetectedEmail || result.DetectedURL

	return result
}

// Phone number patterns (Italian and international)
var phonePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\+39\s*\d{2,3}[\s.-]?\d{3}[\s.-]?\d{4}`),  // +39 xxx xxx xxxx
	regexp.MustCompile(`\b0\d{1,3}[\s.-]?\d{6,8}\b`),               // Italian landline
	regexp.MustCompile(`\b3\d{2}[\s.-]?\d{3}[\s.-]?\d{4}\b`),       // Italian mobile
	regexp.MustCompile(`\b\d{3}[\s.-]?\d{3}[\s.-]?\d{4}\b`),        // Generic format
}

func containsPhoneNumber(text string) bool {
	text = strings.ReplaceAll(text, "\n", " ")
	for _, pattern := range phonePatterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}

var emailPattern = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

func containsEmail(text string) bool {
	return emailPattern.MatchString(text)
}

var urlPatterns = []*regexp.Regexp{
	regexp.MustCompile(`https?://[^\s]+`),
	regexp.MustCompile(`www\.[^\s]+`),
	regexp.MustCompile(`[a-zA-Z0-9-]+\.(com|it|eu|org|net|io)\b`),
}

func containsURL(text string) bool {
	for _, pattern := range urlPatterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}
