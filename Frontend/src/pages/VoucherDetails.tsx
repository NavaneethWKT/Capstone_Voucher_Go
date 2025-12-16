import React, { useState, useEffect } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { apiClient, Voucher } from "../services/apiClient";
import {
  Calendar,
  Package,
  Tag,
  ArrowLeft,
  ShoppingCart,
  Loader2,
} from "lucide-react";
import VoucherCard from "../components/VoucherCard";
import { ImageWithFallback } from "../components/figma/ImageWithFallback";

const voucherImages: Record<string, string> = {
  "coffee shop":
    "https://images.unsplash.com/photo-1453614512568-c4024d13c247?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxjb2ZmZWUlMjBzaG9wfGVufDF8fHx8MTc2NTczODc0OXww&ixlib=rb-4.1.0&q=80&w=1080",
  "online shopping":
    "https://images.unsplash.com/photo-1563013544-824ae1b704d3?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxvbmxpbmUlMjBzaG9wcGluZ3xlbnwxfHx8fDE3NjU3MTg2NjB8MA&ixlib=rb-4.1.0&q=80&w=1080",
  "movie theater":
    "https://images.unsplash.com/photo-1517604931442-7e0c8ed2963c?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxtb3ZpZSUyMHRoZWF0ZXJ8ZW58MXx8fHwxNzY1NzgzOTk2fDA&ixlib=rb-4.1.0&q=80&w=1080",
  "spa wellness":
    "https://images.unsplash.com/photo-1554424518-336ec861b705?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxzcGElMjB3ZWxsbmVzc3xlbnwxfHx8fDE3NjU4MDA2MTZ8MA&ixlib=rb-4.1.0&q=80&w=1080",
  "luxury hotel":
    "https://images.unsplash.com/photo-1561501900-3701fa6a0864?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxsdXh1cnklMjBob3RlbHxlbnwxfHx8fDE3NjU3NzgwNTR8MA&ixlib=rb-4.1.0&q=80&w=1080",
  "italian restaurant":
    "https://images.unsplash.com/photo-1532117472055-4d0734b51f31?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxpdGFsaWFuJTIwcmVzdGF1cmFudHxlbnwxfHx8fDE3NjU4MTIwNjl8MA&ixlib=rb-4.1.0&q=80&w=1080",
  "fashion store":
    "https://images.unsplash.com/photo-1546213290-e1b492ab3eee?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxmYXNoaW9uJTIwc3RvcmV8ZW58MXx8fHwxNzY1ODEyMDcwfDA&ixlib=rb-4.1.0&q=80&w=1080",
  "concert music":
    "https://images.unsplash.com/photo-1631061434620-db65394197e2?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxjb25jZXJ0JTIwbXVzaWN8ZW58MXx8fHwxNzY1NzM0Nzk4fDA&ixlib=rb-4.1.0&q=80&w=1080",
  "yoga class":
    "https://images.unsplash.com/photo-1549576490-b0b4831ef60a?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx5b2dhJTIwY2xhc3N8ZW58MXx8fHwxNzY1NzQwODYzfDA&ixlib=rb-4.1.0&q=80&w=1080",
  "airplane travel":
    "https://images.unsplash.com/photo-1436491865332-7a61a109cc05?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxhaXJwbGFuZSUyMHRyYXZlbHxlbnwxfHx8fDE3NjU3MTAzMjd8MA&ixlib=rb-4.1.0&q=80&w=1080",
  "fast food burger":
    "https://images.unsplash.com/photo-1656439659132-24c68e36b553?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxmYXN0JTIwZm9vZCUyMGJ1cmdlcnxlbnwxfHx8fDE3NjU3NjYzMTR8MA&ixlib=rb-4.1.0&q=80&w=1080",
  "electronics store":
    "https://images.unsplash.com/photo-1571857089849-f6390447191a?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxlbGVjdHJvbmljcyUyMHN0b3JlfGVufDF8fHx8MTc2NTgxMjA3Mnww&ixlib=rb-4.1.0&q=80&w=1080",
  "theme park":
    "https://images.unsplash.com/photo-1502136969935-8d8eef54d77b?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx0aGVtZSUyMHBhcmt8ZW58MXx8fHwxNzY1ODEyMDcyfDA&ixlib=rb-4.1.0&q=80&w=1080",
  "massage therapy":
    "https://images.unsplash.com/photo-1598901986949-f593ff2a31a6?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxtYXNzYWdlJTIwdGhlcmFweXxlbnwxfHx8fDE3NjU3NzU1MDV8MA&ixlib=rb-4.1.0&q=80&w=1080",
  "car rental":
    "https://images.unsplash.com/photo-1565043666747-69f6646db940?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHxjYXIlMjByZW50YWx8ZW58MXx8fHwxNzY1NzMzMTk5fDA&ixlib=rb-4.1.0&q=80&w=1080",
};

export default function VoucherDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [voucher, setVoucher] = useState<Voucher | null>(null);
  const [relatedVouchers, setRelatedVouchers] = useState<Voucher[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchVoucher = async () => {
      if (!id) return;
      setLoading(true);
      setError(null);
      try {
        const voucherId = parseInt(id);
        const voucherData = await apiClient.getVoucherById(voucherId);
        if (!voucherData) {
          setError("Voucher not found");
        } else {
          setVoucher(voucherData);
          // Fetch related vouchers
          const allVouchers = await apiClient.searchVouchers({
            category: voucherData.category,
          });
          setRelatedVouchers(
            allVouchers.filter((v) => v.id !== voucherId).slice(0, 3)
          );
        }
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
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (error || !voucher) {
    return (
      <div className="text-center py-12">
        <Tag className="h-12 w-12 text-gray-400 mx-auto mb-4" />
        <h2 className="text-gray-900 mb-2">Voucher not found</h2>
        <p className="text-gray-600 mb-4">
          {error || "The voucher you're looking for doesn't exist."}
        </p>
        <Link to="/browse" className="text-blue-600 hover:text-blue-700">
          Browse all vouchers
        </Link>
      </div>
    );
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      month: "long",
      day: "numeric",
      year: "numeric",
    });
  };

  return (
    <div className="space-y-8">
      {/* Back Button */}
      <button
        onClick={() => navigate(-1)}
        className="flex items-center space-x-2 text-gray-600 hover:text-gray-900"
      >
        <ArrowLeft className="h-4 w-4" />
        <span>Back</span>
      </button>

      {/* Voucher Details */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Image */}
        <div className="aspect-video rounded-2xl overflow-hidden bg-gray-100">
          <ImageWithFallback
            src={
              voucherImages[voucher.category.toLowerCase()] ||
              voucherImages["coffee shop"]
            }
            alt={voucher.name}
            className="w-full h-full object-cover"
          />
        </div>

        {/* Info */}
        <div className="space-y-6">
          <div>
            <span className="inline-block px-3 py-1 bg-blue-50 text-blue-700 rounded-lg mb-3">
              {voucher.category}
            </span>
            <h1 className="text-gray-900 mb-2">{voucher.name}</h1>
            <p className="text-gray-600">{voucher.description}</p>
          </div>

          <div className="flex items-center space-x-6 py-4 border-y border-gray-200">
            <div className="flex items-center space-x-2 text-gray-600">
              <Calendar className="h-5 w-5" />
              <div>
                <div className="text-sm">Valid Period</div>
                <div className="text-gray-900">
                  {formatDate(voucher.valid_from)} -{" "}
                  {formatDate(voucher.valid_to)}
                </div>
              </div>
            </div>
            <div className="flex items-center space-x-2 text-gray-600">
              <Package className="h-5 w-5" />
              <div>
                <div className="text-sm">Available</div>
                <div className="text-gray-900">{voucher.quantity} left</div>
              </div>
            </div>
          </div>

          <div className="bg-gray-50 rounded-xl p-6">
            <div className="flex items-center justify-between mb-4">
              <span className="text-gray-600">Price</span>
              <span className="text-blue-600">${voucher.price.toFixed(2)}</span>
            </div>
            <button
              onClick={() => navigate(`/purchase/${voucher.id}`)}
              className="w-full bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors flex items-center justify-center space-x-2"
            >
              <ShoppingCart className="h-5 w-5" />
              <span>Purchase Now</span>
            </button>
          </div>

          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <p className="text-yellow-900">
              <strong>Note:</strong> This voucher will be added to your account
              immediately after purchase. Please check the terms and conditions
              before purchasing.
            </p>
          </div>
        </div>
      </div>

      {/* Related Vouchers */}
      {relatedVouchers.length > 0 && (
        <div>
          <h2 className="text-gray-900 mb-4">Similar Vouchers</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {relatedVouchers.map((v) => (
              <VoucherCard voucher={v} />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
