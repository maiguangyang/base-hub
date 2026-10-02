# Figma MCP Tools and Prompts

Quick reference for the Figma MCP toolset and how to use it for design-to-code tasks.

## Core tools

- `get_figma_data`: **Primary tool.** Get comprehensive Figma file data including layout, content, visuals, and component information. 
  - **Inputs Required:** `fileKey` (extracted from the figma.com/design/... URL) and `nodeId` (e.g., `50:14180`).
  - **Usage:** Call this first to understand the structure, styles, texts, and box models of the component.

- `download_figma_images`: **Asset extractor.** Download SVG and PNG images used in a Figma file based on the IDs of image or icon nodes.
  - **Inputs Required:** `fileKey`, `nodes` (list of node objects containing at least `nodeId` and `fileName`), and `localPath` (absolute path to save the files).
  - **Usage:** If `get_figma_data` reveals vector icons or raster images that are too complex for CSS, use this to dump them straight into `web/public/` or `web/src/assets/`.

## URL Parsing Guide

Given a Figma URL:
`https://www.figma.com/design/2AFFxtqoN8FvlEVYJaAqvJ/R-F-App?node-id=50-14180&t=...`

- `fileKey` = `2AFFxtqoN8FvlEVYJaAqvJ`
- `nodeId` = `50:14180` (Note: Replace the hyphen `-` with a colon `:` for the tool arguments).

## Best-practice flow reminder
Use `get_figma_data` → (analyze layout and see if assets are needed) → `download_figma_images` (if needed) → write code.
