{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { nixpkgs, ... }: 
    let
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin"];
      eachSystem = nixpkgs.lib.genAttrs systems;
    in
  {
    devShell = eachSystem(system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      pkgs.mkShell {
        packages = [
          pkgs.go
          pkgs.golangci-lint
          pkgs.opentofu
          pkgs.proton-pass-cli
        ];
      }
    );
  };
}
