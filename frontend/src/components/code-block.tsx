import { Button } from "@/components/ui/button";
import { Check, Copy } from "lucide-react";
import { useState } from "react";

interface CodeBlockProps {
	code: string;
	className?: string;
}

export default function CodeBlock({ code, className }: CodeBlockProps) {
	const [copied, setCopied] = useState(false);

	const copyToClipboard = async () => {
		try {
			await navigator.clipboard.writeText(code);

			setCopied(true);
			setTimeout(() => setCopied(false), 2000);
		} catch (err) {
			console.error("Failed to copy text: ", err);
		}
	};

	return (
		<div className={`relative rounded-md border bg-muted ${className}`}>
			<Button
				variant="ghost"
				size="icon"
				className="absolute right-2 top-2 h-8 w-8 opacity-70 hover:opacity-100"
				onClick={copyToClipboard}
				aria-label="Copy code"
			>
				{copied ? (
					<Check className="h-4 w-4 text-green-500" />
				) : (
					<Copy className="h-4 w-4" />
				)}
			</Button>

			<pre className="overflow-x-auto p-4 text-sm font-mono pr-12">
				<code className="whitespace-pre-wrap break-all">{code}</code>
			</pre>
		</div>
	);
}
