import React, {
  createContext,
  useContext,
  useState,
  ReactNode,
  useEffect,
  useCallback,
} from "react";
import { useAuth } from "./AuthContext";
import {
  apiClient,
  Transaction as ApiTransaction,
} from "../services/apiClient";

export interface Transaction {
  id: string;
  type: "purchase" | "topup";
  amount: number;
  date: Date;
  status: "success" | "failed" | "pending";
  voucherName?: string;
  voucherId?: string;
  paymentId: string;
}

interface WalletContextType {
  balance: number;
  transactions: Transaction[];
  loading: boolean;
  refreshBalance: () => Promise<void>;
  refreshTransactions: () => Promise<void>;
  addFunds: (amount: number) => Promise<void>;
  purchaseVoucher: (
    voucherId: number,
    voucherName: string,
    amount: number
  ) => Promise<void>;
}

const WalletContext = createContext<WalletContextType | undefined>(undefined);

// Convert API transaction to frontend transaction
function convertTransaction(
  apiTxn: ApiTransaction,
  voucherName?: string
): Transaction {
  return {
    id: apiTxn.id.toString(),
    type: apiTxn.transaction_type === "purchase" ? "purchase" : "topup",
    amount:
      apiTxn.transaction_type === "purchase" ? -apiTxn.amount : apiTxn.amount,
    date: new Date(apiTxn.created_at),
    status: apiTxn.payment_status as "success" | "failed" | "pending",
    voucherName,
    voucherId: apiTxn.voucher_id?.toString(),
    paymentId: apiTxn.payment_txn_id || `TXN_${apiTxn.id}`,
  };
}

export function WalletProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [balance, setBalance] = useState<number>(0);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);

  const refreshBalance = useCallback(async () => {
    if (!user) return;
    try {
      const bal = await apiClient.getBalance(user.id);
      setBalance(bal);
    } catch (error) {
      console.error("Failed to fetch balance:", error);
    }
  }, [user]);

  const refreshTransactions = useCallback(async () => {
    if (!user) return;
    try {
      const apiTransactions = await apiClient.listTransactions(user.id);
      // Note: We don't have voucher names in the transaction response
      // You may need to fetch voucher details separately or update the API
      const converted = apiTransactions.map((txn) => convertTransaction(txn));
      setTransactions(converted);
    } catch (error) {
      console.error("Failed to fetch transactions:", error);
    }
  }, [user]);

  useEffect(() => {
    if (user) {
      setLoading(true);
      Promise.all([refreshBalance(), refreshTransactions()]).finally(() => {
        setLoading(false);
      });
    } else {
      setBalance(0);
      setTransactions([]);
    }
  }, [user, refreshBalance, refreshTransactions]);

  const addFunds = async (amount: number) => {
    // Note: Top-up functionality needs to be implemented in the backend
    // For now, this is a placeholder
    throw new Error("Top-up functionality not yet implemented in backend");
  };

  const purchaseVoucher = async (
    voucherId: number,
    voucherName: string,
    amount: number
  ) => {
    if (!user) {
      throw new Error("User not authenticated");
    }

    try {
      const response = await apiClient.buyVoucher(user.id, voucherId);

      // Refresh balance and transactions after purchase
      await Promise.all([refreshBalance(), refreshTransactions()]);

      // Return transaction ID for success page
      return response.transaction.id.toString();
    } catch (error) {
      throw error;
    }
  };

  return (
    <WalletContext.Provider
      value={{
        balance,
        transactions,
        loading,
        refreshBalance,
        refreshTransactions,
        addFunds,
        purchaseVoucher,
      }}
    >
      {children}
    </WalletContext.Provider>
  );
}

export function useWallet() {
  const context = useContext(WalletContext);
  if (context === undefined) {
    throw new Error("useWallet must be used within a WalletProvider");
  }
  return context;
}
