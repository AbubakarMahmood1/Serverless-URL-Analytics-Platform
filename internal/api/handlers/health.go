package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "url-shortener",
	})
}

// Welcome handles GET / - Shows API documentation
func (h *HealthHandler) Welcome(c *fiber.Ctx) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>URL Shortener API</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .endpoint { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { display: inline-block; padding: 3px 8px; border-radius: 3px; font-weight: bold; font-size: 12px; }
        .post { background: #49cc90; color: white; }
        .get { background: #61affe; color: white; }
        .delete { background: #f93e3e; color: white; }
        code { background: #e8e8e8; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🔗 URL Shortener API</h1>
    <p>A serverless URL analytics platform with comprehensive tracking.</p>

    <h2>Available Endpoints</h2>

    <div class="endpoint">
        <span class="method get">GET</span>
        <strong>/health</strong>
        <p>Health check endpoint - Returns service status</p>
    </div>

    <div class="endpoint">
        <span class="method post">POST</span>
        <strong>/api/shorten</strong>
        <p>Create a shortened URL</p>
        <pre><code>{
  "url": "https://example.com/very/long/url",
  "custom_alias": "my-link" (optional)
}</code></pre>
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <strong>/:shortCode</strong>
        <p>Redirect to original URL (e.g., <code>/abc123</code>)</p>
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <strong>/api/analytics/:shortCode</strong>
        <p>Get analytics for a shortened URL</p>
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <strong>/api/qr/:shortCode</strong>
        <p>Generate QR code for a shortened URL<br>
        Query params: <code>?size=256&format=png</code></p>
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <strong>/api/links</strong>
        <p>List all your shortened URLs</p>
    </div>

    <div class="endpoint">
        <span class="method get">GET</span>
        <strong>/api/links/:shortCode</strong>
        <p>Get details about a specific shortened URL</p>
    </div>

    <div class="endpoint">
        <span class="method delete">DELETE</span>
        <strong>/api/links/:shortCode</strong>
        <p>Delete a shortened URL</p>
    </div>

    <hr>
    <p><small>Status: <strong style="color: green;">Running</strong> | Version: 1.0.0</small></p>
</body>
</html>
`
	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}
