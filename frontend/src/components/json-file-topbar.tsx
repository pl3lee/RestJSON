import type React from "react";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Edit2, Code, FileJson } from "lucide-react";
import { toast } from "sonner";

interface JsonFileTopbarProps {
	fileId: string;
	fileName: string;
	onRename: (newName: string) => void;
	onFormat: () => void;
}

export default function JsonFileTopbar({
	fileId,
	fileName,
	onRename = () => {},
	onFormat = () => {},
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
						<div className="space-y-4 py-4 w-full">
							<p className="text-sm text-muted-foreground">
								Use this endpoint to access your JSON data programmatically:
							</p>
							<div className="w-full bg-muted p-3 rounded-md font-mono text-sm overflow-x-auto">
								{endpoint}
							</div>

							<div className="space-y-2">
								<h4 className="text-sm font-medium">Example usage:</h4>
								<pre className="bg-muted p-3 rounded-md text-xs overflow-x-auto">
									{`fetch("${endpoint}")
  .then(response => response.json())
  .then(data => console.log(data))`}
								</pre>
							</div>
						</div>
					</DialogContent>
				</Dialog>

				<Button
					variant="outline"
					size="sm"
					onClick={onFormat}
					className="gap-1"
				>
					<FileJson className="h-4 w-4" />
					Format
				</Button>
			</div>
		</div>
	);
}
