import { useGetTransactions } from "@/features/transactions/api/transactions"
import { Button } from "@mantine/core"

const Transactions = () => {
    const getTransactions = useGetTransactions()
 return <Button onClick={() => getTransactions.mutate()}>Load Transactions</Button>
}

export default Transactions