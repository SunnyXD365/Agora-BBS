param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('Enable', 'Disable')]
    [string]$Action
)

$ErrorActionPreference = 'Stop'
$ruleName = 'Agora-BBS Classroom Demo'
$demoPort = if ($env:DEMO_HTTP_PORT) { [int]$env:DEMO_HTTP_PORT } else { 8080 }

$identity = [Security.Principal.WindowsIdentity]::GetCurrent()
$principal = [Security.Principal.WindowsPrincipal]::new($identity)
if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    throw 'Run this script in a PowerShell window opened as Administrator.'
}

if ($Action -eq 'Enable') {
    Remove-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue
    New-NetFirewallRule -DisplayName $ruleName -Direction Inbound -Action Allow -Protocol TCP -LocalPort $demoPort -RemoteAddress LocalSubnet -Profile Any | Out-Null
    Write-Host "Allowed LocalSubnet access to TCP $demoPort. After the demo, run: .\scripts\demo-firewall.ps1 -Action Disable" -ForegroundColor Green
} else {
    Remove-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue
    Write-Host 'Removed the Agora-BBS LAN demo firewall rule.' -ForegroundColor Green
}
