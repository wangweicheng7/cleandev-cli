class Cleaner < Formula
  desc "Safe-first macOS developer junk cleaner CLI"
  homepage "https://github.com/wangweicheng7/cleandev-cli"
  url "https://github.com/wangweicheng7/cleandev-cli/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "7a5f92b4bdcaa9c3162662d735692c315be99048ac6e29efeaf9599357427fe8"
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
