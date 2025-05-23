class Hctl < Formula
  desc "Tool to control your Home Assistant devices from the command-line"
  homepage "https://github.com/xx4h/hctl"
  url "https://github.com/xx4h/hctl.git",
      tag:      "v0.7.0",
      revision: "55c5cef08e80c34551e0206ab4eeba8af2ecae58"
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
