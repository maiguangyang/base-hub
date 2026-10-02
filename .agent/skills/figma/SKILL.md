---
name: figma
description: Use the Figma MCP server to fetch design context, screenshots, variables, and assets from Figma, and to translate Figma nodes into production code. Trigger when a task involves Figma URLs, node IDs, design-to-code implementation, or Figma MCP setup and troubleshooting.
---

# Figma MCP

Use the Figma MCP server for Figma-driven implementation. For setup and debugging details (env vars, config, verification), see `references/figma-mcp-config.md`.

## Figma MCP Integration Rules
These rules define how to translate Figma inputs into code for this project and must be followed for every Figma-driven change.

### Required flow (do not skip)
1. Extract the `fileKey` and `nodeId` from the provided Figma URL.
2. Run `get_figma_data` to fetch comprehensive Figma file data including layout, content, visuals, and component information for the exact node(s).
3. If the design contains images or vector icons that cannot be created with CSS, run `download_figma_images` to save SVG/PNG assets directly to the project directory.
4. Translate the component data into this project's React + Tailwind conventions. Reuse the project's color tokens, components, and typography wherever possible.
6. Validate against Figma for 1:1 look and behavior before marking complete.

### Implementation rules
- Treat the Figma MCP output (React + Tailwind) as a representation of design and behavior, not as final code style.
- Replace Tailwind utility classes with the project's preferred utilities/design-system tokens when applicable.
- Reuse existing components (e.g., buttons, inputs, typography, icon wrappers) instead of duplicating functionality.
- Use the project's color system, typography scale, and spacing tokens consistently.
- Respect existing routing, state management, and data-fetch patterns already adopted in the repo.
- Strive for 1:1 visual parity with the Figma design. When conflicts arise, prefer design-system tokens and adjust spacing or sizes minimally to match visuals.
- Validate the final UI against the Figma data for both look and behavior.

### Asset handling
- Use the `download_figma_images` tool to proactively pull down necessary SVG icons or raster images mentioned in the Figma node.
- IMPORTANT: DO NOT import/add new generic icon packages if the specific icon exists in Figma.
- Pass the correct absolute `localPath` to the tool so assets are saved directly into the frontend project's public or assets directory.

### Link-based prompting
- The server is link-based: copy the Figma frame/layer link and give that URL to the MCP client when asking for implementation help.
- The client cannot browse the URL but extracts the node ID from the link; always ensure the link points to the exact node/variant you want.

### Conditional pairing with `$ui-ux-pro-max`
- Default behavior: when the user provides a specific Figma frame or node and expects faithful design-to-code implementation, use `$figma` alone as the source-of-truth workflow.
- Also invoke `$ui-ux-pro-max` only when at least one of these is true:
  - the user explicitly asks for UX polish, visual refinement, responsiveness improvements, or accessibility improvements
  - the Figma design does not fully specify interactions, states, or responsive behavior, and implementation requires product/UI judgment
  - the design must be adapted into the repo's existing design system, component library, or token constraints beyond straight visual translation
  - the Figma-based implementation is complete and now needs a follow-up UI/UX review pass
- When both skills are used, the order is:
  1. Use `$figma` first to fetch the design context and implement the closest faithful baseline.
  2. Use `$ui-ux-pro-max` second to refine interaction quality, responsiveness, accessibility, motion, hierarchy, and polish.
  3. If a proposed refinement would materially diverge from the approved Figma intent, keep the Figma behavior/layout by default unless the user explicitly asked for a redesign or broader UX improvement.
- Do not let `$ui-ux-pro-max` silently turn a Figma-accurate implementation request into an open-ended redesign task.

## References
- `references/figma-mcp-config.md` — setup, verification, troubleshooting, and link-based usage reminders.
- `references/figma-tools-and-prompts.md` — tool catalog and usage parameters for `get_figma_data` and `download_figma_images`.
