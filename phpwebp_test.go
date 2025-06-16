package phpwebp_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestPHPWebP(t *testing.T) {
	phpVersion := os.Getenv("PHP_VERSION")
	if phpVersion == "" {
		phpVersion = "8.1"
	}
	
	t.Run(fmt.Sprintf("PHP_%s", phpVersion), func(t *testing.T) {
			pool, err := dockertest.NewPool("")
			require.NoError(t, err)
			
			pool.MaxWait = 120 * time.Second
			
			pwd, err := os.Getwd()
			require.NoError(t, err)
			
			wwwrootPath := filepath.Join(pwd, "wwwroot")
			
		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "php",
			Tag:        phpVersion + "-apache",
			Cmd:        []string{"/bin/bash", "-c", "a2enmod rewrite && /usr/local/bin/apache2-foreground"},
			Mounts: []string{
				wwwrootPath + ":/var/www/html",
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
			require.NoError(t, err)
			
			defer pool.Purge(resource)
			
			hostPort := resource.GetPort("80/tcp")
			baseURL := fmt.Sprintf("http://localhost:%s", hostPort)
			
			err = pool.Retry(func() error {
				resp, err := http.Get(baseURL + "/testing/regular.jpg")
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("status code: %d", resp.StatusCode)
				}
				return nil
			})
			require.NoError(t, err)
			
			// WebP対応エージェント(Accept: image/webp あり)のテスト
			t.Run("WebP_Supported_Agent", func(t *testing.T) {
				// 従来フォーマットからWebPに正常に変換されるケース
				testCases := []struct {
					path         string
					expectedType string
				}{
					{"/testing/regular.jpg", "image/webp"},
					{"/testing/regular.png", "image/webp"},
					{"/testing/animation.gif", "image/webp"},
				}
				
				for _, tc := range testCases {
					t.Run(tc.path, func(t *testing.T) {
						testValidImageRequest(t, baseURL, tc.path, "image/webp,*/*", 
							http.StatusOK, tc.expectedType, tc.expectedType)
					})
				}
				
				// WebPに正常に変換されないケース(非対応のBMP)
				t.Run("/testing/bmp.jpg", func(t *testing.T) {
					testBMPImageRequest(t, baseURL, "/testing/bmp.jpg", "image/webp,*/*",
						http.StatusOK)
				})
				
				// WebPに正常に変換されないケース(非対応のCMYK画像)
				t.Run("/testing/cmyk.jpg", func(t *testing.T) {
					testValidImageRequest(t, baseURL, "/testing/cmyk.jpg", "image/webp,*/*",
						http.StatusOK, "image/jpeg", "image/jpeg")
				})
				
				// WebPに直接アクセスするケース
				webpCases := []string{"/testing/lossy.webp", "/testing/lossless.webp"}
				for _, path := range webpCases {
					t.Run(path, func(t *testing.T) {
						testValidImageRequest(t, baseURL, path, "image/webp,*/*",
							http.StatusOK, "image/webp", "image/webp")
					})
				}
			})
			
			// WebP非対応エージェント(Accept: image/webp なし)のテスト
			t.Run("WebP_Unsupported_Agent", func(t *testing.T) {
				// 従来フォーマットの画像にアクセスするケース
				testCases := []struct {
					path         string
					expectedType string
				}{
					{"/testing/regular.jpg", "image/jpeg"},
					{"/testing/regular.png", "image/png"},
					{"/testing/animation.gif", "image/gif"},
				}
				
				for _, tc := range testCases {
					t.Run(tc.path, func(t *testing.T) {
						testValidImageRequest(t, baseURL, tc.path, "*/*",
							http.StatusOK, tc.expectedType, tc.expectedType)
					})
				}
				
				// WebP画像にアクセスするケース
				webpCases := []struct {
					path         string
					expectedType string
				}{
					{"/testing/lossy.webp", "image/png"},
					{"/testing/lossless.webp", "image/png"},
				}
				
				for _, tc := range webpCases {
					t.Run(tc.path, func(t *testing.T) {
						testValidImageRequest(t, baseURL, tc.path, "*/*",
							http.StatusOK, tc.expectedType, tc.expectedType)
					})
				}
				
			// WebP非対応エージェントとしてアクセスするがPNGに変換できないケース(実体がBMP)
			t.Run("/testing/bmp.webp", func(t *testing.T) {
				if phpVersion == "5.4" || phpVersion == "5.5" {
					testValidImageRequest(t, baseURL, "/testing/bmp.webp", "*/*",
						http.StatusOK, "application/octet-stream", "image/bmp")
				} else {
					// PHP 5.6以降はimage/bmpまたはimage/x-ms-bmpを返す可能性がある
					req, err := http.NewRequest("GET", baseURL+"/testing/bmp.webp", nil)
					require.NoError(t, err)
					req.Header.Set("Accept", "*/*")
					
					client := &http.Client{}
					resp, err := client.Do(req)
					require.NoError(t, err)
					defer resp.Body.Close()
					
					assert.Equal(t, http.StatusOK, resp.StatusCode)
					
					contentType := resp.Header.Get("Content-Type")
					assert.True(t, contentType == "image/bmp" || contentType == "image/x-ms-bmp",
						"Content-Type should be image/bmp or image/x-ms-bmp, got %s", contentType)
					
					body, err := io.ReadAll(resp.Body)
					require.NoError(t, err)
					
					actualMimeType := getImageMimeType(body)
					assert.Equal(t, "image/bmp", actualMimeType,
						"Actual content should be image/bmp")
				}
			})
			})
			
			// 404 Not Foundのテスト
			t.Run("404_Not_Found", func(t *testing.T) {
				testCases := []struct {
					name   string
					accept string
				}{
					{"WebP_Supported_Agent", "image/webp,*/*"},
					{"WebP_Unsupported_Agent", "*/*"},
				}
				
				for _, tc := range testCases {
					t.Run(tc.name, func(t *testing.T) {
						testInvalidRequest(t, baseURL, "/testing/notfound.jpg", tc.accept,
							http.StatusNotFound)
					})
				}
		})
	})
}

func testValidImageRequest(t *testing.T, baseURL, path, acceptHeader string, 
	expectedStatus int, expectedMimeType, expectedDataMimeType string) {
	
	req, err := http.NewRequest("GET", baseURL+path, nil)
	require.NoError(t, err)
	
	req.Header.Set("Accept", acceptHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"%s should return status code %d", path, expectedStatus)
	
	assert.Equal(t, expectedMimeType, resp.Header.Get("Content-Type"),
		"%s should return Content-Type: %s", path, expectedMimeType)
	
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	
	// Verify actual image type
	actualMimeType := getImageMimeType(body)
	assert.Equal(t, expectedDataMimeType, actualMimeType,
		"%s actual content should be %s", path, expectedDataMimeType)
}

func testInvalidImageRequest(t *testing.T, baseURL, path, acceptHeader string,
	expectedStatus int, expectedMimeType string) {
	
	req, err := http.NewRequest("GET", baseURL+path, nil)
	require.NoError(t, err)
	
	req.Header.Set("Accept", acceptHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"%s should return status code %d", path, expectedStatus)
	
	if expectedMimeType != "" {
		assert.Equal(t, expectedMimeType, resp.Header.Get("Content-Type"),
			"%s should return Content-Type: %s", path, expectedMimeType)
	}
}

func testBMPImageRequest(t *testing.T, baseURL, path, acceptHeader string,
	expectedStatus int) {
	
	req, err := http.NewRequest("GET", baseURL+path, nil)
	require.NoError(t, err)
	
	req.Header.Set("Accept", acceptHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"%s should return status code %d", path, expectedStatus)
	
	contentType := resp.Header.Get("Content-Type")
	assert.True(t, contentType == "image/bmp" || contentType == "image/x-ms-bmp",
		"%s should return Content-Type: image/bmp or image/x-ms-bmp, got %s", path, contentType)
}

func testInvalidRequest(t *testing.T, baseURL, path, acceptHeader string,
	expectedStatus int) {
	
	req, err := http.NewRequest("GET", baseURL+path, nil)
	require.NoError(t, err)
	
	req.Header.Set("Accept", acceptHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	assert.Equal(t, expectedStatus, resp.StatusCode,
		"%s should return status code %d", path, expectedStatus)
}

// getImageMimeType detects the MIME type of image data
func getImageMimeType(data []byte) string {
	if len(data) < 12 {
		return "unknown"
	}
	
	// JPEG
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "image/jpeg"
	}
	
	// PNG
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "image/png"
	}
	
	// GIF
	if bytes.HasPrefix(data, []byte("GIF87a")) || bytes.HasPrefix(data, []byte("GIF89a")) {
		return "image/gif"
	}
	
	// WebP
	if len(data) >= 12 && bytes.Equal(data[0:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP")) {
		return "image/webp"
	}
	
	// BMP
	if bytes.HasPrefix(data, []byte{0x42, 0x4D}) {
		return "image/bmp"
	}
	
	return "unknown"
}