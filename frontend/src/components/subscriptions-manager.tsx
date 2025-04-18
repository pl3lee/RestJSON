import {
    checkout,
    getBillingPortalUrl,
    getSubscriptionStatus,
} from "@/lib/api/payment";
import { useMutation, useQuery } from "@tanstack/react-query";
import { CheckCircle } from "lucide-react";
import { toast } from "sonner";
import { Badge } from "./ui/badge";
import { Button } from "./ui/button";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "./ui/card";
import { Skeleton } from "./ui/skeleton";

export const plans = {
    Monthly: {
        priceId: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "price_1R9FuqDGza3FJhYO6DtchdQh"
            : "",
        price: 5,
    },
    Yearly: {
        priceId: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "price_1R9FwkDGza3FJhYOjdukqxiO"
            : "",
        price: 30,
    },
};
export function SubscriptionsManager() {
    const { data: subscriptionStatus, isLoading: subscriptionStatusLoading } =
        useQuery({
            queryKey: ["subscription-status"],
            queryFn: getSubscriptionStatus,
        });
    const billingPortalMutation = useMutation({
        mutationFn: getBillingPortalUrl,
        onSuccess: (portalUrl: string) => {
            window.location.href = portalUrl;
        },
        onError: () => {
            toast.error("Cannot get management portal url, please retry.");
        },
    });

    if (subscriptionStatusLoading) {
        return (
            <div className="container mx-auto py-10">
                <Skeleton className="h-42 w-full" />
            </div>
        );
    }
    return (
        <div className="container mx-auto py-10">
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl">Pro</CardTitle>
                    <CardDescription>
                        Upgrade to Pro for higher usage limits
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {subscriptionStatus ? (
                        <Card className="border border-green-200 bg-green-50">
                            <CardHeader className="pb-2">
                                <div className="flex justify-between items-center">
                                    <CardTitle>Current Subscription</CardTitle>
                                    <Badge className="bg-green-600">
                                        Active
                                    </Badge>
                                </div>
                                <CardDescription>
                                    You are currently on the Pro plan.
                                </CardDescription>
                            </CardHeader>
                            <CardContent>
                                <div className="flex items-center gap-2 text-green-700">
                                    <CheckCircle className="h-5 w-5" />
                                    <span>Your subscription is active.</span>
                                </div>
                            </CardContent>
                            <CardFooter>
                                <Button
                                    variant="outline"
                                    className="w-full sm:w-auto"
                                    disabled={billingPortalMutation.isPending}
                                    onClick={() =>
                                        billingPortalMutation.mutate()
                                    }
                                >
                                    Manage Subscription
                                </Button>
                            </CardFooter>
                        </Card>
                    ) : (
                        <div className="grid gap-6 md:grid-cols-2">
                            <Card>
                                <CardHeader>
                                    <CardTitle>Monthly</CardTitle>
                                    <CardDescription>
                                        Billed monthly
                                    </CardDescription>
                                </CardHeader>
                                <CardContent>
                                    <div className="mb-4">
                                        <span className="text-3xl font-bold">
                                            $5
                                        </span>
                                        <span className="text-gray-500">
                                            {" "}
                                            / month
                                        </span>
                                    </div>
                                    <ul className="space-y-2">
                                        <li className="flex items-center gap-2">
                                            <CheckCircle className="h-5 w-5 text-green-600" />
                                            <span>
                                                Create up to 20 JSON files
                                            </span>
                                        </li>
                                    </ul>
                                </CardContent>
                                <CardFooter>
                                    <Button
                                        onClick={async () =>
                                            await checkout(
                                                plans.Monthly.priceId,
                                            )
                                        }
                                        className="w-full"
                                    >
                                        Subscribe
                                    </Button>
                                </CardFooter>
                            </Card>

                            <Card className="border-2 border-black">
                                <CardHeader>
                                    <div className="flex justify-between items-center">
                                        <CardTitle>Yearly</CardTitle>
                                        <Badge>Save 50%</Badge>
                                    </div>
                                    <CardDescription>
                                        Billed annually
                                    </CardDescription>
                                </CardHeader>
                                <CardContent>
                                    <div className="mb-4">
                                        <span className="text-3xl font-bold">
                                            $30
                                        </span>
                                        <span className="text-gray-500">
                                            {" "}
                                            / year
                                        </span>
                                    </div>
                                    <ul className="space-y-2">
                                        <li className="flex items-center gap-2">
                                            <CheckCircle className="h-5 w-5 text-green-600" />
                                            <span>
                                                Create up to 20 JSON files
                                            </span>
                                        </li>
                                        <li className="flex items-center gap-2">
                                            <CheckCircle className="h-5 w-5 text-green-600" />
                                            <span>6 Months Free</span>
                                        </li>
                                    </ul>
                                </CardContent>
                                <CardFooter>
                                    <Button
                                        onClick={async () =>
                                            await checkout(plans.Yearly.priceId)
                                        }
                                        className="w-full"
                                    >
                                        Subscribe
                                    </Button>
                                </CardFooter>
                            </Card>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
