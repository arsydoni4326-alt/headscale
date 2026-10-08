# Device Registration

This guide covers how to register new devices to your Headscale network.

## QR Code Registration

As an alternative to the CLI command, you can register devices by scanning a QR code:

1. Start Tailscale on your device to initiate registration
2. Note the registration URL displayed
3. Open Headplane web UI
4. Navigate to **Machines → Add Device → Scan QR**
5. Select the user/namespace for the device
6. Scan the QR code shown on the registration page
7. Device is automatically approved

### Requirements

- HTTPS connection (required for camera access)
- Browser with camera support (Chrome, Firefox, Safari, Edge)
- Camera permissions granted

### Security

- QR codes expire with the registration session (typically 5-10 minutes)
- QR codes are single-use only
- No secrets are embedded in QR codes

## CLI Registration

To register a device using the CLI:

```bash
headscale auth register --auth-id <hskey> --user USERNAME
```

The `auth-id` is provided in the registration URL when you start Tailscale on the device.
