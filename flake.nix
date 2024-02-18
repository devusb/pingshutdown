{
  description = "A simple Go package";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixpkgs-unstable";
  inputs.flake-parts.url = "github:hercules-ci/flake-parts";
  inputs.hercules-ci-effects.url = "github:hercules-ci/hercules-ci-effects";

  outputs = inputs@{ self, nixpkgs, flake-parts, hercules-ci-effects }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
    in
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        hercules-ci-effects.flakeModule
      ];

      systems = supportedSystems;
      herculesCI.ciSystems = [ "x86_64-linux" ];

      flake = {
        overlays = {
          default = final: prev: {
            pingshutdown = prev.callPackage ./package.nix { inherit version; };
          };
        };


        nixosModules = {
          pingshutdown = {
            imports = [
              ./module.nix
            ];

            nixpkgs.overlays = [
              self.overlays.default
            ];
          };
        };
      };

      perSystem = { pkgs, ... }: {
        packages = rec {
          pingshutdown = pkgs.buildGoModule {
            pname = "pingshutdown";
            inherit version;

            src = ./.;
            vendorHash = "sha256-n0WW0DuNo5gyhYFWVdzJHS9MTCVRjy1zwd1UydGlqGQ=";
          };
          default = pingshutdown;
        };

        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];

            PINGSHUTDOWN_DELAY = "15s";
            PINGSHUTDOWN_TARGET = "192.0.2.100";
            PINGSHUTDOWN_NOTIFICATION = "true";
            PINGSHUTDOWN_DRYRUN = "true";

            shellHook = ''
              if [ -f ".env" ]; then
                source .env
              fi
            '';
          };
        };
      };

    };
}
