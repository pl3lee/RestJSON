import { getDynamicRoutes } from "@/lib/api/jsonFiles";
import type { Route } from "@/lib/types";
import { cn } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import { AlertCircle, Code, Key } from "lucide-react";
import CodeBlock from "./code-block";
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "./ui/accordion";
import { Alert, AlertDescription, AlertTitle } from "./ui/alert";
import { Badge } from "./ui/badge";
import { Button } from "./ui/button";
import { Card } from "./ui/card";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "./ui/dialog";
import { Skeleton } from "./ui/skeleton";

type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
const methodColors = {
    GET: "bg-emerald-100 text-emerald-700 hover:bg-emerald-200 dark:bg-emerald-900 dark:text-emerald-300",
    POST: "bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900 dark:text-blue-300",
    PUT: "bg-amber-100 text-amber-700 hover:bg-amber-200 dark:bg-amber-900 dark:text-amber-300",
    DELETE: "bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-900 dark:text-red-300",
    PATCH: "bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900 dark:text-purple-300",
};

export function ApiRouteDialog({ fileId }: { fileId: string }) {
    const endpoint = `${import.meta.env.VITE_API_URL}/public/${fileId}`;
    const { data: routesData, isLoading: routesLoading } = useQuery({
        queryFn: async () => await getDynamicRoutes(fileId),
        queryKey: [`dynamic-${fileId}`],
    });
    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="outline" size="sm" className="gap-1">
                    <Code className="h-4 w-4" />
                    API
                </Button>
            </DialogTrigger>
            <DialogContent
                className="overflow-y-auto max-h-[80dvh] min-w-full"
                aria-describedby="dialog-description"
            >
                <DialogHeader>
                    <DialogTitle>API Endpoint</DialogTitle>
                </DialogHeader>
                <div className="space-y-4 py-4" id="dialog-description">
                    <p className="text-sm text-muted-foreground">
                        Use this endpoint to access your JSON data
                        programmatically, including your API key in the
                        Authorization header:
                    </p>
                    <CodeBlock code={endpoint} />
                    <Alert>
                        <Key className="h-4 w-4" />
                        <AlertTitle>API Key</AlertTitle>
                        <AlertDescription>
                            You can get your API key in your account page
                        </AlertDescription>
                    </Alert>

                    <div className="space-y-2">
                        <h4 className="text-base font-medium">
                            Example usage:
                        </h4>
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
                    {routesLoading ? (
                        <Skeleton className="h-24 w-full" />
                    ) : (
                        routesData && (
                            <div className="space-y-2">
                                <h4 className="text-base font-medium">
                                    Dynamic routes generated:
                                </h4>
                                <Alert>
                                    <AlertCircle className="h-4 w-4" />
                                    <AlertTitle>Missing endpoint?</AlertTitle>
                                    <AlertDescription>
                                        Endpoint won't be generated if the
                                        resource name contains whitespace!
                                    </AlertDescription>
                                </Alert>
                                <ApiRouteCollection
                                    fileId={fileId}
                                    routes={routesData}
                                />
                            </div>
                        )
                    )}
                </div>
            </DialogContent>
        </Dialog>
    );
}

function ApiRoute({
    method,
    endpoint,
    description,
    url,
}: {
    method: HttpMethod;
    endpoint: string;
    description: string;
    url: string;
    defaultOpen?: boolean;
}) {
    const methodBadge = (
        <Badge className={cn("font-mono font-bold", methodColors[method])}>
            {method}
        </Badge>
    );
    return (
        <Card className="py-0">
            <Accordion type="single" collapsible>
                <AccordionItem value={`${method}-${endpoint}`}>
                    <AccordionTrigger className="py-4 px-4 hover:no-underline">
                        <div className="flex items-center gap-3 py-0">
                            {methodBadge}
                            <span className="font-medium text-gray-800 dark:text-gray-200">
                                {endpoint}
                            </span>
                        </div>
                    </AccordionTrigger>
                    <AccordionContent className="px-4 pb-4 pt-1">
                        <div className="space-y-3">
                            <p className="text-foreground">{description}</p>
                            <CodeBlock
                                code={url}
                                prependComponent={methodBadge}
                            />
                        </div>
                    </AccordionContent>
                </AccordionItem>
            </Accordion>
        </Card>
    );
}

function ApiRouteCollection({
    fileId,
    routes,
}: { fileId: string; routes: Route[] }) {
    const base = `${import.meta.env.VITE_API_URL}/public/${fileId}`;
    return (
        <div className="space-y-3">
            {routes.map((route) => (
                <ApiRoute
                    key={`${route.method}${base}${route.url}`}
                    method={route.method as HttpMethod}
                    endpoint={route.url}
                    description={route.description}
                    url={`${base}${route.url}`}
                />
            ))}
        </div>
    );
}
