Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:DemoProjectRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$script:DemoComposeFile = Join-Path $script:DemoProjectRoot 'docker-compose.demo.yml'
$script:DemoPort = if ($env:DEMO_HTTP_PORT) { [int]$env:DEMO_HTTP_PORT } else { 8080 }

function Invoke-DemoCompose {
    param([Parameter(ValueFromRemainingArguments = $true)][string[]]$ComposeArguments)
    & docker-compose -f $script:DemoComposeFile @ComposeArguments
    if ($LASTEXITCODE -ne 0) {
        throw "docker-compose failed with exit code $LASTEXITCODE."
    }
}

function Get-DemoLanIPv4 {
    $routeOutput = (& route print -4) -join "`n"
    $defaultRoute = [regex]::Match($routeOutput, '(?m)^\s*0\.0\.0\.0\s+0\.0\.0\.0\s+\S+\s+(?<ip>\d{1,3}(?:\.\d{1,3}){3})\s+\d+\s*$')
    if ($defaultRoute.Success) { return $defaultRoute.Groups['ip'].Value }
    return $null
}

function Test-DemoDocker {
    & docker info *> $null
    if ($LASTEXITCODE -ne 0) {
        throw 'Docker Desktop is not ready. Start Docker Desktop and try again.'
    }
}
