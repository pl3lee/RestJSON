import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { getJSONMetadata, renameJSONFile } from "@/lib/api/jsonFiles";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlertCircle, Check, Edit2, FileJson } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";
import { ApiRouteDialog } from "./api-route-dialog";
import { DeleteFileButton } from "./delete-file-button";

interface JsonFileTopbarProps {
    fileId: string;
    saved: boolean;
}

export default function JsonFileTopbar({ fileId, saved }: JsonFileTopbarProps) {
    const navigate = useNavigate();
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
        onError: (error) => {
            toast.error(error.message);
        },
    });
    const [isRenaming, setIsRenaming] = useState(false);

    const [nameInput, setNameInput] = useState(
        jsonMetadata ? jsonMetadata.fileName : "",
    );

    const handleRenameKeyDown = (e: React.KeyboardEvent | React.FocusEvent) => {
        if (
            e.type === "keydown" &&
            (e as React.KeyboardEvent).key === "Escape"
        ) {
            setNameInput(jsonMetadata!.fileName);
            setIsRenaming(false);
        } else if (
            e.type === "blur" ||
            (e.type === "keydown" && (e as React.KeyboardEvent).key === "Enter")
        ) {
            if (nameInput === "") {
                toast.error("File name cannot be empty!");
                return;
            }
            renameMutation.mutate({
                name: nameInput,
                fileId,
            });
            setIsRenaming(false);
        }
    };

    return (
        <div className="flex justify-between w-full gap-2 px-4 py-2 border-b bg-background items-center flex-row">
            <div className="flex items-center gap-2">
                <FileJson className="h-5 w-5 text-muted-foreground" />

                {isRenaming ? (
                    <Input
                        value={nameInput}
                        onChange={(e) => setNameInput(e.target.value)}
                        onBlur={handleRenameKeyDown}
                        onKeyDown={handleRenameKeyDown}
                        className="h-8"
                        autoFocus
                    />
                ) : jsonMetadataLoading ? (
                    <Skeleton className="h-8 w-32" />
                ) : (
                    <div className="flex items-center gap-2 flex-grow max-w-[30dvw]">
                        <span className="text-sm font-medium md:text-base text-nowrap overflow-hidden text-ellipsis">
                            {jsonMetadata?.fileName}
                        </span>
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => setIsRenaming(true)}
                            className="h-4 w-4"
                        >
                            <Edit2 className="h-4 w-4" />
                            <span className="sr-only">Rename</span>
                        </Button>
                    </div>
                )}
                <div className="flex items-center gap-2 z-50">
                    {saved ? (
                        <div className="flex items-center gap-1">
                            <Check className="h-4 w-4 text-green-500" />
                            <span className="text-sm text-green-500 font-medium">
                                Saved
                            </span>
                        </div>
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
                <ApiRouteDialog fileId={fileId} />
                {jsonMetadata && (
                    <DeleteFileButton
                        fileId={jsonMetadata!.id}
                        onDeleteSuccess={() => navigate("/app")}
                    />
                )}
            </div>
        </div>
    );
}
