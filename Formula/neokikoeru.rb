class Neokikoeru < Formula
  desc "Cloud storage-based web media player for DLsite voice works"
  homepage "https://github.com/vscodev/neokikoeru"
  version "0.38.0"

  if OS.mac?
    if Hardware::CPU.arm? || Hardware::CPU.in_rosetta2?
      url "https://github.com/vscodev/neokikoeru/releases/download/v0.38.0/neokikoeru-macos-arm64.tar.gz"
      sha256 "9d7e2e48ee2683bfe5d2ad6b1e191570bf73a16c18ac0b8f901fba071501fa6c"
    else
      url "https://github.com/vscodev/neokikoeru/releases/download/v0.38.0/neokikoeru-macos-amd64.tar.gz"
      sha256 "d2bf75c30bfc59be37b5326e96f198180ceb94d25ac453097e1660c0a5b383ab"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/vscodev/neokikoeru/releases/download/v0.38.0/neokikoeru-linux-arm64.tar.gz"
      sha256 "d18f09b9e42e005430c478ca163b4d25c76a1512e9c2c7b7d3ea18e3d482dd92"
    else
      url "https://github.com/vscodev/neokikoeru/releases/download/v0.38.0/neokikoeru-linux-amd64.tar.gz"
      sha256 "a615576533552f575307701cc79fa14401acf85b64acc8e768245103b183459f"
    end
  else
    odie "Unsupported platform. Please submit a bug report here: https://github.com/vscodev/neokikoeru/issues\n#{OS.report}"
  end

  def install
    bin.install "neokikoeru"
    generate_completions_from_executable(bin/"neokikoeru", "completion")
  end

  service do
    run [bin/"neokikoeru", "start"]
    keep_alive crashed: true
    environment_variables PATH: std_service_path_env
    log_path var/"log/neokikoeru.log"
    error_log_path var/"log/neokikoeru.log"
  end

  test do
    assert_match "neokikoeru version #{version}", shell_output("#{bin}/neokikoeru -v")
  end
end
