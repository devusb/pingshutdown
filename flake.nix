{
  description = "A simple Go package";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    rec {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          pingshutdown = pkgs.callPackage ./package.nix { inherit version; }; 
          default = pingshutdown;
        });


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

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
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
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.pingshutdown);
    };
}
