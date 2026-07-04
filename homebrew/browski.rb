class Browski < Formula
  desc "Browser Profile Picker"
  homepage "https://github.com/icaliskanoglu/browski"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/icaliskanoglu/browski/releases/download/v#{version}/browski-macos-arm64.dmg"
      sha256 "REPLACE_WITH_ARM64_SHA256"
    else
      url "https://github.com/icaliskanoglu/browski/releases/download/v#{version}/browski-macos-amd64.dmg"
      sha256 "REPLACE_WITH_AMD64_SHA256"
    end
  end

  def install
    prefix.install "Browski.app"
    bin.write_exec_script "#{prefix}/Browski.app/Contents/MacOS/browski"
  end

  def caveats
    <<~EOS
      Browski has been installed. To run it:
        browski [url]

      Or launch the app from Applications folder.

      To set Browski as your default browser:
        1. Open System Preferences → General → Default web browser
        2. Select Browski from the list
    EOS
  end

  test do
    system "#{bin}/browski", "--version"
  end
end
