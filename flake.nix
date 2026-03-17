{
  description = "A command-line tool to control Home Assistant devices from the terminal";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      # Read version from VERSION file, fall back to git short rev or "dev"
      version = let
        versionFile = builtins.readFile ./VERSION;
        trimmed = builtins.replaceStrings [ "\n" "\r" " " ] [ "" "" "" ] versionFile;
      in
        if trimmed != "" then trimmed
        else if (self ? shortRev) then self.shortRev
        else "dev";

      commit = if (self ? rev) then self.rev else "dirty";

      # Format lastModifiedDate: YYYYMMDDHHMMSS -> YYYY-MM-DDTHH:MM:SSZ
      date = let
        raw = self.lastModifiedDate or "19700101000000";
        year = builtins.substring 0 4 raw;
        month = builtins.substring 4 2 raw;
        day = builtins.substring 6 2 raw;
        hour = builtins.substring 8 2 raw;
        min = builtins.substring 10 2 raw;
        sec = builtins.substring 12 2 raw;
      in "${year}-${month}-${day}T${hour}:${min}:${sec}Z";

      # Package builder function - used by both overlay and packages output
      mkHctl = pkgs: pkgs.buildGoModule {
        pname = "hctl";
        inherit version;

        src = ./.;

        vendorHash = "sha256-4FrLHwbZqvRE547PfeZowwNH1+zI4Ut04wtyvUfxGug=";

        env.CGO_ENABLED = 0;

        ldflags = [
          "-s"
          "-w"
          "-X github.com/xx4h/hctl/cmd.version=${version}"
          "-X github.com/xx4h/hctl/cmd.commit=${commit}"
          "-X github.com/xx4h/hctl/cmd.date=${date}"
        ];

        # Skip tests that require network access (sandbox blocks it)
        checkFlags = [ "-skip" "Play" ];

        meta = with pkgs.lib; {
          description = "A command-line tool to control Home Assistant devices";
          homepage = "https://github.com/xx4h/hctl";
          license = licenses.asl20;
          maintainers = [ ];
          mainProgram = "hctl";
        };
      };
    in
    {
      # Overlay for users who want to use their own nixpkgs
      overlays.default = final: prev: {
        hctl = mkHctl final;
      };
    }
    //
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          hctl = mkHctl pkgs;
          default = self.packages.${system}.hctl;
        };

        apps = {
          hctl = flake-utils.lib.mkApp {
            drv = self.packages.${system}.hctl;
          };
          default = self.apps.${system}.hctl;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_25
            golangci-lint
            goreleaser
            gnumake
          ];

          shellHook = ''
            echo "hctl development shell"
            echo "Go: $(go version)"
          '';
        };
      }
    );
}
