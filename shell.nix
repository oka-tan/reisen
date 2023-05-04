{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
	    nativeBuildInputs = [ pkgs.go_1_19 pkgs.graphviz pkgs.gv ];
}
