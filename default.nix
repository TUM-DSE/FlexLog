# we ping nixpkgs sources here according to https://nix.dev/tutorials/towards-reproducibility-pinning-nixpkgs
{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/3c52ea8c9216a0d5b7a7b4d74a9d2e858b06df5c.tar.gz") {}
}:
pkgs.mkShell {
  buildInputs = [
    pkgs.libndctl
    pkgs.glib
  ];
  # add tools here
  nativeBuildInputs = [
    pkgs.gcc10
    pkgs.pkg-config
    pkgs.pandoc
    pkgs.m4
    pkgs.go
    pkgs.tbb
    pkgs.pmdk
    pkgs.fmt
    pkgs.bashInteractive
    pkgs.cmake
    pkgs.dos2unix
    pkgs.binutils-unwrapped
    pkgs.gdb
    pkgs.python
    pkgs.python39Packages.pandas
    pkgs.python39Packages.scikitlearn
    pkgs.python39Packages.shutilwhich
    pkgs.python39Packages.opencv3
    pkgs.valgrind
    pkgs.autoreconfHook
    pkgs.numactl
    (pkgs.python39.withPackages (ps: [
        ps.matplotlib
        ps.pyyaml
    ]))
  ];
}
