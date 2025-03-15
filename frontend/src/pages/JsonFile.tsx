import { JsonFileEditor } from "@/components/json-file-editor";
import JsonFileTopbar from "@/components/json-file-topbar";
import { getJSONMetadata, renameJSONFile } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router";
import { toast } from "sonner";

export function JsonFile() {
	const { fileId } = useParams();
	const queryClient = useQueryClient();
	const { data: jsonMetadata, isLoading: jsonMetadataLoading } = useQuery({
		queryKey: [`jsonmetadata-${fileId}`],
		queryFn: async () => await getJSONMetadata(fileId!),
		enabled: !!fileId,
	});

	const renameMutation = useMutation({
		mutationFn: renameJSONFile,
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: [`jsonmetadata-${fileId}`],
			});
			toast.success("Renamed file successfully");
		},
	});

	if (jsonMetadataLoading) {
		return <div>Loading...</div>;
	}

	return (
		<div className="flex flex-col gap-2">
			<JsonFileTopbar
				fileId={fileId!}
				fileName={jsonMetadata!.fileName}
				onRename={(name: string) =>
					renameMutation.mutate({ fileId: fileId!, name })
				}
				onFormat={() => {}}
			/>

			<JsonFileEditor fileId={fileId!} />
		</div>
	);
}
