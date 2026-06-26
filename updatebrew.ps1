param(
    [Parameter(Mandatory = $true)]
    [string]$Version
)

$repo = "llm-supermarket-org/cli-qwen36plus-go"
$platforms = @("darwin-amd64", "darwin-arm64", "linux-amd64", "linux-arm64")
$formulaPath = "$PSScriptRoot/Formula/rclone-encrypt.rb"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) "rclone-encrypt-release-$Version"

if (-not (Test-Path $tempDir)) {
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
}

Write-Host "Downloading release assets for v$Version ..."
gh release download "v$Version" --repo $repo --pattern "*.tar.gz" --dir $tempDir

# Compute SHA256 for each platform.
$hash = @{}
foreach ($platform in $platforms) {
    $file = Join-Path $tempDir "rclone-encrypt-$platform.tar.gz"
    if (-not (Test-Path $file)) {
        throw "Missing asset: $file"
    }
    $hash[$platform] = (Get-FileHash -Path $file -Algorithm SHA256).Hash.ToLower()
    Write-Host "SHA256 for ${platform}: $($hash[$platform])"
}

# Clean up temp files.
Remove-Item $tempDir -Recurse -Force

# Regenerate the formula wholesale so it stays correct across releases.
$formula = @"
class RcloneEncrypt < Formula
  desc "CLI tool to encrypt and decrypt files using rclone-compatible cryptography"
  homepage "https://github.com/$repo"
  version "$Version"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/$repo/releases/download/v$Version/rclone-encrypt-darwin-arm64.tar.gz"
      sha256 "$($hash['darwin-arm64'])"
    else
      url "https://github.com/$repo/releases/download/v$Version/rclone-encrypt-darwin-amd64.tar.gz"
      sha256 "$($hash['darwin-amd64'])"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/$repo/releases/download/v$Version/rclone-encrypt-linux-arm64.tar.gz"
      sha256 "$($hash['linux-arm64'])"
    else
      url "https://github.com/$repo/releases/download/v$Version/rclone-encrypt-linux-amd64.tar.gz"
      sha256 "$($hash['linux-amd64'])"
    end
  end

  def install
    bin.install "rclone-encrypt-darwin-arm64" => "rclone-encrypt" if OS.mac? && Hardware::CPU.arm?
    bin.install "rclone-encrypt-darwin-amd64" => "rclone-encrypt" if OS.mac? && !Hardware::CPU.arm?
    bin.install "rclone-encrypt-linux-arm64" => "rclone-encrypt" if OS.linux? && Hardware::CPU.arm?
    bin.install "rclone-encrypt-linux-amd64" => "rclone-encrypt" if OS.linux? && !Hardware::CPU.arm?
  end

  test do
    assert_match "rclone-encrypt #{version}", shell_output("#{bin}/rclone-encrypt --version")
  end
end
"@

Set-Content -Path $formulaPath -Value $formula -NoNewline
Write-Host "Wrote $formulaPath for version $Version"
