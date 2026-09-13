. (Join-Path $PSScriptRoot 'demo-common.ps1')

Test-DemoDocker
Push-Location $script:DemoProjectRoot
try {
    Invoke-DemoCompose ps
    $health = Invoke-WebRequest -UseBasicParsing -TimeoutSec 10 "http://127.0.0.1:$script:DemoPort/api/ping"
    if ($health.StatusCode -ne 200) { throw "Gateway health check returned HTTP $($health.StatusCode)." }
    $lanIP = Get-DemoLanIPv4
    Write-Host "Local gateway is healthy: http://localhost:$script:DemoPort" -ForegroundColor Green
    if ($lanIP) { Write-Host "Test from another device: http://${lanIP}:$script:DemoPort" -ForegroundColor Cyan }
} finally {
    Pop-Location
}
