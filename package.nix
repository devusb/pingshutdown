{ buildGoModule, version }:
buildGoModule {
  pname = "pingshutdown";
  inherit version;

  src = ./.;
  vendorHash = "sha256-n0WW0DuNo5gyhYFWVdzJHS9MTCVRjy1zwd1UydGlqGQ=";

  meta.mainProgram = "pingshutdown";
}
