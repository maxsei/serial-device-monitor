{
  description = "dev shell";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/release-23.05";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    let overlays = [ ];
    in flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit overlays system; };
        lib = nixpkgs.lib;
      in rec {
        devShells.default = with pkgs;
          mkShell {
            nativeBuildInputs = [ gf udev pkg-config clang_16 ];
            shellHook = ''
              export PATH="${clang-tools_16}/bin:$PATH"
            '';
          };
        devShell = self.devShells.${system}.default;
      });
}
