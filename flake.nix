{
  description = "A very basic flake";

  inputs = {
      go-dev.url = "github:kijjuy/nix-flakes?dir=go";
  };

  outputs = { self, go-dev }: go-dev.outputs;
}
