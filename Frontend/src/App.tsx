import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { WalletProvider } from './contexts/WalletContext';
import LoginPage from './pages/LoginPage';
import Dashboard from './pages/Dashboard';
import VoucherBrowse from './pages/VoucherBrowse';
import VoucherDetails from './pages/VoucherDetails';
import PurchaseFlow from './pages/PurchaseFlow';
import WalletPage from './pages/WalletPage';
import TransactionHistory from './pages/TransactionHistory';
import Layout from './components/Layout';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
}

function App() {
  return (
    <Router>
      <AuthProvider>
        <WalletProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Layout />
                </ProtectedRoute>
              }
            >
              <Route index element={<Dashboard />} />
              <Route path="browse" element={<VoucherBrowse />} />
              <Route path="voucher/:id" element={<VoucherDetails />} />
              <Route path="purchase/:id" element={<PurchaseFlow />} />
              <Route path="wallet" element={<WalletPage />} />
              <Route path="transactions" element={<TransactionHistory />} />
            </Route>
          </Routes>
        </WalletProvider>
      </AuthProvider>
    </Router>
  );
}

export default App;