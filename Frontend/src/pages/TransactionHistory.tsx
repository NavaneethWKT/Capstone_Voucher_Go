import React from "react";
import { useState } from "react";
import { useWallet } from "../contexts/WalletContext";
import {
  ArrowUpRight,
  ArrowDownRight,
  CheckCircle,
  XCircle,
  Clock,
  FileText,
} from "lucide-react";

export default function TransactionHistory() {
  const { transactions } = useWallet();
  const [filter, setFilter] = useState<"all" | "purchase" | "topup">("all");
  const [selectedTransaction, setSelectedTransaction] = useState<string | null>(
    null
  );

  const filteredTransactions = transactions.filter((t) => {
    if (filter === "all") return true;
    return t.type === filter;
  });

  const formatDate = (date: Date) => {
    return date.toLocaleDateString("en-US", {
      month: "long",
      day: "numeric",
      year: "numeric",
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

  const getStatusText = (status: string) => {
    return status.charAt(0).toUpperCase() + status.slice(1);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "success":
        return "bg-green-100 text-green-800";
      case "failed":
        return "bg-red-100 text-red-800";
      case "pending":
        return "bg-yellow-100 text-yellow-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const selectedTxn = transactions.find((t) => t.id === selectedTransaction);

  return (
    <div className="max-w-5xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-gray-900 mb-2">Transaction History</h1>
        <p className="text-gray-600">
          View all your wallet transactions and voucher purchases
        </p>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-2">
        <button
          onClick={() => setFilter("all")}
          className={`px-4 py-2 rounded-lg transition-colors ${
            filter === "all"
              ? "bg-blue-600 text-white"
              : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50"
          }`}
        >
          All Transactions
        </button>
        <button
          onClick={() => setFilter("purchase")}
          className={`px-4 py-2 rounded-lg transition-colors ${
            filter === "purchase"
              ? "bg-blue-600 text-white"
              : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50"
          }`}
        >
          Purchases
        </button>
        <button
          onClick={() => setFilter("topup")}
          className={`px-4 py-2 rounded-lg transition-colors ${
            filter === "topup"
              ? "bg-blue-600 text-white"
              : "bg-white text-gray-700 border border-gray-300 hover:bg-gray-50"
          }`}
        >
          Top-ups
        </button>
      </div>

      {/* Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl p-4 border border-gray-200">
          <p className="text-gray-600 text-sm mb-1">Total Transactions</p>
          <p className="text-gray-900">{transactions.length}</p>
        </div>
        <div className="bg-white rounded-xl p-4 border border-gray-200">
          <p className="text-gray-600 text-sm mb-1">Total Purchases</p>
          <p className="text-gray-900">
            {transactions.filter((t) => t.type === "purchase").length}
          </p>
        </div>
        <div className="bg-white rounded-xl p-4 border border-gray-200">
          <p className="text-gray-600 text-sm mb-1">Total Top-ups</p>
          <p className="text-gray-900">
            {transactions.filter((t) => t.type === "topup").length}
          </p>
        </div>
      </div>

      {/* Transactions List */}
      <div className="bg-white rounded-2xl border border-gray-200">
        {filteredTransactions.length > 0 ? (
          <div className="divide-y divide-gray-200">
            {filteredTransactions.map((transaction) => (
              <button
                key={transaction.id}
                onClick={() => setSelectedTransaction(transaction.id)}
                className="w-full p-6 hover:bg-gray-50 transition-colors text-left"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-4 flex-1">
                    <div
                      className={`h-12 w-12 rounded-full flex items-center justify-center ${
                        transaction.type === "topup"
                          ? "bg-green-100"
                          : "bg-blue-100"
                      }`}
                    >
                      {transaction.type === "topup" ? (
                        <ArrowDownRight className="h-6 w-6 text-green-600" />
                      ) : (
                        <ArrowUpRight className="h-6 w-6 text-blue-600" />
                      )}
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center space-x-2 mb-1">
                        <p className="text-gray-900">
                          {transaction.type === "topup"
                            ? "Wallet Top-up"
                            : transaction.voucherName}
                        </p>
                        <span
                          className={`px-2 py-0.5 rounded text-xs ${getStatusColor(
                            transaction.status
                          )}`}
                        >
                          {getStatusText(transaction.status)}
                        </span>
                      </div>
                      <div className="flex items-center space-x-3 text-sm text-gray-500">
                        <span>{formatDate(transaction.date)}</span>
                        <span>•</span>
                        <span>{transaction.paymentId}</span>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <span
                      className={`${
                        transaction.amount > 0
                          ? "text-green-600"
                          : "text-gray-900"
                      }`}
                    >
                      {transaction.amount > 0 ? "+" : ""}$
                      {Math.abs(transaction.amount).toFixed(2)}
                    </span>
                  </div>
                </div>
              </button>
            ))}
          </div>
        ) : (
          <div className="p-12 text-center">
            <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-600 mb-2">No transactions found</p>
            <p className="text-gray-500 text-sm">
              {filter === "all"
                ? "Your transactions will appear here"
                : `No ${filter} transactions found`}
            </p>
          </div>
        )}
      </div>

      {/* Transaction Detail Modal */}
      {selectedTxn && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
          onClick={() => setSelectedTransaction(null)}
        >
          <div
            className="bg-white rounded-2xl p-8 max-w-md w-full"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-gray-900">Transaction Details</h2>
              <button
                onClick={() => setSelectedTransaction(null)}
                className="text-gray-400 hover:text-gray-600"
              >
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div className="flex justify-between py-3 border-b border-gray-200">
                <span className="text-gray-600">Transaction ID</span>
                <span className="text-gray-900">{selectedTxn.id}</span>
              </div>

              <div className="flex justify-between py-3 border-b border-gray-200">
                <span className="text-gray-600">Type</span>
                <span className="text-gray-900 capitalize">
                  {selectedTxn.type}
                </span>
              </div>

              <div className="flex justify-between py-3 border-b border-gray-200">
                <span className="text-gray-600">Status</span>
                <div className="flex items-center space-x-2">
                  {getStatusIcon(selectedTxn.status)}
                  <span className="text-gray-900">
                    {getStatusText(selectedTxn.status)}
                  </span>
                </div>
              </div>

              <div className="flex justify-between py-3 border-b border-gray-200">
                <span className="text-gray-600">Amount</span>
                <span
                  className={`${
                    selectedTxn.amount > 0 ? "text-green-600" : "text-gray-900"
                  }`}
                >
                  {selectedTxn.amount > 0 ? "+" : ""}$
                  {Math.abs(selectedTxn.amount).toFixed(2)}
                </span>
              </div>

              <div className="flex justify-between py-3 border-b border-gray-200">
                <span className="text-gray-600">Date & Time</span>
                <span className="text-gray-900">
                  {formatDate(selectedTxn.date)}
                </span>
              </div>

              {selectedTxn.voucherName && (
                <div className="flex justify-between py-3 border-b border-gray-200">
                  <span className="text-gray-600">Voucher</span>
                  <span className="text-gray-900">
                    {selectedTxn.voucherName}
                  </span>
                </div>
              )}

              <div className="flex justify-between py-3 border-b border-gray-200">
                <span className="text-gray-600">Payment ID</span>
                <span className="text-gray-900">{selectedTxn.paymentId}</span>
              </div>
            </div>

            <button
              onClick={() => setSelectedTransaction(null)}
              className="w-full mt-6 bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
