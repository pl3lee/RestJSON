import { useTheme } from "@/components/theme-provider";
import { getJSONFile, updateJSONFile } from "@/lib/api";
import Editor from "@monaco-editor/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";
import { Card, CardHeader } from "./ui/card";
import { AlertCircle, Check } from "lucide-react";

export function JsonFileEditor({ fileId }: { fileId: string }) {
	const { theme } = useTheme();
	const queryClient = useQueryClient();
	const [saved, setSaved] = useState(true);
	const { data: jsonFile, isLoading: jsonFileLoading } = useQuery({
		queryKey: [`jsonfile-${fileId}`],
		queryFn: async () => await getJSONFile(fileId!),
		enabled: !!fileId,
	});
	const jsonString = jsonFile ? JSON.stringify(jsonFile, null, 2) : "";

	const updateMutation = useMutation({
		mutationFn: (variables: { fileId: string; contents: unknown }) =>
			updateJSONFile(variables),
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: [`jsonfile-${fileId}`],
			});
			setSaved(true);
			toast.success("Changes saved successfully.");
		},
	});

	const handleEditorChange = (value: string | undefined) => {
		setSaved(false);
		if (value !== undefined) {
			const normalizedValue = JSON.stringify(JSON.parse(value));
			const normalizedCurrentContents = JSON.stringify(jsonFile);
			if (normalizedValue === normalizedCurrentContents) {
				setSaved(true);
				return;
			}
		}
		debouncedHandleEditorChange(value);
	};
	const debouncedHandleEditorChange = useDebouncedCallback(
		(value: string | undefined) => {
			if (value !== undefined) {
				try {
					const contents = JSON.parse(value);
					updateMutation.mutate({ fileId: fileId!, contents });
				} catch (e) {
					console.error(e);
					return;
				}
			}
		},
		1000,
	);
	if (jsonFileLoading) {
		return <div>Loading...</div>;
	}

	return (
		<Card>
			<CardHeader className="flex flex-row items-center justify-end">
				{saved ? (
					<>
						<Check className="h-4 w-4 text-green-500" />
						<span className="text-sm text-green-500 font-medium">Saved</span>
					</>
				) : (
					<>
						<AlertCircle className="h-4 w-4 text-amber-500" />
						<span className="text-sm text-amber-500 font-medium">Unsaved</span>
					</>
				)}
			</CardHeader>
			<Editor
				height="90vh"
				defaultLanguage="json"
				defaultValue={jsonString}
				theme={theme === "light" ? theme : "vs-dark"}
				onChange={handleEditorChange}
			/>
		</Card>
	);
}
