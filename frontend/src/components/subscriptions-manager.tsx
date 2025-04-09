import { useAuth } from "@/hooks/useAuth";
import { Button } from "./ui/button";
import { Skeleton } from "./ui/skeleton";

export const plans = [
    {
        name: "Monthly",
        link: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "https://buy.stripe.com/test_4gw00KfuidKl79SeUW"
            : "",
        priceId: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "prod_S3MLiWE274RWwB"
            : "",
        price: 5,
    },
    {
        name: "Yearly",
        link: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "https://buy.stripe.com/test_6oE6p83LAdKlbq84gj"
            : "",
        priceId: import.meta.env.VITE_BASE_URL.includes("localhost")
            ? "prod_S3MLiWE274RWwB"
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
                <Button asChild key={plan.name}>
                    <a
                        href={`${plan.link}?prefilled_email=${user.email}`}
                        target="_blank"
                        rel="noreferrer"
                    >
                        Subscribe {plan.name}
                    </a>
                </Button>
            ))}
        </div>
    );
}
