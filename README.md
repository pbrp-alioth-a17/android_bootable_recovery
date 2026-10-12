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

## Standard A17 build workflow

Use the recovery checkout for source, dependencies and outputs:

```bash
source build/envsetup.sh
list_products
list_releases pb_alioth
lunch pb_alioth cp2a user
m -j4 bootimage
```

Replace the product, release and variant with those declared by the target
port. A17's full target spelling (`pb_alioth-cp2a-user`) also works. Register
products in the device's `AndroidProducts.mk`; use its actual release config.
Do not assume that every device uses Alioth's product, kernel or Boot layout.

The recovery layer in `vendor/lineage` loads these additional commands without changing `lunch` or
`m`, modifying platform source files, selecting another compiler, or accessing
a phone:

```bash
breakfast pb_alioth cp2a user       # same lunch selection
recoverybuild bootimage -j4        # wrapper around m -j4 bootimage
recoveryout                        # selected product output directory
```

Choose `bootimage`, `recoveryimage` or `vendorbootimage` according to the device
port's partition layout. The helper requires an explicit image target and does
not infer it from a codename. A generated image is not automatically a flashable
release: run the device's packaging and validation workflow.

For Alioth, the existing metadata-aware packager and installer validator remain
required. From the checkout root, for example:

```bash
recoverypackage device/xiaomi/alioth/tools/package-recovery.py \
  --metadata-boot /path/to/matched-boot.img --security-patch 2026-10
recoverypackage device/xiaomi/alioth/tools/package-installer.py
```

Use verified ROM metadata, not the sample patch value, for an actual release.
Other ports supply their own packaging script and arguments; `recoverypackage`
accepts an explicit repository-relative Python script and forwards arguments.
It does not reuse Alioth's validation or partition assumptions for another
product. None of these helpers uploads, boots or flashes an image.

## PBRP layer on vendor/lineage

Recovery products inherit `vendor/lineage/config/recovery.mk`, not the Lineage
ROM common product. Use `vendor/lineage/config/BoardConfigRecovery.mk` for the
shared recovery board integration, or its scoped Soong/policy includes when a
device already configures the Lineage kernel rules. Recovery tools, versioning,
PBRP flag namespaces and recovery-only policy remain separate from ROM defaults.
The A17 kernel/header infrastructure and standard Lineage tools are reused.

`vendor/lineage/pb_build.sh` is a manifest link to the local A17 wrapper:

```bash
vendor/lineage/pb_build.sh pb_alioth cp2a user bootimage -j4
```

It selects an explicit product/release/variant and image target. It does not set
credentials, synchronize a different checkout, upload, reboot or flash. The
manifest also keeps the TWRP common-product link pointing to the recovery layer.
The obsolete `vendor/pb` project is no longer part of the build tree.
