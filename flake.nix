{
  description = "scalar monorepo dev shell";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      # change this if you are not on x86_64 linux
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          delve
          flutter
          git
        ];

        # keep go build cache inside the repo, not in ~
        shellHook = ''
          export GOPATH="$PWD/.go"
          export PATH="$GOPATH/bin:$PATH"
        '';
      };
    };
}
