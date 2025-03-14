{
  description = "Gofi - TUI Launcher for NixOS";

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
      pname = "gofi-launcher";
      version = "0.1";

      src = ./.;

      vendorHash = null;

      nativeBuildInputs = with pkgs; [go];

      buildInputs = with pkgs; [
        go
      ];
      proxyVendor = true;
    };

    apps.x86_64-linux.default = {
      type = "app";
      program = "${self.packages.x86_64-linux.default}/bin/gofi-launcher";
    };

    devShells.x86_64-linux.default = pkgs.mkShell {
      buildInputs = with pkgs; [go gotools];
    };
  };
}
