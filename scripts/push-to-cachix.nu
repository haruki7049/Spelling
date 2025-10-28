#!/usr/bin/env nu

if ($env | get CACHIX_AUTH_TOKEN? | is-empty) {
  error make { msg: "There are no setting for CACHIX_AUTH_TOKEN environment variable." }
}

nix flake archive --json
  | from json
  | [ ($in | get path), ($in | get inputs | values | get path) ]
  | flatten
  | str join (char nl)
  | cachix push spelling
