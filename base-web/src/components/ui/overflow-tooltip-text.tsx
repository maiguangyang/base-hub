import * as React from "react"

import { cn } from "@/lib/utils"

import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./tooltip"

function getTextContent(node: React.ReactNode): string {
  if (typeof node === "string" || typeof node === "number") {
    return String(node)
  }

  if (Array.isArray(node)) {
    return node.map(getTextContent).join("")
  }

  if (React.isValidElement<{ children?: React.ReactNode }>(node)) {
    return getTextContent(node.props.children)
  }

  return ""
}

type OverflowTooltipTextProps = Omit<React.HTMLAttributes<HTMLElement>, "children" | "title"> & {
  as?: "span" | "div" | "p" | "code"
  children: React.ReactNode
  tooltipContent?: React.ReactNode
}

function OverflowTooltipText({
  as = "span",
  className,
  children,
  tooltipContent,
  ...props
}: OverflowTooltipTextProps) {
  const Component = as
  const fallbackTitle = getTextContent(tooltipContent ?? children).trim()

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Component
            data-slot="overflow-tooltip-text"
            className={cn("block min-w-0 max-w-full truncate", className)}
            title={fallbackTitle || undefined}
            {...props}
          >
            {children}
          </Component>
        </TooltipTrigger>
        <TooltipContent side="top" className="max-w-[32rem] break-all text-left leading-5">
          {tooltipContent ?? children}
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}

export { OverflowTooltipText }
