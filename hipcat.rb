class HipCat < Formula
  desc "Simple command-line Utility to post snippets to HipChat."
  homepage "https://github.com/samdfonseca/hipcat"
  url "https://github.com/samdfonseca/hipcat/archive/v0.1.tar.gz
  version "0.1"
  sha256 "58beac16e8949a793400025ea3ce159220f21cbf3f92bf8e5530d7662d3132e9"

  depends_on "go"

  def install
    platform = `uname`.downcase.strip

    unless ENV["GOPATH"]
      ENV["GOPATH"] = "/tmp"
    end

    system "make"
    bin.install "build/hipcat-0.6-#{platform}-amd64" => "hipcat"
  end

  test do
    assert_equal(0, "/usr/local/bin/hipcat")
  end
end
