import { JsonFileEditor } from "@/components/json-file-editor";
import JsonFileTopbar from "@/components/json-file-topbar";
import { getJSONMetadata, renameJSONFile } from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useParams } from "react-router";
import { toast } from "sonner";

export function JsonFile() {
	const { fileId } = useParams();
	const queryClient = useQueryClient();
	const [saved, setSaved] = useState(true);
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
		onError: (error) => {
			toast.error(error.message);
		},
	});

	if (jsonMetadataLoading) {
		return <div>Loading...</div>;
	}

	return (
		<div className="flex flex-col gap-2">
			<title>{`${jsonMetadata!.fileName} - RestJSON`}</title>
			<JsonFileTopbar fileId={fileId!} saved={saved} />

			<JsonFileEditor fileId={fileId!} setSaved={setSaved} />
		</div>
	);
}
