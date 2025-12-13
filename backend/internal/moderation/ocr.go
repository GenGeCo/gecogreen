package moderation

import (
	"context"
	"regexp"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
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
type OCRService struct {
	client *vision.ImageAnnotatorClient
}

// NewOCRService creates a new OCR service using Google Cloud Vision
// Requires GOOGLE_APPLICATION_CREDENTIALS env var to be set
func NewOCRService(ctx context.Context) (*OCRService, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	return &OCRService{client: client}, nil
}

// Close closes the client connection
func (s *OCRService) Close() error {
	return s.client.Close()
}

// AnalyzeImageFromURL analyzes an image from a URL
func (s *OCRService) AnalyzeImageFromURL(ctx context.Context, imageURL string) (*OCRResult, error) {
	image := vision.NewImageFromURI(imageURL)
	return s.analyzeImage(ctx, image)
}

// AnalyzeImageFromBytes analyzes an image from bytes
func (s *OCRService) AnalyzeImageFromBytes(ctx context.Context, imageBytes []byte) (*OCRResult, error) {
	image := &visionpb.Image{Content: imageBytes}
	return s.analyzeImage(ctx, image)
}

func (s *OCRService) analyzeImage(ctx context.Context, image *visionpb.Image) (*OCRResult, error) {
	// Detect text in image
	annotations, err := s.client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	result := &OCRResult{
		Confidence: 0,
	}

	if len(annotations) == 0 {
		return result, nil
	}

	// First annotation contains all detected text
	fullText := annotations[0].Description
	result.DetectedText = fullText

	// Check for phone numbers
	result.DetectedPhone = containsPhoneNumber(fullText)

	// Check for emails
	result.DetectedEmail = containsEmail(fullText)

	// Check for URLs
	result.DetectedURL = containsURL(fullText)

	// Calculate if suspicious
	result.IsSuspicious = result.DetectedPhone || result.DetectedEmail || result.DetectedURL

	// Confidence based on bounding box coverage
	if annotations[0].BoundingPoly != nil {
		result.Confidence = 0.9 // High confidence if we have bounding info
	} else {
		result.Confidence = 0.7
	}

	return result, nil
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
