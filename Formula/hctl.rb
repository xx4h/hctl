class Hctl < Formula
  desc "Tool to control your Home Assistant devices from the command-line"
  homepage "https://github.com/xx4h/hctl"
  url "https://github.com/xx4h/hctl.git",
      tag:      "v0.5.0",
      revision: "9c3a308600bce627bd79e0d3d1a2f05b7a13c347"
  license "Apache-2.0"
  head "https://github.com/xx4h/hctl.git", branch: "main"

  livecheck do
    url :stable
    regex(/^v?(\d+(?:\.\d+)+)$/i)
  end

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/xx4h/hctl/cmd.version=v#{version}
      -X github.com/xx4h/hctl/cmd.commit=#{Utils.git_head}
      -X github.com/xx4h/hctl/cmd.date=#{time.iso8601}
    ]
    system "go", "build", *std_go_args(ldflags:)

    generate_completions_from_executable(bin/"hctl", "completion")
  end

  test do
    assert_match "Hctl is a CLI tool to control your home automation", shell_output("#{bin}/hctl --help")
  end
end
