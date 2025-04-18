import { useAuth } from "@/hooks/useAuth";
import { Button } from "./ui/button";
import { Skeleton } from "./ui/skeleton";
import { checkout, getSubscriptionStatus } from "@/lib/api/payment";
import { useQuery } from "@tanstack/react-query";
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "./ui/card";
import { Link } from "react-router";

export const plans = [
    {
        name: "Monthly",
        priceId: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "price_1R9FuqDGza3FJhYO6DtchdQh"
            : "",
        price: 5,
    },
    {
        name: "Yearly",
        priceId: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "price_1R9FwkDGza3FJhYOjdukqxiO"
            : "",
        price: 30,
    },
];
export function SubscriptionsManager() {
    const { data: subscriptionStatus, isLoading } = useQuery({
        queryKey: ["subscription-status"],
        queryFn: getSubscriptionStatus,
    });
    if (isLoading) {
        return (
            <div className="container mx-auto py-10 space-y-6">
                <Skeleton className="h-8 w-48" />
                <Skeleton className="h-10 w-40" />
            </div>
        );
    }
    return (
        <div className="container mx-auto py-10">
            {subscriptionStatus ? (
                // --- Subscribed State ---
                <Card className="max-w-md mx-auto">
                    <CardHeader>
                        <CardTitle>Subscription Active</CardTitle>
                        <CardDescription>
                            You have an active subscription. Manage your billing
                            details and view invoices below.
                        </CardDescription>
                    </CardHeader>
                    <CardFooter>
                        <Button asChild>
                            <Link to="https://billing.stripe.com/p/login/test_4gw29W08WaTF2kM8ww">
                                Manage Subscription
                            </Link>
                        </Button>
                    </CardFooter>
                </Card>
            ) : (
                // --- Not Subscribed State ---
                <>
                    <div className="text-center mb-8">
                        <h2 className="text-2xl font-semibold tracking-tight">
                            Choose Your Plan
                        </h2>
                        <p className="text-muted-foreground">
                            Select the plan that works best for you.
                        </p>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-2xl mx-auto">
                        {plans.map((plan) => (
                            <Card key={plan.name}>
                                <CardHeader>
                                    <CardTitle>{plan.name}</CardTitle>
                                </CardHeader>
                                <CardContent>
                                    <p className="text-3xl font-bold">
                                        ${plan.price}
                                        <span className="text-sm font-normal text-muted-foreground">
                                            /
                                            {plan.name === "Monthly"
                                                ? "month"
                                                : "year"}
                                        </span>
                                    </p>
                                </CardContent>
                                <CardFooter>
                                    <Button
                                        key={plan.name}
                                        onClick={async () =>
                                            await checkout(plan.priceId)
                                        }
                                    >
                                        Subscribe {plan.name}
                                    </Button>
                                </CardFooter>
                            </Card>
                        ))}
                    </div>
                </>
            )}
        </div>
    );
}
