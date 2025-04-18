import { DeleteFileButton } from "@/components/delete-file-button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { useAuth } from "@/hooks/useAuth";
import { createJSONFile, getAllJSONMetadata } from "@/lib/api/jsonFiles";
import type { FileMetadata } from "@/lib/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import {
    AlertCircle,
    CalendarDays,
    File,
    FileIcon,
    PlusCircle,
} from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { toast } from "sonner";

export function App() {
    const navigate = useNavigate();

    const {
        user,
        isLoading: isLoadingUser,
        isLoggedIn,
        isError: userError,
    } = useAuth();
    const {
        data: jsonFiles,
        isLoading: isLoadingFiles,
        isError: filesError,
    } = useQuery({
        queryKey: ["jsonfiles"],
        queryFn: getAllJSONMetadata,
        enabled: !!user,
    });
    if (!isLoggedIn || userError) {
        navigate("/auth");
        return null;
    }

    if (filesError) {
        return (
            <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>
                    Cannot fetch your files, please refresh or try again later.
                </AlertDescription>
            </Alert>
        );
    }

    return (
        <div className="flex flex-col gap-5">
            <title>Your JSON files - RestJSON</title>
            {isLoadingUser || !jsonFiles ? (
                <Skeleton className="h-8 w-full" />
            ) : (
                <>
                    {jsonFiles.length > 0 && <NewJsonForm firstFile={false} />}
                    <div className="flex flex-col gap-2">
                        {isLoadingFiles ? (
                            <>
                                <Skeleton className="h-24 w-full" />
                                <Skeleton className="h-24 w-full" />
                                <Skeleton className="h-24 w-full" />
                            </>
                        ) : (
                            <JsonFilesList jsonFiles={jsonFiles} />
                        )}
                    </div>
                </>
            )}
        </div>
    );
}

function NewJsonForm({ firstFile }: { firstFile: boolean }) {
    const queryClient = useQueryClient();
    const [newFileName, setNewFileName] = useState("");
    const createMutation = useMutation({
        mutationFn: createJSONFile,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["jsonfiles"] });
            setNewFileName("");
        },
        onError: (error) => {
            toast.error(error.message);
        },
    });
    if (firstFile) {
        return (
            <div className="flex flex-col items-center justify-center py-20 text-center">
                <div className="mb-4 rounded-full bg-primary/10 p-3">
                    <FileIcon className="h-10 w-10 text-primary" />
                </div>
                <h2 className="text-2xl font-semibold tracking-tight mb-2">
                    No JSON files yet
                </h2>
                <p className="text-muted-foreground mb-6 max-w-md">
                    Create your first JSON file to get started with RestJSON.
                </p>
                <div className="max-w-sm w-full">
                    <form
                        className="flex gap-2"
                        onSubmit={(e) => {
                            e.preventDefault();
                            if (newFileName === "") {
                                toast.error("File name cannot be empty!");
                                return;
                            }
                            createMutation.mutate(newFileName);
                        }}
                    >
                        <Input
                            type="text"
                            placeholder="Enter file name"
                            value={newFileName}
                            onChange={(e) => setNewFileName(e.target.value)}
                        />
                        <Button
                            type="submit"
                            disabled={createMutation.isPending}
                        >
                            <PlusCircle className="h-4 w-4 mr-2" />
                            Create
                        </Button>
                    </form>
                </div>
            </div>
        );
    }

    return (
        <form
            className="flex flex-row w-full gap-2"
            onSubmit={(e) => {
                e.preventDefault();
                if (newFileName === "") {
                    toast.error("File name cannot be empty!");
                    return;
                }
                createMutation.mutate(newFileName);
            }}
        >
            <Input
                value={newFileName}
                onChange={(e) => setNewFileName(e.target.value)}
                placeholder="Enter new file name"
            />
            <Button type="submit" disabled={createMutation.isPending}>
                <PlusCircle className="h-4 w-4 mr-2" />
                Create JSON File
            </Button>
        </form>
    );
}

function JsonFilesList({ jsonFiles }: { jsonFiles: FileMetadata[] }) {
    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return format(date, "MMM d, yyyy 'at' h:mm a");
    };
    if (jsonFiles.length === 0) {
        return <NewJsonForm firstFile={true} />;
    }

    return jsonFiles.map((file) => (
        <Card key={file.id}>
            <CardContent className="flex flex-row justify-between items-center">
                <Link to={`/app/jsonfile/${file.id}`} className="w-full h-full">
                    <div className="flex items-center gap-2">
                        <File className="h-5 w-5 text-primary" />
                        <span className="font-medium">{file.fileName}</span>
                    </div>
                    <div className="mt-2 text-sm text-foreground">
                        <div className="flex items-center gap-1">
                            <CalendarDays className="h-3.5 w-3.5" />
                            <span>Modified: {formatDate(file.modifiedAt)}</span>
                        </div>
                    </div>
                </Link>
                <DeleteFileButton fileId={file.id} />
            </CardContent>
        </Card>
    ));
}
