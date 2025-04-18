import { paymentSuccessful } from "@/lib/api/payment";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CheckCircle } from "lucide-react";
import { useNavigate } from "react-router";

export function Success() {
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    // const paymentSuccessMutation = useMutation({
    //     mutationFn: paymentSuccessful,
    //     onSuccess: () => {
    //         navigate("/app/account");
    //         queryClient.invalidateQueries({
    //             queryKey: ["subscription-status"],
    //         });
    //     },
    // });
    // paymentSuccessMutation.mutate();
    return (
        <div className="flex flex-col items-center justify-center min-h-screen gap-5 text-center">
            <title>Payment Successful - RestJSON</title>
            <CheckCircle className="w-16 h-16 text-green-500" />
            <h1 className="text-2xl font-bold">Payment Successful!</h1>
            <p className="text-gray-600">
                Thank you for your purchase. You will be redirected shortly.
            </p>
        </div>
    );
}
