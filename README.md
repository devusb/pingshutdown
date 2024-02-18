# pingshutdown
`pingshutdown` can be used to shutdown a host when another remote host is no longer responding to ping.

A potential use of this is safely shutting down a host connected to a UPS on power failure. While many UPSes work with `apcupsd` or `nut` to provide this notification, `pingshutdown` can be used in any hardware configuration. 

`pingshutdown` also provides the ability to notify on pending shutdown via [Pushover](https://pushover.net), and offers a web interface to see timer status as well as temporarily disable the shutdown functionality.

## Installation
`pingshutdown` is intended to run as a `systemd` unit, either installed manually or using the included NixOS module. 

### Manual Installation
To install manually, build the application as 
```
go build ./cmd/pingshutdown
```
and put the binary in your `PATH` (e.g. `/usr/local/bin`). 

Then, copy the included `pingshutdown.service` into `/etc/systemd/system` and create an `EnvironmentFile` in `/etc/default/pingshutdown.env` containing your desired configuration (see `pingshutdown.env` for example syntax).

### NixOS Module Installation
To install using the NixOS module, add this repository to your `flake.nix` as
```
inputs.pingshutdown.url = "github:devusb/pingshutdown";
```
and import the module in your NixOS configuration as
```
imports = [
    inputs.pingshutdown.nixosModules.pingshutdown
];
```
then, enable the service in your NixOS configuration as
```
services.pingshutdown = {
  enable = true;
  environmentFile = /run/secrets/pushover;
  settings = {
    PINGSHUTDOWN_DELAY = "10m";
    PINGSHUTDOWN_TARGET = "192.168.20.1";
    PINGSHUTDOWN_NOTIFICATION = "true";
    PINGSHUTDOWN_DRYRUN = "false";
    PINGSHUTDOWN_STATUSPORT = "9081";
  };
};
```
where `environmentFile` can contain additional settings such as a Pushover token to be used for notifications.

## Configuration
Configuration is via environment variables, options are listed below

- `PINGSHUTDOWN_DELAY` - amount of time to wait before initiating system shutdown (default `5m`)
- `PINGSHUTDOWN_TARGET` - remote host to ping (default `www.google.com`)
- `PINGSHUTDOWN_NOTIFICATON` - enable notification of shutdown status via Pushover (default `false`)
- `PINGSHUTDOWN_NOTIFICATONTOKEN` - Pushover application API token
- `PINGSHUTDOWN_NOTIFICATIONUSER` - Pushover user key to be notified
- `PINGSHUTDOWN_DRYRUN` - whether to actually shut down machine -- `true` will only initiate countdown and notification, but will not shut down (default `false`)
- `PINGSHUTDOWN_STATUSPORT` - port to serve status web interface (default `8081`)

## Usage
When started, `pingshutdown` will begin pinging the target, and upon ping failures occuring for 10 seconds, will begin a timer after which the host system will be shut down. 

When the timer begins, an optional notification will be set to the Pushover user specified via the `PINGSHUTDOWN_NOTIFICATONUSER` environment variable. 

A web interface is available at the port configure via `PINGSHUTDOWN_STATUSPORT` that shows timer status and enables locking out (disabling) the shutdown functionality, as might be necessary during network maintenance.
