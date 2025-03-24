import { useTheme } from "@/components/theme-provider";
import { getJSONFile, updateJSONFile } from "@/lib/api";
import Editor from "@monaco-editor/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useDebouncedCallback } from "use-debounce";
import { Skeleton } from "./ui/skeleton";

export function JsonFileEditor({
	fileId,
	setSaved,
}: { fileId: string; setSaved: (saved: boolean) => void }) {
	const { theme } = useTheme();
	const queryClient = useQueryClient();
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
			queryClient.invalidateQueries({
				queryKey: [`dynamic-${fileId}`],
			});
			setSaved(true);
			toast.success("Changes saved successfully.");
		},
		onError: (error) => {
			toast.error(error.message);
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
		return <Skeleton className="h-90vh w-full" />;
	}

	return (
		<div>
			<Editor
				height="90vh"
				defaultLanguage="json"
				defaultValue={jsonString}
				theme={theme === "light" ? theme : "vs-dark"}
				onChange={handleEditorChange}
				options={{
					minimap: { enabled: false },
					formatOnPaste: true,
					formatOnType: true,
					scrollBeyondLastLine: false,
					automaticLayout: true,
					tabSize: 2,
					wordWrap: "on",
					wrappingIndent: "deepIndent",
				}}
			/>
		</div>
	);
}
