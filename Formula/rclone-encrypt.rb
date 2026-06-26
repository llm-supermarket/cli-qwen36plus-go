class RcloneEncrypt < Formula
  desc "CLI tool to encrypt and decrypt files using rclone-compatible cryptography"
  homepage "https://github.com/llm-supermarket/cli-qwen36plus-go"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/llm-supermarket/cli-qwen36plus-go/releases/download/v1.0.0/rclone-encrypt-darwin-arm64.tar.gz"
      sha256 "ea34798932928a7988af1fa924b76813961a3e3ea7ba92cb56f77ee45fbc0c3e"
    else
      url "https://github.com/llm-supermarket/cli-qwen36plus-go/releases/download/v1.0.0/rclone-encrypt-darwin-amd64.tar.gz"
      sha256 "cc99fd4a8c93175c4045551588ae63f8add621b958674fc98b001aeda10e3cac"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/llm-supermarket/cli-qwen36plus-go/releases/download/v1.0.0/rclone-encrypt-linux-arm64.tar.gz"
      sha256 "7a4dea32317a60f98ac6b955ac76ab7b9a8790cd0171c0de47a892bfcedef8d3"
    else
      url "https://github.com/llm-supermarket/cli-qwen36plus-go/releases/download/v1.0.0/rclone-encrypt-linux-amd64.tar.gz"
      sha256 "ab2a5630342c82f054b8adfdfbabd17ea7ccab8f4471b2b07480bbde38d4111a"
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