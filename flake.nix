{
  description = "dev shell";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/release-23.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    let overlays = [ ];
    in flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit overlays system;
        };
      in rec {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            gf
            udev
            pkg-config
            clang
          ];
        };
        devShell = self.devShells.${system}.default;
        packages = {
          default = pkgs.callPackage ./default.nix { };
        };
      });
}
