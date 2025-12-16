const BASE_URL = "http://localhost:8080/api/v1";

// Types
export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface Voucher {
  id: number;
  name: string;
  description: string;
  category: string;
  price: number;
  quantity: number;
  valid_from: string;
  valid_to: string;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: number;
  user_id: number;
  voucher_id?: number;
  amount: number;
  transaction_type: string;
  payment_status: string;
  payment_txn_id?: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  message: string;
}

export interface SearchVouchersParams {
  category?: string;
  min_price?: number;
  max_price?: number;
}

export interface BuyVoucherRequest {
  user_id: number;
  voucher_id: number;
}

export interface BuyVoucherResponse {
  transaction: Transaction;
  message: string;
}

// API Client
class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;

    const config: RequestInit = {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        const errorData = await response
          .json()
          .catch(() => ({ error: "Unknown error" }));
        throw new Error(
          errorData.error || `HTTP error! status: ${response.status}`
        );
      }

      const data = await response.json();
      return data as T;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("Network error occurred");
    }
  }

  // Login
  async login(email: string, password: string): Promise<LoginResponse> {
    // Note: Login endpoint needs to be added to the client backend at POST /api/v1/auth/login
    // For now, this will fail until the endpoint is implemented
    const response = await this.request<{
      success: boolean;
      user: User;
      message: string;
    }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ Email: email, Password: password }),
    });
    console.log(response);

    if (!response.success) {
      throw new Error("Login failed");
    }

    return {
      user: response.user,
      message: response.message,
    };
  }

  // Search Vouchers
  async searchVouchers(params: SearchVouchersParams = {}): Promise<Voucher[]> {
    const queryParams = new URLSearchParams();
    if (params.category) queryParams.append("category", params.category);
    if (params.min_price !== undefined)
      queryParams.append("min_price", params.min_price.toString());
    if (params.max_price !== undefined)
      queryParams.append("max_price", params.max_price.toString());

    const queryString = queryParams.toString();
    const endpoint = `/vouchers/search${queryString ? `?${queryString}` : ""}`;

    const response = await this.request<{
      success: boolean;
      vouchers: Voucher[];
    }>(endpoint);

    if (!response.success) {
      throw new Error("Failed to fetch vouchers");
    }

    return response.vouchers;
  }

  // Get Voucher by ID (from search results)
  async getVoucherById(voucherId: number): Promise<Voucher | null> {
    // Since there's no direct endpoint, we'll search and filter
    const vouchers = await this.searchVouchers();
    return vouchers.find((v) => v.id === voucherId) || null;
  }

  // Buy Voucher
  async buyVoucher(
    userId: number,
    voucherId: number
  ): Promise<BuyVoucherResponse> {
    const response = await this.request<{
      success: boolean;
      transaction: Transaction;
      message: string;
    }>("/vouchers/buy", {
      method: "POST",
      body: JSON.stringify({
        user_id: userId,
        voucher_id: voucherId,
      }),
    });

    if (!response.success) {
      throw new Error("Failed to purchase voucher");
    }

    return {
      transaction: response.transaction,
      message: response.message,
    };
  }

  // Get Wallet Balance
  async getBalance(userId: number): Promise<number> {
    const response = await this.request<{ success: boolean; balance: number }>(
      `/wallet/balance/${userId}`
    );

    if (!response.success) {
      throw new Error("Failed to fetch balance");
    }

    return response.balance;
  }

  // List Transactions
  async listTransactions(userId: number): Promise<Transaction[]> {
    const response = await this.request<{
      success: boolean;
      transactions: Transaction[];
    }>(`/transactions/${userId}`);

    if (!response.success) {
      throw new Error("Failed to fetch transactions");
    }

    return response.transactions;
  }
}

// Export singleton instance
export const apiClient = new ApiClient(BASE_URL);
