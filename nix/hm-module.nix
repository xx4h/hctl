# Home Manager module for hctl.
#
# Used via the flake's `homeManagerModules.default` (or `.hctl`) output.
# `self` is captured at flake evaluation so the default package resolves
# from this flake's outputs, regardless of whether the user has applied
# our overlay to their nixpkgs.
self:
{ config, lib, pkgs, ... }:

let
  cfg = config.programs.hctl;
  defaultPackage = self.packages.${pkgs.stdenv.hostPlatform.system}.hctl;
in
{
  options.programs.hctl = {
    enable = lib.mkEnableOption
      "hctl, a command-line tool to control Home Assistant devices";

    package = lib.mkOption {
      type = lib.types.package;
      default = defaultPackage;
      defaultText = lib.literalExpression
        "hctl.packages.\${pkgs.stdenv.hostPlatform.system}.hctl";
      description = "The hctl package to install.";
    };

    enableBashIntegration = lib.mkEnableOption "Bash shell completion" // {
      default = true;
    };

    enableZshIntegration = lib.mkEnableOption "Zsh shell completion" // {
      default = true;
    };

    enableFishIntegration = lib.mkEnableOption "Fish shell completion" // {
      default = true;
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];

    # The package already drops completion files under
    # $out/share/{bash-completion,zsh/site-functions,fish/vendor_completions.d},
    # which HM-managed shells pick up automatically when their respective
    # `enableCompletion` option is on. The hooks below additionally source
    # completions directly from the binary, matching the idiom used by
    # `programs.atuin` and similar modules — useful when a shell isn't
    # configured to auto-load vendor completions.
    programs.bash.initExtra = lib.mkIf cfg.enableBashIntegration ''
      source <(${lib.getExe cfg.package} completion bash)
    '';

    programs.zsh.initContent = lib.mkIf cfg.enableZshIntegration ''
      source <(${lib.getExe cfg.package} completion zsh)
    '';

    programs.fish.interactiveShellInit = lib.mkIf cfg.enableFishIntegration ''
      ${lib.getExe cfg.package} completion fish | source
    '';
  };
}
