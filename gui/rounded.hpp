/* SPDX-License-Identifier: GPL-3.0-or-later */
#pragma once

#include <algorithm>
#include <cmath>
#include "minuitwrp/minui.h"

// Uses the caller's color and clip; radius is in framebuffer pixels.
inline void gui_fill_rounded(int x, int y, int width, int height, int radius,
                             bool top = true, bool bottom = true)
{
    if (width <= 0 || height <= 0) return;
    radius = std::max(0, std::min(radius, std::min(width, height) / 2));
    gr_fill(x, y + (top ? radius : 0), width,
            height - (top ? radius : 0) - (bottom ? radius : 0));
    for (int row = 0; row < radius; ++row) {
        const double dy = radius - row - 0.5;
        const int inset = radius - static_cast<int>(std::sqrt(radius * radius - dy * dy));
        if (top) gr_fill(x + inset, y + row, width - 2 * inset, 1);
        if (bottom) gr_fill(x + inset, y + height - row - 1, width - 2 * inset, 1);
    }
}
