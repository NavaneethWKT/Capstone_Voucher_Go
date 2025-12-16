import React from "react";
import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { apiClient, Voucher } from "../services/apiClient";
import { useWallet } from "../contexts/WalletContext";
import {
  CheckCircle,
  AlertCircle,
  Loader2,
  ArrowLeft,
  Wallet,
  Tag,
} from "lucide-react";

type PurchaseStep = "review" | "confirm" | "processing" | "success" | "error";

export default function PurchaseFlow() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { balance, purchaseVoucher } = useWallet();
  const [step, setStep] = useState<PurchaseStep>("review");
  const [error, setError] = useState("");
  const [transactionId, setTransactionId] = useState("");
  const [voucher, setVoucher] = useState<Voucher | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchVoucher = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const voucherId = parseInt(id);
        const voucherData = await apiClient.getVoucherById(voucherId);
        setVoucher(voucherData);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load voucher");
      } finally {
        setLoading(false);
      }
    };

    fetchVoucher();
  }, [id]);

  if (loading) {
    return (
      <div className="max-w-2xl mx-auto text-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600 mx-auto" />
      </div>
    );
  }

  if (!voucher) {
    return (
      <div className="max-w-2xl mx-auto text-center py-12">
        <Tag className="h-12 w-12 text-gray-400 mx-auto mb-4" />
        <h2 className="text-gray-900 mb-2">Voucher not found</h2>
        <Link to="/browse" className="text-blue-600 hover:text-blue-700">
          Browse all vouchers
        </Link>
      </div>
    );
  }

  const handleConfirmPurchase = async () => {
    if (balance < voucher.price) {
      setError("Insufficient balance. Please add funds to your wallet.");
      setStep("error");
      return;
    }

    setStep("processing");
    setError("");

    try {
      const txnId = await purchaseVoucher(
        voucher.id,
        voucher.name,
        voucher.price
      );
      setTransactionId(txnId);
      setStep("success");
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Purchase failed. Please try again."
      );
      setStep("error");
    }
  };

  const remainingBalance = balance - voucher.price;

  return (
    <div className="max-w-2xl mx-auto">
      {/* Back Button (only show before success) */}
      {step !== "success" && (
        <button
          onClick={() => navigate(-1)}
          className="flex items-center space-x-2 text-gray-600 hover:text-gray-900 mb-6"
        >
          <ArrowLeft className="h-4 w-4" />
          <span>Back</span>
        </button>
      )}

      {/* Step Indicator */}
      <div className="mb-8">
        <div className="flex items-center justify-center space-x-2">
          {["review", "confirm", "processing", "success"].map((s, index) => (
            <div key={s} className="flex items-center">
              <div
                className={`h-2 w-2 rounded-full ${
                  step === s ? "bg-blue-600" : "bg-gray-300"
                }`}
              />
              {index < 3 && <div className="h-0.5 w-8 bg-gray-300" />}
            </div>
          ))}
        </div>
      </div>

      {/* Review Step */}
      {step === "review" && (
        <div className="bg-white rounded-2xl p-8 border border-gray-200">
          <h2 className="text-gray-900 mb-6">Review Your Purchase</h2>

          <div className="space-y-4 mb-6">
            <div className="flex justify-between py-3 border-b border-gray-200">
              <span className="text-gray-600">Voucher</span>
              <span className="text-gray-900">{voucher.name}</span>
            </div>
            <div className="flex justify-between py-3 border-b border-gray-200">
              <span className="text-gray-600">Category</span>
              <span className="text-gray-900">{voucher.category}</span>
            </div>
            <div className="flex justify-between py-3 border-b border-gray-200">
              <span className="text-gray-600">Price</span>
              <span className="text-gray-900">${voucher.price.toFixed(2)}</span>
            </div>
            <div className="flex justify-between py-3 border-b border-gray-200">
              <span className="text-gray-600">Valid Until</span>
              <span className="text-gray-900">
                {new Date(voucher.valid_to).toLocaleDateString("en-US", {
                  month: "long",
                  day: "numeric",
                  year: "numeric",
                })}
              </span>
            </div>
          </div>

          <button
            onClick={() => setStep("confirm")}
            className="w-full bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors"
          >
            Continue to Payment
          </button>
        </div>
      )}

      {/* Confirm Step */}
      {step === "confirm" && (
        <div className="bg-white rounded-2xl p-8 border border-gray-200">
          <h2 className="text-gray-900 mb-6">Confirm Payment</h2>

          <div className="bg-blue-50 rounded-xl p-6 mb-6">
            <div className="flex items-center space-x-3 mb-4">
              <Wallet className="h-6 w-6 text-blue-600" />
              <span className="text-gray-900">Wallet Payment</span>
            </div>

            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600">Current Balance</span>
                <span className="text-gray-900">${balance.toFixed(2)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Amount to Deduct</span>
                <span className="text-red-600">
                  -${voucher.price.toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between pt-3 border-t border-blue-200">
                <span className="text-gray-900">Remaining Balance</span>
                <span
                  className={
                    remainingBalance >= 0 ? "text-green-600" : "text-red-600"
                  }
                >
                  ${remainingBalance.toFixed(2)}
                </span>
              </div>
            </div>
          </div>

          {remainingBalance < 0 && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6 flex items-start space-x-3">
              <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-red-900 mb-2">Insufficient balance</p>
                <Link
                  to="/wallet"
                  className="text-red-600 hover:text-red-700 underline"
                >
                  Add funds to your wallet
                </Link>
              </div>
            </div>
          )}

          <div className="flex gap-3">
            <button
              onClick={() => setStep("review")}
              className="flex-1 px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Back
            </button>
            <button
              onClick={handleConfirmPurchase}
              disabled={remainingBalance < 0}
              className="flex-1 bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Confirm Purchase
            </button>
          </div>
        </div>
      )}

      {/* Processing Step */}
      {step === "processing" && (
        <div className="bg-white rounded-2xl p-12 border border-gray-200 text-center">
          <Loader2 className="h-12 w-12 text-blue-600 animate-spin mx-auto mb-4" />
          <h2 className="text-gray-900 mb-2">Processing Payment</h2>
          <p className="text-gray-600">
            Please wait while we process your purchase...
          </p>
        </div>
      )}

      {/* Success Step */}
      {step === "success" && (
        <div className="bg-white rounded-2xl p-12 border border-gray-200 text-center">
          <div className="h-16 w-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <CheckCircle className="h-10 w-10 text-green-600" />
          </div>
          <h2 className="text-gray-900 mb-2">Purchase Successful!</h2>
          <p className="text-gray-600 mb-6">
            Your voucher has been added to your account
          </p>

          <div className="bg-gray-50 rounded-xl p-6 mb-6 text-left">
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-gray-600">Transaction ID</span>
                <span className="text-gray-900">{transactionId}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Voucher</span>
                <span className="text-gray-900">{voucher.name}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Amount Paid</span>
                <span className="text-gray-900">
                  ${voucher.price.toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">New Balance</span>
                <span className="text-green-600">
                  ${remainingBalance.toFixed(2)}
                </span>
              </div>
            </div>
          </div>

          <div className="flex flex-col sm:flex-row gap-3">
            <Link
              to="/wallet"
              className="flex-1 px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-center"
            >
              View Transactions
            </Link>
            <Link
              to="/"
              className="flex-1 bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors text-center"
            >
              Back to Home
            </Link>
          </div>
        </div>
      )}

      {/* Error Step */}
      {step === "error" && (
        <div className="bg-white rounded-2xl p-12 border border-gray-200 text-center">
          <div className="h-16 w-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <AlertCircle className="h-10 w-10 text-red-600" />
          </div>
          <h2 className="text-gray-900 mb-2">Purchase Failed</h2>
          <p className="text-gray-600 mb-6">{error}</p>

          <div className="flex flex-col sm:flex-row gap-3">
            {error.includes("Insufficient") ? (
              <Link
                to="/wallet"
                className="flex-1 bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors text-center"
              >
                Add Funds
              </Link>
            ) : (
              <button
                onClick={() => setStep("confirm")}
                className="flex-1 bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors"
              >
                Try Again
              </button>
            )}
            <Link
              to="/"
              className="flex-1 px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-center"
            >
              Back to Home
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}
