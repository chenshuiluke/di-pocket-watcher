import { apiClient } from "@/lib/api-client";
import { useMutation, useQueryClient } from "@tanstack/react-query"

export const useGetTransactions = () => {

    return useMutation({
        mutationFn: async ()=> {
            await getTransactions()
        }
    })
}

const getTransactions = async () => {
    await apiClient.get("/transaction-analysis/email")
}