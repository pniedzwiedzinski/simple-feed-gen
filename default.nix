let
  pkgs = import <nixpkgs> {};
in
  with pkgs;
  mkShell {
    buildInputs = [
      stdenv
      go
    ];
    shellHook = ''
      export GO111MODULE=on
    '';
  }
