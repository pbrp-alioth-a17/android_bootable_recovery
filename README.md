**Pitch Black Recovery Project (PBRP) <br>**
**Based on TWRP**

You can find a compiling guide [here](http://forum.xda-developers.com/showthread.php?t=1943625 "Guide").

## Android 17 device ports

Use the pinned recovery manifest and its own build tree. Device-specific kernels,
HALs, firmware, fstabs and SELinux policy must come from the target device; the
Alioth configuration is not a generic device template.

### Display configuration

Set `TARGET_SCREEN_WIDTH` and `TARGET_SCREEN_HEIGHT` (or `DEVICE_RESOLUTION`).
All automatically selected portrait displays use the bundled `portrait_hdpi`
layout, whose declared 1080×1920 coordinates are scaled by the theme engine.
This includes widths below 720 pixels; no separate `portrait_mdpi` assets exist.
`TW_CUSTOM_THEME` continues to select a device-provided layout. Landscape and
square/watch displays require a suitable custom theme; the bundled theme does
not claim support for those shapes. Configure existing `TW_X_OFFSET`,
`TW_Y_OFFSET`, `TW_W_OFFSET` and `TW_H_OFFSET` where the framebuffer needs them.
Validate text, touch targets and keyboard placement on the target display.
Original PBRP, bundled Material and external themes retain separate resources.

### Encryption compatibility

The shared recovery and recovery-local vold retain FBE policy v1/v2 handling,
metadata-encrypted data support, and Gatekeeper/Weaver integration. Password-data
parsing accepts the older serialized fields and ignores newer optional trailing
fields, matching Android's forward-compatible format. Weaver records contain a
version byte and a big-endian slot integer; protector filenames can be padded
to 16 hexadecimal digits. Malformed/truncated records must fail without reading
past their buffers. These paths read existing credentials and keys; they must
not migrate stored key blobs merely to make recovery decryption succeed.

This does not make every encryption implementation interchangeable. A port needs
matching kernel filesystem/crypto support, vendor security services, fstab
options, init configuration, policy and compatible OS/security-patch metadata.
Legacy full-disk encryption (FDE) has no working backend in this branch and must
not be advertised as supported. Vendor-specific wrapped-key or credential
formats require separate implementation and validation. Never offer formatting
as a required verification step for existing user data.

The six-resolution theme selection and sanitizer-backed serialized credential
checks are host validation only. Build and boot each device port, then verify
its actual credential types, decrypted filenames, file access, MTP, display and
touch behavior before declaring that port supported.
