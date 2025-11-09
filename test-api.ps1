# Comprehensive API Test Suite for Windows
# This tests all endpoints and provides detailed feedback

$baseUrl = "http://localhost:8080"
$testResults = @()
$passed = 0
$failed = 0

function Write-TestHeader {
    param([string]$title)
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "TEST: $title" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
}

function Write-TestResult {
    param(
        [string]$testName,
        [bool]$success,
        [string]$message,
        [object]$data = $null
    )

    $script:testResults += [PSCustomObject]@{
        Test = $testName
        Status = if ($success) { "PASS" } else { "FAIL" }
        Message = $message
        Data = $data
    }

    if ($success) {
        $script:passed++
        Write-Host "✓ PASS: $testName" -ForegroundColor Green
        if ($message) {
            Write-Host "  $message" -ForegroundColor Gray
        }
    } else {
        $script:failed++
        Write-Host "✗ FAIL: $testName" -ForegroundColor Red
        if ($message) {
            Write-Host "  $message" -ForegroundColor Yellow
        }
    }

    if ($data) {
        Write-Host "  Response: $($data | ConvertTo-Json -Compress)" -ForegroundColor DarkGray
    }
}

# Test 1: Health Check
Write-TestHeader "Health Check Endpoint"
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get -TimeoutSec 5
    if ($response.status -eq "healthy") {
        Write-TestResult "Health Check" $true "Service is healthy" $response
    } else {
        Write-TestResult "Health Check" $false "Unexpected status: $($response.status)" $response
    }
} catch {
    Write-TestResult "Health Check" $false "Failed to connect: $($_.Exception.Message)"
}

# Test 2: Shorten URL (Random Code)
Write-TestHeader "Shorten URL with Random Code"
$testUrl = "https://example.com/test-$(Get-Random)"
try {
    $body = @{
        url = $testUrl
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "$baseUrl/api/shorten" -Method Post -Body $body -ContentType "application/json" -TimeoutSec 10

    if ($response.short_url -and $response.short_code) {
        Write-TestResult "Shorten URL (Random)" $true "Created short URL: $($response.short_url)" $response
        $global:randomShortCode = $response.short_code
    } else {
        Write-TestResult "Shorten URL (Random)" $false "Missing short_url or short_code" $response
    }
} catch {
    Write-TestResult "Shorten URL (Random)" $false "Request failed: $($_.Exception.Message)"
}

# Test 3: Shorten URL (Custom Alias)
Write-TestHeader "Shorten URL with Custom Alias"
$customAlias = "test-$(Get-Random -Maximum 9999)"
try {
    $body = @{
        url = "https://example.com/custom-test"
        custom_alias = $customAlias
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "$baseUrl/api/shorten" -Method Post -Body $body -ContentType "application/json" -TimeoutSec 10

    if ($response.short_code -eq $customAlias) {
        Write-TestResult "Shorten URL (Custom)" $true "Created with custom alias: $customAlias" $response
        $global:customShortCode = $customAlias
    } else {
        Write-TestResult "Shorten URL (Custom)" $false "Alias mismatch. Expected: $customAlias, Got: $($response.short_code)" $response
    }
} catch {
    Write-TestResult "Shorten URL (Custom)" $false "Request failed: $($_.Exception.Message)"
}

# Test 4: Invalid URL
Write-TestHeader "Validation - Invalid URL"
try {
    $body = @{
        url = "not-a-valid-url"
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "$baseUrl/api/shorten" -Method Post -Body $body -ContentType "application/json" -TimeoutSec 10 -ErrorAction Stop
    Write-TestResult "Invalid URL Validation" $false "Should have rejected invalid URL" $response
} catch {
    if ($_.Exception.Response.StatusCode -eq 400) {
        Write-TestResult "Invalid URL Validation" $true "Correctly rejected invalid URL"
    } else {
        Write-TestResult "Invalid URL Validation" $false "Unexpected error: $($_.Exception.Message)"
    }
}

# Test 5: Duplicate Custom Alias
Write-TestHeader "Validation - Duplicate Alias"
if ($global:customShortCode) {
    try {
        $body = @{
            url = "https://example.com/another-url"
            custom_alias = $global:customShortCode
        } | ConvertTo-Json

        $response = Invoke-RestMethod -Uri "$baseUrl/api/shorten" -Method Post -Body $body -ContentType "application/json" -TimeoutSec 10 -ErrorAction Stop
        Write-TestResult "Duplicate Alias" $false "Should have rejected duplicate alias" $response
    } catch {
        if ($_.Exception.Response.StatusCode -eq 400) {
            Write-TestResult "Duplicate Alias" $true "Correctly rejected duplicate alias"
        } else {
            Write-TestResult "Duplicate Alias" $false "Unexpected error: $($_.Exception.Message)"
        }
    }
} else {
    Write-TestResult "Duplicate Alias" $false "Skipped: No custom code from previous test"
}

# Test 6: URL Redirect
Write-TestHeader "URL Redirect"
if ($global:randomShortCode) {
    try {
        $response = Invoke-WebRequest -Uri "$baseUrl/$($global:randomShortCode)" -MaximumRedirection 0 -ErrorAction SilentlyContinue

        if ($response.StatusCode -eq 302 -or $response.StatusCode -eq 301) {
            $location = $response.Headers.Location
            Write-TestResult "URL Redirect" $true "Redirects to: $location" @{redirect=$location}
        } else {
            Write-TestResult "URL Redirect" $false "Expected 302/301, got $($response.StatusCode)"
        }
    } catch {
        # Check if it's a redirect (which PowerShell treats as an error by default)
        if ($_.Exception.Response.StatusCode -eq 302 -or $_.Exception.Response.StatusCode -eq 301) {
            $location = $_.Exception.Response.Headers.Location
            Write-TestResult "URL Redirect" $true "Redirects to: $location" @{redirect=$location}
        } else {
            Write-TestResult "URL Redirect" $false "Request failed: $($_.Exception.Message)"
        }
    }
} else {
    Write-TestResult "URL Redirect" $false "Skipped: No short code from previous test"
}

# Test 7: Non-existent Short Code
Write-TestHeader "404 - Non-existent Code"
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/nonexistent123456" -Method Get -TimeoutSec 5 -ErrorAction Stop
    Write-TestResult "404 Not Found" $false "Should have returned 404" $response
} catch {
    if ($_.Exception.Response.StatusCode -eq 404) {
        Write-TestResult "404 Not Found" $true "Correctly returned 404"
    } else {
        Write-TestResult "404 Not Found" $false "Unexpected status: $($_.Exception.Response.StatusCode)"
    }
}

# Test 8: Get Link Info
Write-TestHeader "Get Link Information"
if ($global:customShortCode) {
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/links/$($global:customShortCode)" -Method Get -TimeoutSec 5

        if ($response.short_code -eq $global:customShortCode) {
            Write-TestResult "Get Link Info" $true "Retrieved link information" $response
        } else {
            Write-TestResult "Get Link Info" $false "Unexpected response" $response
        }
    } catch {
        Write-TestResult "Get Link Info" $false "Request failed: $($_.Exception.Message)"
    }
} else {
    Write-TestResult "Get Link Info" $false "Skipped: No short code available"
}

# Test 9: Analytics
Write-TestHeader "Get Analytics"
if ($global:randomShortCode) {
    try {
        $response = Invoke-RestMethod -Uri "$baseUrl/api/analytics/$($global:randomShortCode)" -Method Get -TimeoutSec 5

        if ($response.short_code -eq $global:randomShortCode) {
            Write-TestResult "Get Analytics" $true "Analytics data: $($response.total_clicks) clicks" $response
        } else {
            Write-TestResult "Get Analytics" $false "Unexpected response" $response
        }
    } catch {
        Write-TestResult "Get Analytics" $false "Request failed: $($_.Exception.Message)"
    }
} else {
    Write-TestResult "Get Analytics" $false "Skipped: No short code available"
}

# Test 10: QR Code Generation
Write-TestHeader "QR Code Generation"
if ($global:customShortCode) {
    try {
        $response = Invoke-WebRequest -Uri "$baseUrl/api/qr/$($global:customShortCode)?size=256" -Method Get -TimeoutSec 5

        if ($response.StatusCode -eq 200 -and $response.Headers.'Content-Type' -like '*image*') {
            $size = $response.Content.Length
            Write-TestResult "QR Code Generation" $true "Generated QR code (${size} bytes)" @{size=$size;contentType=$response.Headers.'Content-Type'}
        } else {
            Write-TestResult "QR Code Generation" $false "Unexpected response type"
        }
    } catch {
        Write-TestResult "QR Code Generation" $false "Request failed: $($_.Exception.Message)"
    }
} else {
    Write-TestResult "QR Code Generation" $false "Skipped: No short code available"
}

# Test 11: List User Links
Write-TestHeader "List All User Links"
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/api/links" -Method Get -TimeoutSec 5

    if ($response.urls) {
        $count = $response.urls.Count
        Write-TestResult "List User Links" $true "Retrieved $count links" @{count=$count}
    } else {
        Write-TestResult "List User Links" $true "No links found (empty list is valid)" $response
    }
} catch {
    Write-TestResult "List User Links" $false "Request failed: $($_.Exception.Message)"
}

# Test 12: Rate Limiting
Write-TestHeader "Rate Limiting"
try {
    $requests = 0
    $rateLimited = $false

    # Make rapid requests
    for ($i = 0; $i -lt 15; $i++) {
        try {
            Invoke-RestMethod -Uri "$baseUrl/health" -Method Get -TimeoutSec 2 | Out-Null
            $requests++
        } catch {
            if ($_.Exception.Response.StatusCode -eq 429) {
                $rateLimited = $true
                break
            }
        }
        Start-Sleep -Milliseconds 50
    }

    if ($rateLimited) {
        Write-TestResult "Rate Limiting" $true "Rate limit triggered after $requests requests"
    } else {
        Write-TestResult "Rate Limiting" $false "Rate limit not triggered (made $requests requests)" @{note="Rate limit may be configured higher"}
    }
} catch {
    Write-TestResult "Rate Limiting" $false "Test failed: $($_.Exception.Message)"
}

# Test 13: CORS Headers
Write-TestHeader "CORS Headers"
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/health" -Method Get -TimeoutSec 5

    $corsHeaders = @{}
    if ($response.Headers.'Access-Control-Allow-Origin') {
        $corsHeaders.AllowOrigin = $response.Headers.'Access-Control-Allow-Origin'
    }
    if ($response.Headers.'Access-Control-Allow-Methods') {
        $corsHeaders.AllowMethods = $response.Headers.'Access-Control-Allow-Methods'
    }

    if ($corsHeaders.Count -gt 0) {
        Write-TestResult "CORS Headers" $true "CORS headers present" $corsHeaders
    } else {
        Write-TestResult "CORS Headers" $false "No CORS headers found"
    }
} catch {
    Write-TestResult "CORS Headers" $false "Request failed: $($_.Exception.Message)"
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "TEST SUMMARY" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Total Tests: $($passed + $failed)" -ForegroundColor White
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor Red

if ($failed -eq 0) {
    Write-Host "`n✓ ALL TESTS PASSED!" -ForegroundColor Green
    Write-Host "Your URL Shortener is working correctly!" -ForegroundColor Green
} else {
    Write-Host "`n✗ SOME TESTS FAILED" -ForegroundColor Red
    Write-Host "Please check the errors above and ensure:" -ForegroundColor Yellow
    Write-Host "  1. Services are running (docker-compose ps)" -ForegroundColor Yellow
    Write-Host "  2. DynamoDB tables are created" -ForegroundColor Yellow
    Write-Host "  3. Redis is connected" -ForegroundColor Yellow
    Write-Host "  4. AWS credentials are valid in .env" -ForegroundColor Yellow
}

Write-Host "`n========================================`n" -ForegroundColor Cyan

# Export results to file
$testResults | Export-Csv -Path "test-results.csv" -NoTypeInformation
Write-Host "Detailed results saved to: test-results.csv" -ForegroundColor Gray

# Return exit code
exit $failed
