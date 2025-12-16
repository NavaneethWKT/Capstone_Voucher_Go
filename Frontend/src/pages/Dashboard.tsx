import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useWallet } from "../contexts/WalletContext";
import { useAuth } from "../contexts/AuthContext";
import { Wallet, User, Mail, ShoppingBag } from "lucide-react";

export default function Dashboard() {
  const { balance, transactions } = useWallet();
  const { user } = useAuth();

  const purchaseCount = transactions.filter(
    (t) => t.type === "purchase"
  ).length;
  const totalSpent = transactions
    .filter((t) => t.type === "purchase")
    .reduce((sum, t) => sum + Math.abs(t.amount), 0);

  return (
    <div className="space-y-8">
      {/* Welcome Section */}
      <div className="bg-gradient-to-r from-blue-600 to-indigo-600 rounded-2xl p-8 text-white">
        <div className="max-w-4xl">
          <h1 className="text-white mb-2">Welcome back, {user?.name}!</h1>
          <p className="text-blue-100 mb-6">
            Discover amazing deals and vouchers for restaurants, shopping,
            entertainment, and more.
          </p>
        </div>
      </div>

      {/* Profile & Stats Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Profile Card */}
        <div className="bg-white rounded-xl p-6 border border-gray-200">
          <div className="flex items-center space-x-3 mb-4">
            <div className="h-12 w-12 bg-blue-100 rounded-full flex items-center justify-center">
              <User className="h-6 w-6 text-blue-600" />
            </div>
            <div>
              <h3 className="text-gray-900">{user?.name}</h3>
              <p className="text-gray-600 text-sm">{user?.email}</p>
            </div>
          </div>
          <div className="pt-4 border-t border-gray-200 space-y-2">
            <div className="flex items-center space-x-2 text-gray-600 text-sm">
              <Mail className="h-4 w-4" />
              <span>Account ID: {user?.id || "N/A"}</span>
            </div>
          </div>
        </div>

        {/* Wallet Balance */}
        <div className="bg-white rounded-xl p-6 border border-gray-200">
          <div className="flex items-center justify-between mb-2">
            <span className="text-gray-600">Wallet Balance</span>
            <Wallet className="h-5 w-5 text-green-600" />
          </div>
          <div className="text-green-600 mb-2">${balance.toFixed(2)}</div>
          <Link
            to="/wallet"
            className="text-blue-600 hover:text-blue-700 text-sm"
          >
            View Wallet →
          </Link>
        </div>

        {/* Purchase Stats */}
        <div className="bg-white rounded-xl p-6 border border-gray-200">
          <div className="flex items-center justify-between mb-4">
            <span className="text-gray-600">Your Activity</span>
            <ShoppingBag className="h-5 w-5 text-purple-600" />
          </div>
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Total Purchases</span>
              <span className="text-gray-900">{purchaseCount}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Total Spent</span>
              <span className="text-gray-900">${totalSpent.toFixed(2)}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
