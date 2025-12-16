import React from "react";
import { Link } from "react-router-dom";
import { useWallet } from "../contexts/WalletContext";
import {
  Wallet,
  ArrowUpRight,
  ArrowDownRight,
  CheckCircle,
  XCircle,
  Clock,
} from "lucide-react";
import TransactionHistory from "./TransactionHistory";

export default function WalletPage() {
  const { balance, transactions } = useWallet();

  const formatDate = (date: Date) => {
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "success":
        return <CheckCircle className="h-5 w-5 text-green-600" />;
      case "failed":
        return <XCircle className="h-5 w-5 text-red-600" />;
      case "pending":
        return <Clock className="h-5 w-5 text-yellow-600" />;
      default:
        return null;
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-gray-900 mb-2">Wallet</h1>
        <p className="text-gray-600">
          View your wallet balance and transaction history
        </p>
      </div>

      {/* Balance Card */}
      <div className="bg-gradient-to-br from-blue-600 to-indigo-600 rounded-2xl p-8 text-white">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-blue-100 mb-2">Available Balance</p>
            <h2 className="text-white">${balance.toFixed(2)}</h2>
          </div>
          <div className="h-16 w-16 bg-white/20 rounded-2xl flex items-center justify-center">
            <Wallet className="h-8 w-8" />
          </div>
        </div>
      </div>

      <TransactionHistory />
    </div>
  );
}
