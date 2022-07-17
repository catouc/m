{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils/7e5bf3925";
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };

  description = "m - a whole new organisation";

  outputs = { self, nixpkgs, flake-utils }:
  flake-utils.lib.eachDefaultSystem (system:
  let
    pkgs = nixpkgs.legacyPackages.${system};
  in
    rec {
      packages = flake-utils.lib.flattenTree {
        m = pkgs.buildGo118Module {
          pname = "m";
          version = "v0.1.0";
          modSha256 = pkgs.lib.fakeSha256;
          vendorSha256 = null;
          src = ./.;

          meta = {
            description = "A CLI that manages content";
            homepage = "https://github.com/catouc/m";
            license = pkgs.lib.licenses.mit;
            maintainers = [ "catouc" ];
            platforms = pkgs.lib.platforms.linux;
          };
        };
      };

      defaultPackage = packages.m;
      defaultApp = packages.m;

      devShell = pkgs.mkShell {
        buildInputs = [
          pkgs.go_1_18
          pkgs.buf
          pkgs.protoc-gen-go
          pkgs.protoc-gen-connect-go
        ];
      };
    }
  );
}
