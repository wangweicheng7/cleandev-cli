class Cleaner < Formula
  desc "Safe-first macOS developer junk cleaner CLI"
  homepage "https://github.com/your-org/cleandev-cli"
  url "https://github.com/your-org/cleandev-cli/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_RELEASE_TARBALL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/cleaner"
  end

  test do
    output = shell_output("#{bin}/cleaner doctor")
    assert_match "Cleaner Doctor", output
  end
end
