# Build script for Windows PowerShell
param(
    [string]$Mode = "release"
)

Write-Host "Building Decomposed Key-Value Store Services..." -ForegroundColor Green

function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Error "Go is not installed. Please install Go 1.21 or later."
    exit 1
}

# Create bin directory
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Build flags
$buildFlags = @("-ldflags", "-w -s")
if ($Mode -eq "debug") {
    $buildFlags = @()
    Write-Status "Building in debug mode..."
} else {
    Write-Status "Building in release mode..."
}

# Build KV service
Write-Status "Building KV service..."
$args = @("build") + $buildFlags + @("-o", "bin\kv-service.exe", ".\cmd\kv-service")
& go @args
if ($LASTEXITCODE -eq 0) {
    Write-Status "✓ KV service built successfully"
} else {
    Write-Error "✗ Failed to build KV service"
    exit 1
}

# Build API service
Write-Status "Building API service..."
$args = @("build") + $buildFlags + @("-o", "bin\api-service.exe", ".\cmd\api-service")
& go @args
if ($LASTEXITCODE -eq 0) {
    Write-Status "✓ API service built successfully"
} else {
    Write-Error "✗ Failed to build API service"
    exit 1
}

Write-Status "All services built successfully!"
Write-Status "Binaries are available in the .\bin\ directory:"
Get-ChildItem -Path "bin" | Format-Table Name, Length, LastWriteTime

Write-Host ""
Write-Status "To run the services:"
Write-Host "  KV Service:  .\bin\kv-service.exe"
Write-Host "  API Service: .\bin\api-service.exe"