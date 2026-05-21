param(
    [Parameter(Position = 0)]
    [ValidateSet("start", "stop", "status", "remove")]
    [string]$Action = "start"
)

$ErrorActionPreference = "Stop"

$ScriptDir = $PSScriptRoot
$PodYaml = Join-Path $ScriptDir "pod.yaml"
$PodName = "neonya"

function Assert-PodmanAvailable {
    if (-not (Get-Command podman -ErrorAction SilentlyContinue)) {
        Write-Error "podman not found. Please install podman first."
        exit 1
    }
}

function Start-Pod {
    Write-Host "Starting $PodName pod..." -ForegroundColor Cyan
    Push-Location $ScriptDir
    try {
        podman play kube $PodYaml --replace
        Write-Host "Pod $PodName started successfully." -ForegroundColor Green
    } finally {
        Pop-Location
    }
}

function Stop-Pod {
    Write-Host "Stopping $PodName pod..." -ForegroundColor Cyan
    $running = podman pod ps --filter "name=$PodName" --format "{{.Name}}" 2>$null
    if ($running) {
        podman pod stop $PodName
        podman pod rm $PodName
        Write-Host "Pod $PodName stopped and removed." -ForegroundColor Green
    } else {
        Write-Host "Pod $PodName is not running." -ForegroundColor Yellow
    }
}

function Show-Status {
    Write-Host "Pod $PodName status:" -ForegroundColor Cyan
    $podInfo = podman pod ps --filter "name=$PodName" --format "table {{.Name}}`t{{.Status}}`t{{.Created}}`t{{.InfraId}}" 2>$null
    if ($podInfo) {
        Write-Host $podInfo
        podman ps --filter "pod=$PodName" --format "table {{.Names}}`t{{.Image}}`t{{.Status}}`t{{.Ports}}"
    } else {
        Write-Host "Pod $PodName not found." -ForegroundColor Yellow
    }
}

function Remove-Pod {
    Write-Host "Removing $PodName pod and all its resources..." -ForegroundColor Cyan
    Push-Location $ScriptDir
    try {
        podman play kube $PodYaml --down 2>$null
        Write-Host "Pod $PodName removed." -ForegroundColor Green
    } finally {
        Pop-Location
    }
}

Assert-PodmanAvailable

switch ($Action) {
    "start"  { Start-Pod }
    "stop"   { Stop-Pod }
    "status" { Show-Status }
    "remove" { Remove-Pod }
}
