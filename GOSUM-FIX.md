# Quick Fix Guide for go.sum Error

## Problem
You're seeing this error when building Docker:
```
"/go.sum": not found
```

## Why This Happens
We added AWS Lambda dependencies to `go.mod`, but the `go.sum` file (which tracks exact versions) wasn't created because you don't have Go installed locally.

## ✅ Solution (Choose One)

### Option 1: Automatic Fix (Easiest!)

Run this script to generate `go.sum` using Docker:

**PowerShell:**
```powershell
.\fix-gosum.ps1
```

**Command Prompt:**
```cmd
.\fix-gosum.bat
```

This uses Docker to run `go mod tidy` and create the `go.sum` file.

### Option 2: Manual Docker Command

```powershell
docker run --rm -v ${PWD}:/app -w /app golang:1.23-alpine go mod tidy
```

### Option 3: Updated run.bat (Already Done!)

Just run:
```cmd
.\run.bat
```

It will automatically detect missing `go.sum` and generate it!

## After Fixing

Once `go.sum` is created, you can:

```powershell
# Build and start services
docker-compose up -d

# Or rebuild if needed
docker-compose build
docker-compose up -d
```

## Why We Need go.sum

- `go.mod` lists dependencies (like package.json)
- `go.sum` locks exact versions (like package-lock.json)
- Docker needs both to build reproducibly

## Prevention

The updated `Dockerfile` now handles missing `go.sum` automatically, but it's better to have it committed to the repo.
