# ============================
# Headless Chrome with CJK fonts
# ============================
# The upstream image ships only DejaVu, so Chinese, Japanese and Korean pages
# render as empty boxes. Noto CJK adds the missing glyphs; the upstream
# entrypoint and sandbox flags stay unchanged.
FROM chromedp/headless-shell:latest

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        fontconfig \
        fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/*
