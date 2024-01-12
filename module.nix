{ config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.pingshutdown;
in
{
  options.services.pingshutdown = with types; {
    enable = mkEnableOption (mdDoc "Shut down host when remote connectivity not available");
    package = mkPackageOption pkgs "pingshutdown" { };
    environmentFile = mkOption {
      type = nullOr path;
      default = null;
      example = "/run/secrets/pingshutdown";
    };
    settings = mkOption {
      type = types.submodule (settings: {
        freeformType = attrsOf str;
      });
      default = { };
    };
  };

  config = mkIf cfg.enable {
    systemd.services.pingshutdown = {
      description = "pingshutdown";
      after = [ "network.target" "network-online.target" ];
      wants = [ "network-online.target" ];
      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        ExecStart = "${getExe cfg.package}";
        RuntimeDirectory = "pingshutdown";
        RuntimeDirectoryMode = "0700";
        DynamicUser = true;

        EnvironmentFile = [
          cfg.environmentFile
          (pkgs.writeText "pingshutdown-settings" (generators.toKeyValue { } cfg.settings))
        ];
      };
    };
  };
}
