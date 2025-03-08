{
  description = "TUI Launcher for NixOS";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
  in {
    packages.x86_64-linux.default = pkgs.buildGoModule {
      pname = "tui-launcher";
      version = "0.1";

      src = ./.;

      vendorHash = null;

      nativeBuildInputs = with pkgs; [go];

      buildInputs = with pkgs; [
        go
        bubbletea
        fuzzy
      ];
    };

    apps.x86_64-linux.default = {
      type = "app";
      program = "${self.packages.x86_64-linux.default}/bin/tui-launcher";
    };

    devShells.x86_64-linux.default = pkgs.mkShell {
      buildInputs = with pkgs; [go bubbletea fuzzy];
    };
  };
}
