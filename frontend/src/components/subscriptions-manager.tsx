import { useAuth } from "@/hooks/useAuth";
import { Button } from "./ui/button";
import { Skeleton } from "./ui/skeleton";
import { checkout } from "@/lib/api/payment";

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
    const { user, isLoading } = useAuth();
    if (isLoading || !user) {
        return <Skeleton className="h-2 w-5" />;
    }
    return (
        <div>
            {plans.map((plan) => (
                <Button
                    key={plan.name}
                    onClick={async () => await checkout(plan.priceId)}
                >
                    Subscribe {plan.name}
                </Button>
            ))}
        </div>
    );
}
