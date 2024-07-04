import { defineConfig } from "astro/config";
import tailwind from "@astrojs/tailwind";
import mdx from "@astrojs/mdx";
import sitemap from "@astrojs/sitemap";
import embeds from "astro-embed/integration";
import sentry from "@sentry/astro";
import spotlightjs from "@spotlightjs/astro";
import expressiveCode from "astro-expressive-code";
import { pluginCollapsibleSections } from "@expressive-code/plugin-collapsible-sections";
import react from "@astrojs/react";

import icon from "astro-icon";

// https://astro.build/config
export default defineConfig({
  site:
    process.env.NODE_ENV === "development"
      ? "http://localhost:4321"
      : "https://david.cotelessa.com",
  integrations: [
    react(),
    tailwind(),
    sitemap(),
    embeds(),
    expressiveCode({
      plugins: [pluginCollapsibleSections()],
    }),
    mdx(),
    sentry(),
    spotlightjs(),
    icon(),
  ],
});
