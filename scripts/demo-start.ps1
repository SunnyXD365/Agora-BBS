. (Join-Path $PSScriptRoot 'demo-common.ps1')

Test-DemoDocker
Push-Location $script:DemoProjectRoot
try {
    Write-Host 'Building and starting the Agora-BBS LAN demo...' -ForegroundColor Cyan
    Invoke-DemoCompose up -d --build --wait

    Write-Host 'Loading the repeatable demo accounts and topics...' -ForegroundColor Cyan
    Invoke-DemoCompose run --rm agora-demo-seed

    $health = Invoke-WebRequest -UseBasicParsing -TimeoutSec 10 "http://127.0.0.1:$script:DemoPort/api/ping"
    if ($health.StatusCode -ne 200) { throw "Gateway health check returned HTTP $($health.StatusCode)." }

    $lanIP = Get-DemoLanIPv4
    Write-Host ''
    Write-Host 'Agora-BBS demo is ready.' -ForegroundColor Green
    Write-Host "Local URL: http://localhost:$script:DemoPort"
    if ($lanIP) {
        Write-Host "LAN URL: http://${lanIP}:$script:DemoPort" -ForegroundColor Green
    } else {
        Write-Warning 'Could not detect the LAN IPv4 address. Run ipconfig and use the current WLAN IPv4 address.'
    }
    Write-Host 'Test the LAN URL from another computer or phone before the presentation.'
    Write-Host 'If it is unreachable, check campus-network client isolation and the documented firewall rule.'
} finally {
    Pop-Location
}
