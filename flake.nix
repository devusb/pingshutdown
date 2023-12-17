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
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          pingshutdown = pkgs.buildGoModule {
            pname = "pingshutdown";
            inherit version;
            
            src = ./.;
            vendorHash = "sha256-n0WW0DuNo5gyhYFWVdzJHS9MTCVRjy1zwd1UydGlqGQ=";
          };
          default = pingshutdown;
        });
      
      devShells = forAllSystems (system:
        let 
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];

            PINGSHUTDOWN_DELAY = "15s";
            PINGSHUTDOWN_TARGET = "192.0.2.100";

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
