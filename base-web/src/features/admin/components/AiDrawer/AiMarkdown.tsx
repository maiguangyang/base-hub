import Markdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

export function AiMarkdown({ children }: { children: string }) {
  return <div className="min-w-0 break-words text-sm leading-6 [&>*+*]:mt-3">
    <Markdown remarkPlugins={[remarkGfm]} components={{
      h1: ({ children: content }) => <h1 className="text-lg font-semibold">{content}</h1>,
      h2: ({ children: content }) => <h2 className="text-base font-semibold">{content}</h2>,
      h3: ({ children: content }) => <h3 className="font-semibold">{content}</h3>,
      ul: ({ children: content }) => <ul className="list-disc pl-5">{content}</ul>,
      ol: ({ children: content }) => <ol className="list-decimal pl-5">{content}</ol>,
      a: ({ children: content, href }) => <a href={href} target="_blank" rel="noopener noreferrer" className="text-primary underline underline-offset-2">{content}</a>,
      code: ({ children: content }) => <code className="rounded bg-background/70 px-1 py-0.5 font-mono text-xs">{content}</code>,
      pre: ({ children: content }) => <pre className="overflow-x-auto rounded-md bg-background/70 p-3 text-xs">{content}</pre>,
      table: ({ children: content }) => <div className="overflow-x-auto"><table className="w-full border-collapse text-left">{content}</table></div>,
      th: ({ children: content }) => <th className="border border-border px-2 py-1 font-semibold">{content}</th>,
      td: ({ children: content }) => <td className="border border-border px-2 py-1">{content}</td>,
    }}>{children}</Markdown>
  </div>;
}
