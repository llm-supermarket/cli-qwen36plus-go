param(
    [Parameter(Mandatory = $true)]
    [string]$Version
)

$repo = "llm-supermarket/cli-qwen36plus-go"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) "rclone-encrypt-release-$Version"

if (-not (Test-Path $tempDir)) {
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
}

Write-Host "Downloading Windows release asset for v$Version ..."
gh release download "v$Version" --repo $repo --pattern "rclone-encrypt-windows-amd64.exe" --dir $tempDir

$exePath = Join-Path $tempDir "rclone-encrypt-windows-amd64.exe"

if (-not (Test-Path $exePath)) {
    throw "Unable to locate rclone-encrypt-windows-amd64.exe at $exePath"
}

$hash = (Get-FileHash -Path $exePath -Algorithm SHA256).Hash.ToLower()

Write-Host "Hash: $hash"

$url = "https://github.com/$repo/releases/download/v$Version/rclone-encrypt-windows-amd64.exe"

$manifestPath = Join-Path $PSScriptRoot "rclone-encrypt.json"

$manifest = Get-Content -Path $manifestPath -Raw | ConvertFrom-Json

$manifest.version = $Version

$manifest.architecture."64bit".url = $url

$manifest.architecture."64bit".hash = $hash

$manifest | ConvertTo-Json -Depth 10 | Set-Content -Path $manifestPath -NoNewline

# Clean up temp files.
Remove-Item $tempDir -Recurse -Force

Write-Host "Updated rclone-encrypt.json to v$Version"
