. (Join-Path $PSScriptRoot 'demo-common.ps1')

Test-DemoDocker
Push-Location $script:DemoProjectRoot
try {
    Invoke-DemoCompose down
    Write-Host 'Agora-BBS demo containers stopped; the demo database volume was preserved.' -ForegroundColor Green
} finally {
    Pop-Location
}
