import React from "react";
import { Link } from "react-router-dom";
import { Voucher } from "../services/apiClient";
import { Calendar, Package } from "lucide-react";
import { ImageWithFallback } from "./figma/ImageWithFallback";

interface VoucherCardProps {
  voucher: Voucher;
}

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

export default function VoucherCard({ voucher }: VoucherCardProps) {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  };

  // Map category to image key
  const getImageKey = (category: string): string => {
    const categoryMap: Record<string, string> = {
      restaurant: "italian restaurant",
      shopping: "online shopping",
      entertainment: "movie theater",
      wellness: "spa wellness",
      travel: "luxury hotel",
    };
    return categoryMap[category.toLowerCase()] || "coffee shop";
  };

  return (
    <Link
      to={`/voucher/${voucher.id}`}
      className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
    >
      <div className="aspect-video w-full overflow-hidden bg-gray-100">
        <ImageWithFallback
          src={voucherImages[getImageKey(voucher.category)]}
          alt={voucher.name}
          className="w-full h-full object-cover"
        />
      </div>
      <div className="p-4">
        <div className="flex items-start justify-between mb-2">
          <div className="flex-1">
            <h3 className="text-gray-900 mb-1">{voucher.name}</h3>
            <span className="inline-block px-2 py-1 bg-blue-50 text-blue-700 rounded text-sm capitalize">
              {voucher.category}
            </span>
          </div>
          <div className="text-blue-600 ml-2">${voucher.price.toFixed(2)}</div>
        </div>

        <div className="flex items-center justify-between text-gray-600 mt-3 pt-3 border-t border-gray-100">
          <div className="flex items-center space-x-1">
            <Calendar className="h-4 w-4" />
            <span className="text-sm">
              Until {formatDate(voucher.valid_to)}
            </span>
          </div>
          <div className="flex items-center space-x-1">
            <Package className="h-4 w-4" />
            <span className="text-sm">{voucher.quantity} left</span>
          </div>
        </div>
      </div>
    </Link>
  );
}
