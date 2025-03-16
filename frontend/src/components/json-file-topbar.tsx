import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
	AlertCircle,
	Check,
	Code,
	Edit2,
	FileJson,
	Key,
	KeyRound,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import CodeBlock from "./code-block";
import { Alert, AlertDescription, AlertTitle } from "./ui/alert";
import { Badge } from "./ui/badge";

interface JsonFileTopbarProps {
	fileId: string;
	fileName: string;
	saved: boolean;
	onRename: (newName: string) => void;
}

export default function JsonFileTopbar({
	fileId,
	fileName,
	saved,
	onRename = () => {},
}: JsonFileTopbarProps) {
	const [isRenaming, setIsRenaming] = useState(false);

	const [nameInput, setNameInput] = useState(fileName);

	const handleRename = () => {
		if (isRenaming) {
			if (nameInput === "") {
				toast.error("File name cannot be empty!");
				return;
			}
			onRename(nameInput);
			setIsRenaming(false);
		} else {
			setIsRenaming(true);
		}
	};

	const handleKeyDown = (e: React.KeyboardEvent) => {
		if (e.key === "Enter") {
			onRename(nameInput);
			setIsRenaming(false);
		} else if (e.key === "Escape") {
			setNameInput(fileName);
			setIsRenaming(false);
		}
	};

	const endpoint = `${import.meta.env.VITE_API_URL}/public/${fileId}`;

	return (
		<div className="flex items-center justify-between w-full px-4 py-2 border-b bg-background">
			<div className="flex items-center gap-2">
				<FileJson className="h-5 w-5 text-muted-foreground" />

				{isRenaming ? (
					<Input
						value={nameInput}
						onChange={(e) => setNameInput(e.target.value)}
						onBlur={handleRename}
						onKeyDown={handleKeyDown}
						className="h-8 w-64"
						autoFocus
					/>
				) : (
					<div className="flex items-center gap-2">
						<span className="font-medium">{fileName}</span>
						<Button
							variant="ghost"
							size="icon"
							onClick={handleRename}
							className="h-8 w-8"
						>
							<Edit2 className="h-4 w-4" />
							<span className="sr-only">Rename</span>
						</Button>
					</div>
				)}
				<div className="flex items-center gap-2 z-50">
					{saved ? (
						<Badge variant="outline" className="">
							<Check className="h-4 w-4 text-green-500" />
							<span className="text-sm text-green-500 font-medium">Saved</span>
						</Badge>
					) : (
						<Badge variant="outline" className="">
							<AlertCircle className="h-4 w-4 text-amber-500" />
							<span className="text-sm text-amber-500 font-medium">
								Unsaved
							</span>
						</Badge>
					)}
				</div>
			</div>

			<div className="flex items-center gap-2">
				<Dialog>
					<DialogTrigger asChild>
						<Button variant="outline" size="sm" className="gap-1">
							<Code className="h-4 w-4" />
							API
						</Button>
					</DialogTrigger>
					<DialogContent>
						<DialogHeader>
							<DialogTitle>API Endpoint</DialogTitle>
						</DialogHeader>
						<div className="space-y-4 py-4">
							<p className="text-sm text-muted-foreground">
								Use this endpoint to access your JSON data programmatically,
								including your API key in the Authorization header:
							</p>
							<CodeBlock code={endpoint} />

							<div className="space-y-2">
								<h4 className="text-sm font-medium">Example usage:</h4>
								<CodeBlock
									code={`const res = await fetch("${endpoint}", {
	headers: {
		Authorization: "Bearer YOUR_API_KEY"
	}
})
const data = await res.json()
console.log(data)
`}
								/>
							</div>
							<Alert>
								<Key className="h-4 w-4" />
								<AlertTitle>API Key</AlertTitle>
								<AlertDescription>
									You can get your API key in your account page
								</AlertDescription>
							</Alert>
						</div>
					</DialogContent>
				</Dialog>
			</div>
		</div>
	);
}
